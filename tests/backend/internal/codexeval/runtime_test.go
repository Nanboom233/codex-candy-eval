package codexeval

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func waitIdle(t *testing.T) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		idle := len(active) == 0
		mu.Unlock()
		if idle {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("workers did not exit")
}
func TestCodexEvalRuntime(t *testing.T) {
	mu.Lock()
	loaded = false
	stateLoadError = nil
	storageError = ""
	candyResults = map[string][]candyResult{}
	fingerprintResults = map[string][]fingerprintResult{}
	pelicanResults = map[string][]pelicanResult{}
	mu.Unlock()
	oldDelay := retryDelay
	retryDelay = time.Millisecond
	t.Cleanup(func() { retryDelay = oldDelay })
	started := make(chan struct{})
	block := false
	pelicanAnswer := "<html><svg></svg></html>"
	var fingerprintCalls atomic.Int32
	send := func(ctx context.Context, id int64, body []byte, headers http.Header) (*http.Response, error) {
		if id != 7 {
			t.Errorf("sent to unexpected account %d", id)
		}
		if block {
			select {
			case started <- struct{}{}:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
			<-ctx.Done()
			return nil, ctx.Err()
		}
		var req map[string]any
		_ = json.Unmarshal(body, &req)
		answer := "21"
		if input, ok := req["input"].([]any); ok {
			last := input[len(input)-1].(map[string]any)["content"].([]any)[0].(map[string]any)["text"].(string)
			if strings.Contains(last, "HTML") {
				answer = pelicanAnswer
			} else if strings.Contains(last, "355") {
				answer = strings.Repeat("42,", 320)
				if fingerprintCalls.Add(1) == 1 {
					answer = "cannot answer"
				}
			}
		}
		data, _ := json.Marshal(map[string]any{"type": "response.completed", "response": map[string]any{"status": "completed", "output": []any{map[string]any{"type": "message", "content": []any{map[string]any{"type": "output_text", "text": answer}}}}, "usage": map[string]int{"input_tokens": 2, "output_tokens": 3}}})
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("data: " + string(data) + "\n\n"))}, nil
	}
	Initialize(t.TempDir(), send, func(_ context.Context, id int64) (string, error) {
		if id != 7 {
			return "", fmt.Errorf("missing account %d", id)
		}
		return "selected", nil
	})
	for _, kind := range []string{"candy", "fingerprint", "pelican"} {
		request := StartRequest{Type: kind, Targets: []Target{{1, 7}, {1, 7}, {1, 8}}, Model: "fixture", Effort: "none", Runs: 2}
		summary, err := Start(context.Background(), request)
		if err != nil {
			t.Fatal(err)
		}
		if len(summary.Started) != 1 || len(summary.Unavailable) != 1 || len(summary.Busy) != 0 {
			t.Fatalf("bad admission: %+v", summary)
		}
		waitIdle(t)
	}
	mu.Lock()
	if len(candyResults["7"]) != 2 || !candyResults["7"][0].OK || candyResults["7"][0].GroupID != 1 || candyResults["7"][0].Effort != "none" {
		t.Errorf("candy: %+v", candyResults)
	}
	if len(fingerprintResults["7"]) != 1 || len(fingerprintResults["7"][0].Samples) != 3 || fingerprintResults["7"][0].Status != "completed" || fingerprintCalls.Load() != 4 || fingerprintResults["7"][0].InputTokens != 8 {
		t.Errorf("fingerprint: %+v", fingerprintResults)
	}
	if len(pelicanResults["7"]) != 1 || pelicanResults["7"][0].HTML == "" {
		t.Errorf("pelican: %+v", pelicanResults)
	}
	mu.Unlock()
	fp := fingerprintResults["7"][0]
	if fp.summary().Samples != nil || len(fp.Samples) != 3 {
		t.Fatal("summary changed stored samples")
	}
	if detail, status := Record("fingerprint", 7, fp.ID); status != 200 || !strings.Contains(string(detail), `"samples":[`) {
		t.Fatalf("fingerprint detail: %s", detail)
	}
	for _, answer := range []string{"No HTML in this answer", "<html>" + strings.Repeat("x", pelicanMaxHTML) + "</html>"} {
		pelicanAnswer = answer
		result := evaluatePelican(credential{ID: "7", Target: Target{1, 7}, Name: "selected", ctx: context.Background()}, "fixture", "none")
		if result.Status != "failed" || result.Error == "" || result.HTML != "" {
			t.Fatalf("invalid HTML reported completed: %+v", result)
		}
	}
	block = true
	summary, err := Start(context.Background(), StartRequest{Type: "candy", Targets: []Target{{2, 7}}, Model: "fixture", Runs: 3})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("sender not started")
	}
	busy, err := Start(context.Background(), StartRequest{Type: "pelican", Targets: []Target{{1, 7}}, Model: "fixture"})
	if err != nil || len(busy.Busy) != 1 {
		t.Fatalf("cross-type busy failed: %+v %v", busy, err)
	}
	if !Cancel(summary.Started[0].ID) {
		t.Fatal("cancel not accepted")
	}
	waitIdle(t)
	mu.Lock()
	last := candyResults["7"][len(candyResults["7"])-1]
	mu.Unlock()
	if last.Status != "cancelled" || last.OK {
		t.Fatalf("cancel saved as success: %+v", last)
	}
	summary, err = Start(context.Background(), StartRequest{Type: "fingerprint", Targets: []Target{{2, 7}}, Model: "fixture", Concurrency: 1})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("fingerprint sender not started")
	}
	if !Cancel(summary.Started[0].ID) {
		t.Fatal("fingerprint cancel not accepted")
	}
	waitIdle(t)
	if latest := fingerprintResults["7"][1]; latest.Status != "cancelled" || len(latest.Samples) != 1 {
		t.Fatalf("fingerprint cancellation: %+v", latest)
	}
	// Exercise the same original atomic history path, including failed clear rollback.
	mu.Lock()
	for range 25 {
		candyResults["7"] = appendRecent(candyResults["7"], last)
	}
	err = saveStateLocked()
	mu.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	var saved persistedState
	raw, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &saved); err != nil || len(saved.Results["7"]) != 20 {
		t.Fatalf("history retention: %v", err)
	}
	oldPath := statePath
	statePath = filepath.Join(t.TempDir(), "missing", "history.json")
	_, status := Clear("candy")
	if status != 500 || len(candyResults["7"]) != 20 {
		t.Fatal("clear failure discarded history")
	}
	statePath = oldPath
	legacy := `{"results":{"7":[{"account_id":7,"group_id":1,"run_id":"old-candy","model":"preserved","ok":true,"answer":"21"}]},"pelicans":{"7":[{"id":"old-pelican","account_id":7,"group_id":1,"run_id":"old-run","model":"preserved","html":"<html>old</html>"}]},"fingerprints":{"7":[{"attribution":"retired"}]},"modeltraces":{"7":[{"id":"retired"}]}}`
	if err = os.WriteFile(statePath, []byte(legacy), 0600); err != nil {
		t.Fatal(err)
	}
	loaded = false
	fingerprintResults = map[string][]fingerprintResult{}
	loadState()
	raw, err = os.ReadFile(statePath)
	saved = persistedState{}
	if err != nil || json.Unmarshal(raw, &saved) != nil {
		t.Fatal("cannot read migrated history")
	}
	if saved.SchemaVersion != 1 || len(saved.Fingerprints) != 0 || strings.Contains(string(raw), "retired") || saved.Results["7"][0].RunID != "old-candy" || saved.Results["7"][0].Status != "" || saved.Results["7"][0].Attempts != 0 || saved.Pelicans["7"][0].HTML != "<html>old</html>" {
		t.Fatalf("minimal migration changed retained records: %s", raw)
	}
	loaded = false
	candyResults = map[string][]candyResult{}
	pelicanResults = map[string][]pelicanResult{}
	loadState()
	if len(candyResults["7"]) != 1 || len(pelicanResults["7"]) != 1 {
		t.Fatal("migrated records did not reload")
	}
	invalidMarker := `{"schema_version":null,"results":{"7":[]}}`
	if err = os.WriteFile(statePath, []byte(invalidMarker), 0600); err != nil {
		t.Fatal(err)
	}
	loaded = false
	loadState()
	if len(candyResults["7"]) != 1 || storageError == "" {
		t.Fatal("invalid marker replaced loaded history")
	}
	_, status = Clear("all")
	raw, _ = os.ReadFile(statePath)
	if status != 500 || string(raw) != invalidMarker {
		t.Fatal("invalid marker was migrated or overwritten")
	}
	stateLoadError = nil
	if err = os.WriteFile(statePath, []byte("null"), 0600); err != nil {
		t.Fatal(err)
	}
	loaded = false
	loadState()
	_, status = Clear("all")
	if status != 500 {
		t.Fatal("corrupt history was overwritten")
	}
	raw, _ = os.ReadFile(statePath)
	if string(raw) != "null" {
		t.Fatal("original corrupt file changed")
	}
}
func TestCodexEvalStreamSemantics(t *testing.T) {
	text := "event: usage\ndata: {\"type\":\"response.in_progress\",\"response\":{\"usage\":{\"input_tokens\":7}}}\n\n" +
		"data: {\"type\":\"response.output_item.done\",\"item\":{\"type\":\"message\",\"content\":[{\"type\":\"output_text\",\"text\":\"21\"}]}}\n\n" +
		"data: {\"type\":\"response.done\",\"response\":{\"status\":\"completed\",\"output\":[],\"usage\":{\"output_tokens\":3}}}"
	out, err := readResponseStream(strings.NewReader(text))
	if err != nil || out.Usage.InputTokens != 7 || out.Usage.OutputTokens != 3 || len(out.Output) != 1 {
		t.Fatalf("stream: %+v %v", out, err)
	}
	out, err = readResponseStream(strings.NewReader("data: {\"type\":\"response.in_progress\",\"response\":{\"usage\":{\"input_tokens\":9}}}\n\ndata: {\"type\":\"error\",\"message\":\"actual failure\"}\n\n"))
	if err == nil || !strings.Contains(err.Error(), "actual failure") || out.Usage.InputTokens != 9 {
		t.Fatalf("failed usage: %+v %v", out, err)
	}
}
