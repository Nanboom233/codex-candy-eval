package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const binaryPath = "/app/sub2api"
const stagedPath = "/app/sub2api.codex-eval-new"
const readinessTimeout = 120 * time.Second

//go:embed payload.gz LICENSE/*
var assets embed.FS

var elfHeader = []byte{0x7f, 'E', 'L', 'F', 2, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0x3e, 0}
var installedVersion = regexp.MustCompile(`\b0\.2\.(16|17)\b`)
var candidateVersion = regexp.MustCompile(`\b0\.2\.17\b`)

type container struct {
	ID    string `json:"Id"`
	Name  string
	State struct {
		Running, Restarting, Paused, OOMKilled bool
		Status, Error                          string
		ExitCode                               int
	}
	RestartCount    int
	Config          struct{ Env []string }
	HostConfig      struct{ NetworkMode string }
	NetworkSettings struct {
		Ports    map[string][]struct{ HostIp, HostPort string }
		Networks map[string]struct{ IPAddress string }
	}
}

type installer struct {
	in          io.Reader
	out, errOut io.Writer
	timeout     time.Duration
	poll        time.Duration
}

func (i *installer) docker(ctx context.Context, quiet bool, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", args...)
	// Only the installation context cancels Docker. Terminal signals must not interrupt rollback commands.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if !quiet {
		fmt.Fprint(i.out, stdout.String())
		fmt.Fprint(i.errOut, stderr.String())
	}
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	if err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = strings.TrimSpace(stdout.String())
		}
		return "", fmt.Errorf("docker %s failed: %w: %s", strings.Join(args[:min(2, len(args))], " "), err, message)
	}
	return stdout.String(), nil
}

func (i *installer) inspect(ctx context.Context, cid string) (container, error) {
	var items []container
	output, err := i.docker(ctx, true, "inspect", cid)
	if err != nil {
		return container{}, err
	}
	if err := json.Unmarshal([]byte(output), &items); err != nil {
		return container{}, fmt.Errorf("Read Docker inspect: %w", err)
	}
	if len(items) != 1 {
		return container{}, errors.New("Expected one container in Docker inspect.")
	}
	return items[0], nil
}

func (i *installer) selectContainer(ctx context.Context) (container, error) {
	ids, err := i.docker(ctx, true, "ps", "-q")
	if err != nil {
		return container{}, err
	}
	var all, matches []container
	for _, cid := range strings.Fields(ids) {
		item, err := i.inspect(ctx, cid)
		if err != nil {
			return container{}, err
		}
		all = append(all, item)
		if strings.Contains(strings.ToLower(strings.TrimPrefix(item.Name, "/")), "sub2api") {
			matches = append(matches, item)
		}
	}
	if len(all) == 0 {
		return container{}, errors.New("No running containers.")
	}
	shown := matches
	if len(shown) == 0 {
		shown = all
	}
	type input struct {
		line string
		err  error
	}
	lines := make(chan input)
	go func() {
		scanner := bufio.NewScanner(i.in)
		for scanner.Scan() {
			select {
			case lines <- input{line: scanner.Text()}:
			case <-ctx.Done():
				return
			}
		}
		err := scanner.Err()
		if err == nil {
			err = io.EOF
		}
		select {
		case lines <- input{err: err}:
		case <-ctx.Done():
		}
	}()
	for {
		for n, item := range shown {
			fmt.Fprintf(i.out, "[%d] %s (%s)\n", n+1, strings.TrimPrefix(item.Name, "/"), shortID(item.ID))
		}
		more := len(shown) < len(all)
		if more {
			fmt.Fprintf(i.out, "[%d] More\n", len(shown)+1)
		}
		for {
			fmt.Fprint(i.out, "Select: ")
			var value input
			select {
			case <-ctx.Done():
				return container{}, ctx.Err()
			case value = <-lines:
			}
			if value.err != nil {
				return container{}, value.err
			}
			n, err := strconv.Atoi(strings.TrimSpace(value.line))
			if err != nil {
				continue
			}
			if more && n == len(shown)+1 {
				shown = all
				break
			}
			if n > 0 && n <= len(shown) {
				return shown[n-1], nil
			}
		}
	}
}

func shortID(cid string) string { return cid[:min(12, len(cid))] }

func endpoint(item container) (string, error) {
	port := 8080
	for _, entry := range item.Config.Env {
		if value, ok := strings.CutPrefix(entry, "SERVER_PORT="); ok {
			parsed, err := strconv.Atoi(value)
			if err != nil {
				return "", fmt.Errorf("Invalid SERVER_PORT: %w", err)
			}
			port = parsed
		}
	}
	if bindings := item.NetworkSettings.Ports[fmt.Sprintf("%d/tcp", port)]; len(bindings) > 0 {
		host := bindings[0].HostIp
		if host == "" || host == "0.0.0.0" || host == "::" {
			host = "127.0.0.1"
		}
		if _, err := strconv.Atoi(bindings[0].HostPort); err != nil {
			return "", err
		}
		return net.JoinHostPort(host, bindings[0].HostPort), nil
	}
	if item.HostConfig.NetworkMode == "host" {
		return net.JoinHostPort("127.0.0.1", strconv.Itoa(port)), nil
	}
	for _, network := range item.NetworkSettings.Networks {
		if network.IPAddress != "" {
			return net.JoinHostPort(network.IPAddress, strconv.Itoa(port)), nil
		}
	}
	return "", errors.New("No reachable application address in Docker inspect.")
}

func requestReady(ctx context.Context, address, path string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+path, nil)
	if err != nil {
		return err
	}
	// Inspect supplies a direct container address. Do not send local checks through a host proxy.
	client := &http.Client{
		Transport:     &http.Transport{DisableKeepAlives: true},
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", response.StatusCode)
	}
	var data struct {
		Status string
		Data   struct {
			NeedsSetup *bool `json:"needs_setup"`
		}
	}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return err
	}
	if path == "/health" && data.Status != "ok" {
		return fmt.Errorf("Health status: %s", data.Status)
	}
	if path == "/setup/status" && (data.Data.NeedsSetup == nil || *data.Data.NeedsSetup) {
		return errors.New("Application setup is incomplete.")
	}
	return nil
}

func (i *installer) readiness(ctx context.Context, cid string) error {
	started := time.Now()
	waitCtx, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()
	lastError, previous := "not ready", ""
	nextProgress := time.Time{}
	var item container
	for waitCtx.Err() == nil {
		var err error
		item, err = i.inspect(waitCtx, cid)
		if err != nil {
			if waitCtx.Err() != nil {
				break
			}
			return err
		}
		state := item.State
		if !state.Running && !state.Restarting {
			detail := state.Error
			if detail == "" {
				detail = "no runtime error reported"
			}
			return fmt.Errorf("Container %s after %.1fs: exit=%d, OOMKilled=%t, error=%s; last request: %s", state.Status, time.Since(started).Seconds(), state.ExitCode, state.OOMKilled, detail, lastError)
		}
		if state.Paused {
			return fmt.Errorf("Container paused after %.1fs; last request: %s", time.Since(started).Seconds(), lastError)
		}
		address := "no address yet"
		if state.Restarting {
			err = errors.New("Container is restarting")
		} else {
			address, err = endpoint(item)
			if err != nil {
				address = "no address yet"
			}
			if err == nil {
				for _, path := range []string{"/health", "/setup/status"} {
					if err = requestReady(waitCtx, address, path); err != nil {
						err = fmt.Errorf("%s at %s: %w", path, address, err)
						break
					}
				}
			}
		}
		if err == nil {
			fmt.Fprintf(i.out, "Application ready after %.1fs: health OK, setup complete.\n", time.Since(started).Seconds())
			return nil
		}
		lastError = err.Error()
		current := fmt.Sprintf("state=%s, restarts=%d, address=%s", state.Status, item.RestartCount, address)
		if time.Now().After(nextProgress) || current != previous {
			fmt.Fprintf(i.out, "Application is starting; waiting for health and setup [%.1f/%gs, %s].\n", time.Since(started).Seconds(), i.timeout.Seconds(), current)
			nextProgress, previous = time.Now().Add(5*time.Second), current
		}
		timer := time.NewTimer(i.poll)
		select {
		case <-waitCtx.Done():
			timer.Stop()
		case <-timer.C:
		}
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return fmt.Errorf("Readiness timed out after %.1fs: state=%s, restarts=%d, last request: %s", time.Since(started).Seconds(), item.State.Status, item.RestartCount, lastError)
}

func checkELF(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	header := make([]byte, len(elfHeader))
	if _, err := io.ReadFull(f, header); err != nil {
		return errors.New("Expected a static Linux/amd64 ELF executable.")
	}
	if !bytes.Equal(header, elfHeader) {
		return errors.New("Expected a static Linux/amd64 ELF executable.")
	}
	return nil
}

func extractPayload() (string, func(), error) {
	directory, err := os.MkdirTemp("/tmp", "sub2api-codex-eval-payload-")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(directory) }
	path := filepath.Join(directory, "sub2api")
	err = func() error {
		compressed, err := assets.Open("payload.gz")
		if err != nil {
			return err
		}
		defer compressed.Close()
		reader, err := gzip.NewReader(compressed)
		if err != nil {
			return err
		}
		defer reader.Close()
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0755)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(f, reader)
		closeErr := f.Close()
		return errors.Join(copyErr, closeErr)
	}()
	if err != nil {
		cleanup()
		return "", nil, err
	}
	return path, cleanup, nil
}

func (i *installer) replace(ctx context.Context, cid, source string) error {
	item, err := i.inspect(ctx, cid)
	if err != nil {
		return err
	}
	if !item.State.Running {
		_, err = i.docker(ctx, false, "cp", source, cid+":"+binaryPath)
		return err
	}
	if _, err = i.docker(ctx, false, "cp", source, cid+":"+stagedPath); err != nil {
		return err
	}
	_, err = i.docker(ctx, false, "exec", "-u", "0", cid, "sh", "-c", "chmod 755 "+stagedPath+" && mv -f "+stagedPath+" "+binaryPath)
	return err
}

func (i *installer) validateCandidate(ctx context.Context, cid, source string) (result error) {
	if err := checkELF(source); err != nil {
		return err
	}
	if _, err := i.docker(ctx, false, "cp", source, cid+":"+stagedPath); err != nil {
		return err
	}
	defer func() {
		// An interrupt must not skip removal of the staged candidate.
		_, err := i.docker(context.Background(), false, "exec", "-u", "0", cid, "rm", "-f", stagedPath)
		result = errors.Join(result, err)
	}()
	if _, err := i.docker(ctx, false, "exec", "-u", "0", cid, "chmod", "755", stagedPath); err != nil {
		return err
	}
	version, err := i.docker(ctx, true, "exec", cid, stagedPath, "-version")
	if err != nil {
		return err
	}
	if !candidateVersion.MatchString(version) {
		return fmt.Errorf("Expected candidate Sub2API v0.2.17, received: %s", strings.TrimSpace(version))
	}
	return nil
}

func (i *installer) install(ctx context.Context, source string) error {
	item, err := i.selectContainer(ctx)
	if err != nil {
		return err
	}
	cid := item.ID
	fmt.Fprintf(i.out, "Container: %s (%s)\n", strings.TrimPrefix(item.Name, "/"), shortID(cid))
	if !item.State.Running {
		return errors.New("Select a running container.")
	}
	fmt.Fprintln(i.out, "[1/4] Check v0.2.17 and the local binary.")
	version, err := i.docker(ctx, true, "exec", cid, binaryPath, "-version")
	if err != nil {
		return err
	}
	if !installedVersion.MatchString(version) {
		return fmt.Errorf("Expected Sub2API v0.2.16 or v0.2.17, received: %s", strings.TrimSpace(version))
	}
	if err := i.validateCandidate(ctx, cid, source); err != nil {
		return err
	}
	fmt.Fprintln(i.out, "[2/4] Back up the current executable.")
	directory, err := os.MkdirTemp("/tmp", "sub2api-codex-eval-"+shortID(cid)+"-")
	if err != nil {
		return err
	}
	backup := filepath.Join(directory, "sub2api")
	if _, err := i.docker(ctx, false, "cp", cid+":"+binaryPath, backup); err != nil {
		return err
	}
	if err := checkELF(backup); err != nil {
		return errors.New("The existing executable must also be Linux/amd64.")
	}
	fmt.Fprintf(i.out, "Backup: %s\n", backup)
	fmt.Fprintln(i.out, "[3/4] Replace the binary and restart the selected container.")
	err = i.replace(ctx, cid, source)
	if err == nil {
		_, err = i.docker(ctx, false, "restart", cid)
	}
	if err == nil {
		fmt.Fprintln(i.out, "[4/4] Wait for native health and setup status.")
		err = i.readiness(ctx, cid)
	}
	if err != nil {
		fmt.Fprintln(i.errOut, "Installation failed; restoring the original executable.")
		// Rollback uses a separate context so Ctrl-C cannot cancel restoration.
		restoreCtx := context.Background()
		restoreErr := i.replace(restoreCtx, cid, backup)
		if restoreErr == nil {
			_, restoreErr = i.docker(restoreCtx, false, "restart", cid)
		}
		if restoreErr == nil {
			restoreErr = i.readiness(restoreCtx, cid)
		}
		if restoreErr != nil {
			return fmt.Errorf("Installation failed: %w. Restoration failed: %v. Backup: %s", err, restoreErr, backup)
		}
		fmt.Fprintln(i.errOut, "Original executable restored.")
		return err
	}
	fmt.Fprintln(i.out, "Installed. Open 账户评估 in the native administrator sidebar.")
	return nil
}

func run(ctx context.Context, in io.Reader, out, errOut io.Writer) error {
	source, cleanup, err := extractPayload()
	if err != nil {
		return err
	}
	defer cleanup()
	i := installer{in: in, out: out, errOut: errOut, timeout: readinessTimeout, poll: time.Second}
	return i.install(ctx, source)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := run(ctx, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
