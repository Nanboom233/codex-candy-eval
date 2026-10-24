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
const stagedPath = "/app/sub2api.plugins-new"
const readinessTimeout = 120 * time.Second

//go:embed payload.gz licenses
var assets embed.FS

var elfHeader = []byte{0x7f, 'E', 'L', 'F', 2, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0x3e, 0}
var buildVersion string
var versionRecord = regexp.MustCompile(`^Sub2API (\S+) \(commit: [^,\r\n]+, built: [^)\r\n]+\)$`)
var versionToken = regexp.MustCompile(`^v?[0-9]+\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)
var standardLogPrefix = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}(?:\.\d+)? `)

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
	in           io.Reader
	out, errOut  io.Writer
	timeout      time.Duration
	poll         time.Duration
	reader       *bufio.Reader
	probeTimeout time.Duration
}

func (i *installer) readLine(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if i.reader == nil {
		i.reader = bufio.NewReader(i.in)
	}
	type input struct {
		line string
		err  error
	}
	result := make(chan input, 1)
	go func() {
		line, err := i.reader.ReadString('\n')
		result <- input{line, err}
	}()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case value := <-result:
		if err := ctx.Err(); err != nil {
			return "", err
		}
		return strings.TrimSpace(value.line), value.err
	}
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
			line, err := i.readLine(ctx)
			if err != nil {
				return container{}, err
			}
			n, err := strconv.Atoi(line)
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

type versionResult struct {
	version string
	err     error
}

func parseVersion(stdout, stderr string) (string, error) {
	version := ""
	for _, stream := range []string{stdout, stderr} {
		for _, line := range strings.Split(stream, "\n") {
			message := strings.TrimSpace(line)
			fields := strings.Split(message, "\t")
			if len(fields) >= 4 && fields[1] == "INFO" && fields[2] == "stdlog" {
				message = fields[3]
			} else {
				message = standardLogPrefix.ReplaceAllString(message, "")
			}
			record := versionRecord.FindStringSubmatch(message)
			if record == nil {
				continue
			}
			if !versionToken.MatchString(record[1]) {
				return "", fmt.Errorf("Unresolved executable version: %s", record[1])
			}
			current := strings.TrimPrefix(record[1], "v")
			if version != "" && current != version {
				return "", errors.New("Conflicting executable version records")
			}
			version = current
		}
	}
	if version == "" {
		return "", errors.New("No executable version record")
	}
	return version, nil
}

func (i *installer) probeVersion(ctx context.Context, cid, path string) versionResult {
	timeout := i.probeTimeout
	if timeout == 0 {
		timeout = 15 * time.Second
	}
	probeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	command := exec.CommandContext(probeCtx, "docker", "exec", cid, path, "-version")
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	err := command.Run()
	if probeCtx.Err() != nil {
		err = probeCtx.Err()
	}
	if err != nil {
		return versionResult{err: fmt.Errorf("Version command failed: %w. %s", err, strings.TrimSpace(stdout.String()+"\n"+stderr.String()))}
	}
	version, err := parseVersion(stdout.String(), stderr.String())
	return versionResult{version: version, err: err}
}

func (result versionResult) String() string {
	if result.err != nil {
		return result.err.Error()
	}
	return result.version
}

func (i *installer) confirm(ctx context.Context) error {
	fmt.Fprint(i.out, "Continue installation? [y/N] ")
	line, err := i.readLine(ctx)
	if err != nil {
		return fmt.Errorf("Installation canceled: %w", err)
	}
	switch strings.ToLower(line) {
	case "y", "yes":
		return nil
	default:
		return errors.New("Installation canceled.")
	}
}

func (i *installer) confirmVersions(ctx context.Context, installed, candidate versionResult) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if installed.err == nil && candidate.err == nil && installed.version == buildVersion && candidate.version == buildVersion {
		return nil
	}
	fmt.Fprintf(i.errOut, "Version warning. Build: %s\nInstalled: %s\nCandidate: %s\n", buildVersion, installed, candidate)
	return i.confirm(ctx)
}

func (i *installer) validateCandidate(ctx context.Context, cid, source string) (version versionResult, result error) {
	if err := checkELF(source); err != nil {
		return version, err
	}
	if _, err := i.docker(ctx, false, "cp", source, cid+":"+stagedPath); err != nil {
		return version, err
	}
	defer func() {
		// An interrupt must not skip removal of the staged candidate.
		_, err := i.docker(context.Background(), false, "exec", "-u", "0", cid, "rm", "-f", stagedPath)
		result = errors.Join(result, err)
	}()
	if _, err := i.docker(ctx, false, "exec", "-u", "0", cid, "chmod", "755", stagedPath); err != nil {
		return version, err
	}
	return i.probeVersion(ctx, cid, stagedPath), nil
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
	fmt.Fprintln(i.out, "[1/4] Check executable versions and the local binary.")
	installed := i.probeVersion(ctx, cid, binaryPath)
	if err := ctx.Err(); err != nil {
		return err
	}
	candidate, err := i.validateCandidate(ctx, cid, source)
	if err != nil {
		return err
	}
	if err := i.confirmVersions(ctx, installed, candidate); err != nil {
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
	fmt.Fprintln(i.out, "Installed. The application is ready.")
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
