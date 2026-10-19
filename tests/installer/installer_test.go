package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

type fakeState struct {
	Containers                 []container
	Calls                      [][]string
	Installed, Candidate, Mode string
	Current                    string
	Restarts, SinceRestart     int
	Address                    string
	BadBackup                  bool
}

func readState(path string) fakeState {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var state fakeState
	if err := json.Unmarshal(data, &state); err != nil {
		panic(err)
	}
	return state
}

func writeState(path string, state fakeState) {
	data, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		panic(err)
	}
}

// The test executable acts as Docker. Each command records its arguments and state.
func TestDockerHelper(t *testing.T) {
	path := os.Getenv("INSTALLER_FAKE_DOCKER_STATE")
	if path == "" {
		return
	}
	var args []string
	for n, value := range os.Args {
		if value == "--" {
			args = os.Args[n+1:]
			break
		}
	}
	state := readState(path)
	state.Calls = append(state.Calls, args)
	output, err := fakeCommand(&state, args)
	writeState(path, state)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(output)
	os.Exit(0)
}

func fakeCommand(s *fakeState, args []string) (string, error) {
	switch args[0] {
	case "ps":
		var ids []string
		for _, item := range s.Containers {
			ids = append(ids, item.ID)
		}
		return strings.Join(ids, "\n"), nil
	case "inspect":
		for _, item := range s.Containers {
			if args[1] != item.ID {
				continue
			}
			if s.Restarts > 0 {
				s.SinceRestart++
				if s.Mode == "stopped" && s.Current == "candidate" {
					item.State.Running, item.State.Status = false, "exited"
				}
				if s.Mode == "restart" && s.Current == "candidate" && s.SinceRestart == 1 {
					item.State.Running, item.State.Restarting, item.State.Status = false, true, "restarting"
				}
				_, port, _ := net.SplitHostPort(s.Address)
				item.NetworkSettings.Ports = map[string][]struct{ HostIp, HostPort string }{"8080/tcp": {{"0.0.0.0", port}}}
				item.RestartCount = s.Restarts
			}
			data, _ := json.Marshal([]container{item})
			return string(data), nil
		}
		return "", fmt.Errorf("Unknown immutable ID: %s", args[1])
	case "cp":
		if strings.HasSuffix(args[1], ":"+binaryPath) {
			data := append([]byte{}, elfHeader...)
			if s.BadBackup {
				data[0] = 0
			}
			return "", os.WriteFile(args[2], data, 0755)
		}
		isBackup := !strings.Contains(args[1], "payload-")
		if isBackup && s.Mode == "restore-failure" {
			return "", fmt.Errorf("Synthetic restore copy failure")
		}
		if strings.HasSuffix(args[2], ":"+binaryPath) {
			if isBackup {
				s.Current = "original"
			} else {
				s.Current = "candidate"
			}
		}
		return "", nil
	case "exec":
		if args[1] == "-u" {
			switch args[4] {
			case "chmod":
				if s.Mode == "chmod-failure" {
					return "", fmt.Errorf("Synthetic chmod failure")
				}
			case "sh":
				if s.Current == "candidate" {
					s.Current = "original"
				} else {
					s.Current = "candidate"
				}
			}
			return "", nil
		}
		if args[2] == stagedPath {
			return "Sub2API " + s.Candidate, nil
		}
		return "Sub2API " + s.Installed, nil
	case "restart":
		s.Restarts++
		s.SinceRestart = 0
		return args[1] + "\n", nil
	}
	return "", fmt.Errorf("Unexpected Docker command: %v", args)
}

type fixture struct {
	i           *installer
	statePath   string
	out, errOut bytes.Buffer
	cancel      context.CancelFunc
	ctx         context.Context
	source      string
}

func newFixture(t *testing.T, state fakeState, input string) *fixture {
	t.Helper()
	f := &fixture{}
	f.ctx, f.cancel = context.WithCancel(context.Background())
	t.Cleanup(f.cancel)
	root := t.TempDir()
	f.statePath = filepath.Join(root, "state.json")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := readState(f.statePath)
		if current.Current == "candidate" {
			if current.Mode == "redirect" && r.URL.Path == "/health" {
				http.Redirect(w, r, "/redirected", http.StatusFound)
				return
			}
			if current.Mode == "interrupt" {
				f.cancel()
				http.Error(w, "interrupted", 503)
				return
			}
			if current.Mode == "timeout" || current.Mode == "restore-failure" {
				if r.URL.Path == "/health" {
					fmt.Fprint(w, `{"status":"ok"}`)
				} else {
					fmt.Fprint(w, `{"data":{"needs_setup":true}}`)
				}
				return
			}
		}
		if r.URL.Path == "/health" || r.URL.Path == "/redirected" {
			fmt.Fprint(w, `{"status":"ok"}`)
		} else {
			fmt.Fprint(w, `{"data":{"needs_setup":false}}`)
		}
	}))
	t.Cleanup(server.Close)
	state.Address = strings.TrimPrefix(server.URL, "http://")
	if state.Containers == nil {
		state.Containers = []container{runningContainer("immutable-target-123456", "/Sub2API-main")}
	}
	if state.Installed == "" {
		state.Installed = "0.2.16"
	}
	if state.Candidate == "" {
		state.Candidate = "0.2.17"
	}
	state.Current = "original"
	writeState(f.statePath, state)
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	quote := func(value string) string { return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'" }
	shim := "#!/bin/sh\nexec " + quote(executable) + " -test.run '^TestDockerHelper$' -- \"$@\"\n"
	if err := os.WriteFile(filepath.Join(root, "docker"), []byte(shim), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", root+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("INSTALLER_FAKE_DOCKER_STATE", f.statePath)
	f.i = &installer{in: strings.NewReader(input), out: &f.out, errOut: &f.errOut, timeout: 600 * time.Millisecond, poll: 10 * time.Millisecond}
	payloadDirectory := filepath.Join(root, "payload-test")
	if err := os.Mkdir(payloadDirectory, 0700); err != nil {
		t.Fatal(err)
	}
	f.source = filepath.Join(payloadDirectory, "sub2api")
	if err := os.WriteFile(f.source, elfHeader, 0755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		for _, line := range strings.Split(f.out.String(), "\n") {
			if path, ok := strings.CutPrefix(line, "Backup: "); ok {
				os.RemoveAll(filepath.Dir(path))
			}
		}
	})
	return f
}

func runningContainer(id, name string) container {
	item := container{ID: id, Name: name}
	item.State.Running, item.State.Status = true, "running"
	return item
}

func hasCall(calls [][]string, values ...string) bool {
	for _, call := range calls {
		if strings.Join(call, "\x00") == strings.Join(values, "\x00") {
			return true
		}
	}
	return false
}

func TestSelection(t *testing.T) {
	cases := []struct {
		name, input, wanted string
		items               []container
		more                bool
	}{
		{"filtered", "\n99\n1\n", "target", []container{runningContainer("other", "/db"), runningContainer("target", "/SUB2API-main")}, true},
		{"numbered-more", "2\n1\n", "other", []container{runningContainer("other", "/db"), runningContainer("target", "/sub2api-main")}, true},
		{"no-match", "2\n", "second", []container{runningContainer("first", "/db"), runningContainer("second", "/server")}, false},
		{"all-match", "1\n", "target", []container{runningContainer("target", "/sub2api")}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t, fakeState{Containers: tc.items}, tc.input)
			item, err := f.i.selectContainer(f.ctx)
			if err != nil || item.ID != tc.wanted {
				t.Fatalf("Selected %q, error %v", item.ID, err)
			}
			if strings.Contains(f.out.String(), "More") != tc.more {
				t.Fatal(f.out.String())
			}
			if strings.Contains(f.out.String(), "Cancel") {
				t.Fatal(f.out.String())
			}
		})
	}
	t.Run("empty-does-not-select", func(t *testing.T) {
		f := newFixture(t, fakeState{}, "\n")
		if _, err := f.i.selectContainer(f.ctx); err != io.EOF {
			t.Fatalf("Expected EOF, got %v", err)
		}
	})
	t.Run("none", func(t *testing.T) {
		f := newFixture(t, fakeState{Containers: []container{}}, "1\n")
		if _, err := f.i.selectContainer(f.ctx); err == nil || !strings.Contains(err.Error(), "No running containers") {
			t.Fatal(err)
		}
	})
}

func TestInstallTransaction(t *testing.T) {
	cases := []struct{ name, installed, mode, expected string }{
		{"upgrade", "0.2.16", "", ""},
		{"reinstall", "0.2.17", "", ""},
		{"restart-address-change", "0.2.16", "restart", ""},
		{"redirect-rejected", "0.2.16", "redirect", "HTTP 302"},
		{"setup-timeout-restores", "0.2.16", "timeout", "Readiness timed out"},
		{"stopped-restores-directly", "0.2.16", "stopped", "Container exited"},
		{"interrupt-restores", "0.2.16", "interrupt", "context canceled"},
		{"restore-failure-keeps-path", "0.2.16", "restore-failure", "Restoration failed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t, fakeState{Installed: tc.installed, Mode: tc.mode}, "1\n")
			err := f.i.install(f.ctx, f.source)
			if tc.expected == "" && err != nil {
				t.Fatal(err, f.out.String(), f.errOut.String())
			}
			if tc.expected != "" && (err == nil || !strings.Contains(err.Error(), tc.expected)) {
				t.Fatalf("Expected %s, got %v", tc.expected, err)
			}
			s := readState(f.statePath)
			cid := s.Containers[0].ID
			if !hasCall(s.Calls, "exec", "-u", "0", cid, "rm", "-f", stagedPath) {
				t.Fatal("Candidate was not cleaned")
			}
			backupIndex, mutationIndex := -1, -1
			for n, call := range s.Calls {
				if call[0] == "cp" && call[1] == cid+":"+binaryPath {
					backupIndex = n
				}
				if call[0] == "exec" && len(call) > 4 && call[4] == "sh" && mutationIndex < 0 {
					mutationIndex = n
				}
				for _, arg := range call {
					if arg == "Sub2API-main" || strings.HasPrefix(arg, "Sub2API-main:") {
						t.Fatal("Used mutable name", call)
					}
				}
			}
			if backupIndex < 0 || mutationIndex <= backupIndex {
				t.Fatal("Missing backup before replacement", s.Calls)
			}
			if tc.expected == "" {
				if s.Current != "candidate" || !strings.Contains(f.out.String(), "Installed.") {
					t.Fatal(f.out.String(), s)
				}
			} else if tc.mode != "restore-failure" {
				if s.Current != "original" || !strings.Contains(f.errOut.String(), "Original executable restored.") {
					t.Fatal(f.errOut.String(), s)
				}
			}
			if tc.mode == "stopped" {
				found := false
				for _, call := range s.Calls {
					if call[0] == "cp" && call[2] == cid+":"+binaryPath {
						found = true
					}
				}
				if !found {
					t.Fatal("Stopped restoration did not copy directly")
				}
			}
			if tc.mode == "restart" && !strings.Contains(f.out.String(), "state=restarting") {
				t.Fatal(f.out.String())
			}
			if tc.mode == "restore-failure" && !strings.Contains(err.Error(), "Backup: /tmp/") {
				t.Fatal(err)
			}
			for _, line := range strings.Split(f.out.String(), "\n") {
				if path, ok := strings.CutPrefix(line, "Backup: "); ok {
					info, err := os.Stat(filepath.Dir(path))
					if err != nil || info.Mode().Perm() != 0700 {
						t.Fatal("Backup directory is not private", err)
					}
					if err := checkELF(path); err != nil {
						t.Fatal("Backup was not retained", err)
					}
				}
			}
		})
	}
}

func TestRejectionBeforeReplacement(t *testing.T) {
	cases := []struct {
		name, installed, candidate, mode, expected string
		badELF, badBackup                          bool
	}{
		{name: "installed-version", installed: "0.2.18", expected: "Expected Sub2API"},
		{name: "candidate-version", candidate: "0.2.16", expected: "Expected candidate"},
		{name: "candidate-elf", badELF: true, expected: "Linux/amd64"},
		{name: "backup-elf", badBackup: true, expected: "existing executable"},
		{name: "chmod-cleanup", mode: "chmod-failure", expected: "Synthetic chmod failure"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t, fakeState{Installed: tc.installed, Candidate: tc.candidate, Mode: tc.mode, BadBackup: tc.badBackup}, "1\n")
			if tc.badELF {
				if err := os.WriteFile(f.source, []byte("wrong architecture"), 0755); err != nil {
					t.Fatal(err)
				}
			}
			err := f.i.install(f.ctx, f.source)
			if err == nil || !strings.Contains(err.Error(), tc.expected) {
				t.Fatalf("Expected %s, got %v", tc.expected, err)
			}
			s := readState(f.statePath)
			if s.Restarts != 0 || s.Current != "original" {
				t.Fatal("Rejected installation mutated executable", s)
			}
			if tc.name == "candidate-version" || tc.name == "chmod-cleanup" {
				if !hasCall(s.Calls, "exec", "-u", "0", s.Containers[0].ID, "rm", "-f", stagedPath) {
					t.Fatal("Missing candidate cleanup")
				}
			}
			// The backup-ELF rejection happens before its path is printed.
			for _, call := range s.Calls {
				if call[0] == "cp" && strings.HasSuffix(call[1], ":"+binaryPath) {
					t.Cleanup(func() { os.RemoveAll(filepath.Dir(call[2])) })
				}
			}
		})
	}
}

func TestEntryAndEmbeddedPayload(t *testing.T) {
	f := newFixture(t, fakeState{}, "1\n")
	if err := run(f.ctx, strings.NewReader("1\n"), &f.out, &f.errOut); err != nil {
		t.Fatal(err)
	}
	for _, call := range readState(f.statePath).Calls {
		if call[0] == "cp" && strings.Contains(call[1], "payload-") {
			if _, err := os.Stat(filepath.Dir(call[1])); !os.IsNotExist(err) {
				t.Fatal("Entry retained extracted payload", err)
			}
		}
	}
	licenses, err := assets.ReadDir("LICENSE")
	if err != nil || len(licenses) != 5 {
		t.Fatal("Missing embedded license texts", err, licenses)
	}
	for _, entry := range licenses {
		data, err := assets.ReadFile("LICENSE/" + entry.Name())
		if err != nil || len(data) == 0 {
			t.Fatal("Empty license", entry.Name(), err)
		}
	}
}

func TestEndpoint(t *testing.T) {
	item := runningContainer("id", "/sub2api")
	item.Config.Env = []string{"SERVER_PORT=9000"}
	item.HostConfig.NetworkMode = "host"
	if address, err := endpoint(item); err != nil || address != "127.0.0.1:9000" {
		t.Fatal(address, err)
	}
	item.HostConfig.NetworkMode = "bridge"
	item.NetworkSettings.Networks = map[string]struct{ IPAddress string }{"default": {"172.17.0.2"}}
	if address, err := endpoint(item); err != nil || address != "172.17.0.2:9000" {
		t.Fatal(address, err)
	}
	item.NetworkSettings.Ports = map[string][]struct{ HostIp, HostPort string }{"9000/tcp": {{"::", strconv.Itoa(19000)}}}
	if address, err := endpoint(item); err != nil || address != "127.0.0.1:19000" {
		t.Fatal(address, err)
	}
}
