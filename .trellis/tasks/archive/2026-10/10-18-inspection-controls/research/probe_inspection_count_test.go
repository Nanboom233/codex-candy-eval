package codexeval

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestResearchInspectionDegradedWrongCalls(t *testing.T) {
	for _, tc := range []struct {
		name         string
		listPriority int
		getPriority  int
		wantObserved int
	}{
		{"stable_degraded", 99, 99, 1},
		{"normal_snapshot_now_degraded", 10, 99, 3},
	} {
		for _, adjust := range []bool{false, true} {
			t.Run(tc.name+map[bool]string{false: "/recommend", true: "/adjust"}[adjust], func(t *testing.T) {
				mu.Lock()
				loaded = false
				stateLoadError = nil
				storageError = ""
				candyResults = map[string][]candyResult{}
				mu.Unlock()
				calls, updates, notifications := 0, 0, 0
				Initialize(t.TempDir(), func(context.Context, int64, []byte, http.Header) (*http.Response, error) {
					calls++
					body := "data: {\"type\":\"response.completed\",\"response\":{\"status\":\"completed\",\"output\":[{\"type\":\"message\",\"content\":[{\"type\":\"output_text\",\"text\":\"20\"}]}]}}\n\n"
					return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body))}, nil
				}, func(context.Context, int64) (string, error) { return "research account", nil })
				listed := InspectionAccount{Target: Target{AccountID: 41}, Name: "research account", Priority: tc.listPriority,
					Available: true, Enabled: true, GroupIDs: []int64{101}, Type: "oauth", PlanType: "pro"}
				current := listed
				current.Priority = tc.getPriority
				i := newInspector(InspectionHooks{
					List: func(context.Context, InspectionSettings) ([]InspectionAccount, error) {
						return []InspectionAccount{listed}, nil
					},
					Get:            func(context.Context, int64) (InspectionAccount, error) { return current, nil },
					UpdatePriority: func(context.Context, int64, int) error { updates++; return nil },
					Notify:         func(context.Context, InspectionSettings, string, string) error { notifications++; return nil },
				})
				defer i.cancel()
				i.poll = time.Millisecond
				s := DefaultInspectionSettings()
				s.Enabled, s.AdjustPriority, s.GroupID = true, adjust, 101
				i.configure(s)
				i.mu.Lock()
				due := time.Now().Add(-time.Second)
				i.state.NextRunAt = &due
				i.mu.Unlock()
				if err := i.startRun(true); err != nil {
					t.Fatal(err)
				}
				<-i.runDone
				state := i.snapshot()
				if calls != tc.wantObserved || len(state.Results) != 1 || state.Results[0].Tests != tc.wantObserved {
					t.Fatalf("unexpected observed calls=%d state=%+v", calls, state)
				}
				if updates != 0 || notifications != 0 {
					t.Fatalf("unexpected writes=%d notifications=%d", updates, notifications)
				}
				t.Logf("list_priority=%d current_priority=%d adjust=%v calls=%d tests=%d outcome=%s target=%v error=%q",
					tc.listPriority, tc.getPriority, adjust, calls, state.Results[0].Tests,
					state.Results[0].Outcome, state.Results[0].TargetPriority, state.Results[0].Error)
			})
		}
	}
}
