package codexeval

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestCodexInspectionConfirmationAndBoundaries(t *testing.T) {
	mu.Lock()
	loaded = false
	stateLoadError = nil
	storageError = ""
	candyResults = map[string][]candyResult{}
	mu.Unlock()
	var lock sync.Mutex
	calls := map[int64]int{}
	updates := map[int64]int{}
	accounts := []InspectionAccount{
		{Target: Target{AccountID: 1}, Name: "codex-wrong", Available: true, Priority: 10},
		{Target: Target{AccountID: 2}, Name: "codex-correct", Available: true, Priority: 10},
		{Target: Target{AccountID: 3}, Name: "codex-mixed", Available: true, Priority: 10},
		{Target: Target{AccountID: 4}, Name: "codex-disabled", Available: false},
		{Target: Target{AccountID: 5}, Name: "codex-free", Available: true, PlanType: "free"},
		{Target: Target{AccountID: 6}, Name: "outside-old-filter", Available: true},
		{Target: Target{AccountID: 7}, Name: "codex-recovery", Available: true, Priority: 99},
		{Target: Target{AccountID: 8}, Name: "codex-failed", Available: true, Priority: 10},
	}
	for n := range accounts {
		accounts[n].GroupIDs = []int64{101}
		accounts[n].Type = "oauth"
		accounts[n].Enabled = true
		accounts[n].AutoInspectionEnabled = true
		if accounts[n].PlanType == "" {
			accounts[n].PlanType = "pro"
		}
	}
	Initialize(t.TempDir(), func(ctx context.Context, id int64, _ []byte, _ http.Header) (*http.Response, error) {
		lock.Lock()
		calls[id]++
		n := calls[id]
		lock.Unlock()
		if id == 8 {
			return nil, fmt.Errorf("fake upstream failed")
		}
		answer := "20"
		if id == 2 || id == 7 || (id == 3 && n == 2) {
			answer = "21"
		}
		body := fmt.Sprintf("data: {\"type\":\"response.completed\",\"response\":{\"status\":\"completed\",\"output\":[{\"type\":\"message\",\"content\":[{\"type\":\"output_text\",\"text\":%q}]}]}}\n\n", answer)
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
	}, func(context.Context, int64) (string, error) { return "fake", nil })
	notified := 0
	hooks := InspectionHooks{
		List: func(_ context.Context, s InspectionSettings) ([]InspectionAccount, error) {
			if s.GroupID != 101 {
				t.Errorf("wrong selected group: %d", s.GroupID)
			}
			return accounts, nil
		},
		Get: func(_ context.Context, id int64) (InspectionAccount, error) {
			for _, a := range accounts {
				if a.AccountID == id {
					return a, nil
				}
			}
			panic("unexpected account")
		},
		UpdatePriority: func(_ context.Context, id int64, priority int) error { updates[id] = priority; return nil },
		Notify: func(context.Context, InspectionSettings, string, string) error {
			notified++
			return fmt.Errorf("notification failed")
		},
	}
	i := newInspector(hooks)
	defer i.cancel()
	i.poll = time.Millisecond
	i.recovery = time.Millisecond
	s := DefaultInspectionSettings()
	s.GroupID, s.AccountPlan = 101, "Pro ALL"
	s.AdjustPriority = true
	i.configure(s)
	if err := i.start(); err != nil {
		t.Fatal(err)
	}
	<-i.runDone
	if calls[1] != 3 || calls[2] != 1 || calls[3] != 3 || calls[7] != 3 || calls[8] != 1 || calls[4] != 0 || calls[5] != 0 || calls[6] != 1 {
		t.Fatalf("wrong test counts: %v", calls)
	}
	if len(updates) != 2 || updates[1] != 99 || updates[7] != 10 {
		t.Fatalf("wrong priority updates: %v", updates)
	}
	if notified != 2 {
		t.Fatalf("notifications=%d", notified)
	}
	state := i.snapshot()
	if state.Running || state.Done != 6 {
		t.Fatalf("unexpected state: %+v", state)
	}
	for _, r := range state.Results {
		if r.Outcome == "updated" && r.NotificationError == "" {
			t.Fatal("missing notification failure")
		}
	}

	// Default recommendation mode does not write priorities or send notifications.
	clear(updates)
	notified = 0
	s.AdjustPriority = false
	i.configure(s)
	if err := i.start(); err != nil {
		t.Fatal(err)
	}
	<-i.runDone
	if len(updates) != 0 || notified != 0 {
		t.Fatal("recommendations changed accounts")
	}
}

func TestCodexInspectionInvalidInitialSettingsBlockStart(t *testing.T) {
	settings := DefaultInspectionSettings()
	settings.Enabled = true
	settings.DegradedPriority = 5
	i := newInspector(InspectionHooks{
		Load: func(context.Context) (InspectionSettings, error) { return settings, nil },
	})
	defer i.cancel()

	i.loadInitialSettings()
	state := i.snapshot()
	if state.Error == "" || state.NextRunAt != nil || i.settingsLoaded {
		t.Fatalf("invalid settings were configured or scheduled: %+v", state)
	}
	if err := i.start(); err == nil {
		t.Fatal("manual inspection started with invalid loaded settings")
	}
}

func TestCodexInspectionCancelAndSchedule(t *testing.T) {
	mu.Lock()
	loaded = false
	stateLoadError = nil
	storageError = ""
	candyResults = map[string][]candyResult{}
	mu.Unlock()
	started := make(chan struct{}, 1)
	Initialize(t.TempDir(), func(ctx context.Context, _ int64, _ []byte, _ http.Header) (*http.Response, error) {
		started <- struct{}{}
		<-ctx.Done()
		return nil, ctx.Err()
	}, func(context.Context, int64) (string, error) { return "fake", nil })
	a := InspectionAccount{Target: Target{AccountID: 42}, GroupIDs: []int64{101}, Type: "oauth", Available: true, Enabled: true, AutoInspectionEnabled: true}
	i := newInspector(InspectionHooks{List: func(context.Context, InspectionSettings) ([]InspectionAccount, error) {
		return []InspectionAccount{a}, nil
	}, Get: func(context.Context, int64) (InspectionAccount, error) { return a, nil }, UpdatePriority: func(context.Context, int64, int) error { t.Error("cancelled run wrote account"); return nil }})
	defer i.cancel()
	i.poll = time.Millisecond
	s := DefaultInspectionSettings()
	s.Enabled = true
	s.GroupID = 101
	s.AdjustPriority = true
	i.configure(s)
	before := i.snapshot().NextRunAt
	if before == nil || time.Until(*before) < 15*time.Minute-time.Second || time.Until(*before) > 35*time.Minute {
		t.Fatal("schedule outside bounds")
	}
	if i.snapshot().NextRunAt != before {
		t.Fatal("state read reset schedule")
	}
	s.NtfyTopic = "changed"
	i.configure(s)
	if i.snapshot().NextRunAt != before {
		t.Fatal("autosave reset the pending schedule")
	}
	if err := i.start(); err != nil {
		t.Fatal(err)
	}
	<-started
	if err := i.start(); err == nil {
		t.Fatal("overlapping run accepted")
	}
	s.Enabled = false
	i.configure(s)
	if !i.snapshot().Running {
		t.Fatal("disable cancelled active round")
	}
	i.stop()
	<-i.runDone
	state := i.snapshot()
	if state.Running || state.Phase != "cancelled" || state.NextRunAt != nil {
		t.Fatalf("wrong stop state: %+v", state)
	}
	s.MinWaitMinutes = 17
	s.MaxWaitMinutes = 17
	if inspectionDelay(s) != 17*time.Minute {
		t.Fatal("equal bounds not fixed")
	}
}

func TestCodexInspectionNotifications(t *testing.T) {
	var ntfy, bark bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/topic" {
			user, password, ok := r.BasicAuth()
			if !ok || user != "fixture" || password != "secret" {
				t.Error("missing ntfy auth")
			}
			ntfy = true
			w.WriteHeader(503)
			return
		}
		if r.URL.Path == "/private-bark-key" {
			bark = true
			_ = r.ParseForm()
			if r.Form.Get("body") != "message" {
				t.Error("missing Bark body")
			}
		}
	}))
	defer server.Close()
	s := DefaultInspectionSettings()
	s.NtfyServer = server.URL
	s.NtfyTopic = "topic"
	s.NtfyUsername = "fixture"
	s.NtfyPassword = "secret"
	s.BarkURL = server.URL + "/private-bark-key"
	err := SendInspectionNotification(context.Background(), s, "title", "message")
	if err == nil || !strings.Contains(err.Error(), "ntfy: HTTP 503") || strings.Contains(err.Error(), "private-bark-key") || !ntfy || !bark {
		t.Fatalf("unexpected notification result: %v", err)
	}
	public := s.Public()
	if public.NtfyPassword != "" || public.BarkURL != "" || !public.NtfyPasswordConfigured || !public.BarkURLConfigured {
		t.Fatal("secret exposed")
	}
}

func TestCodexInspectionRestartSchedulesFreshWait(t *testing.T) {
	s := DefaultInspectionSettings()
	s.Enabled = true
	s.GroupID = 101
	loaded := make(chan struct{})
	i := newInspector(InspectionHooks{Load: func(context.Context) (InspectionSettings, error) { close(loaded); return s, nil }})
	i.loadInitialSettings()
	go i.loop()
	defer func() { i.cancel(); <-i.done }()
	<-loaded
	state := i.snapshot()
	if state.NextRunAt == nil || state.Running || state.StartedAt != nil || len(state.Results) != 0 {
		t.Fatalf("restart resumed work or missed schedule: %+v", state)
	}
	first := *state.NextRunAt
	s.MinWaitMinutes = 16
	i.configure(s)
	if next := i.snapshot().NextRunAt; next == nil || next.Equal(first) {
		t.Fatal("configuration did not draw a fresh wait")
	}
	s.Enabled = false
	i.configure(s)
	if err := i.startRun(true); err != nil || i.snapshot().Running {
		t.Fatal("disabled stale timer started a round")
	}
}

func TestCodexInspectionCandySelection(t *testing.T) {
	s := DefaultInspectionSettings()
	s.GroupID = 202
	a := InspectionAccount{Target: Target{AccountID: 33}, GroupIDs: []int64{101, 202}, Type: "oauth", PlanType: "pro_lite", Available: true, Enabled: true, AutoInspectionEnabled: true}
	if !inspectionEligible(a, s) {
		t.Fatal("shared group member was excluded")
	}
	s.AccountPlan = "Pro ALL"
	if !inspectionEligible(a, s) {
		t.Fatal("Pro ALL lost a Pro alias")
	}
	s.AccountPlan = "Pro 100"
	if !inspectionEligible(a, s) {
		t.Fatal("exact plan lost a Pro alias")
	}
	s.AccountPlan = "Free"
	if inspectionEligible(a, s) {
		t.Fatal("a different plan matched")
	}
	a.PlanType = "free"
	if !inspectionEligible(a, s) {
		t.Fatal("the Candy Free filter was overridden")
	}
	a.Type = "apikey"
	if inspectionEligible(a, s) {
		t.Fatal("API keys matched an OAuth plan")
	}
	s.AccountPlan = "all"
	s.AccountType = "oauth"
	if inspectionEligible(a, s) {
		t.Fatal("account type filter was ignored")
	}
	s.AccountType = "all"
	a.GroupIDs = []int64{101}
	if inspectionEligible(a, s) {
		t.Fatal("an account outside the saved group matched")
	}
	a.GroupIDs = []int64{202}
	a.Available = false
	if inspectionEligible(a, s) {
		t.Fatal("an unavailable account matched")
	}
}

func TestCodexInspectionDegradedWrongOnce(t *testing.T) {
	for _, adjust := range []bool{false, true} {
		for _, listedPriority := range []int{99, 10} {
			t.Run(fmt.Sprintf("adjust=%t/list=%d", adjust, listedPriority), func(t *testing.T) {
				mu.Lock()
				loaded, stateLoadError, storageError = false, nil, ""
				candyResults = map[string][]candyResult{}
				mu.Unlock()
				calls := 0
				Initialize(t.TempDir(), func(context.Context, int64, []byte, http.Header) (*http.Response, error) {
					calls++
					return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("data: {\"type\":\"response.completed\",\"response\":{\"status\":\"completed\",\"output\":[{\"type\":\"message\",\"content\":[{\"type\":\"output_text\",\"text\":\"20\"}]}]}}\n\n"))}, nil
				}, func(context.Context, int64) (string, error) { return "degraded", nil })
				listed := InspectionAccount{Target: Target{AccountID: 91}, GroupIDs: []int64{101}, Type: "oauth", Available: true, Enabled: true, AutoInspectionEnabled: true, Priority: listedPriority}
				current := listed
				current.Priority = 99
				i := newInspector(InspectionHooks{
					List: func(context.Context, InspectionSettings) ([]InspectionAccount, error) {
						return []InspectionAccount{listed}, nil
					},
					Get: func(context.Context, int64) (InspectionAccount, error) { return current, nil },
					UpdatePriority: func(context.Context, int64, int) error {
						t.Error("degraded wrong account priority changed")
						return nil
					},
					Notify: func(context.Context, InspectionSettings, string, string) error {
						t.Error("degraded wrong account notified")
						return nil
					},
				})
				defer i.cancel()
				i.poll = time.Millisecond
				s := DefaultInspectionSettings()
				s.GroupID, s.AdjustPriority = 101, adjust
				i.configure(s)
				if err := i.start(); err != nil {
					t.Fatal(err)
				}
				<-i.runDone
				r := i.snapshot().Results[0]
				if calls != 1 || r.Tests != 1 || len(candyResults["91"]) != 1 || r.Outcome != "unchanged" || r.TargetPriority != nil || r.Error != "" || current.Priority != 99 {
					t.Fatalf("calls=%d result=%+v priority=%d", calls, r, current.Priority)
				}
			})
		}
	}
}

func TestCodexInspectionParticipationStopsLaterSends(t *testing.T) {
	mu.Lock()
	loaded, stateLoadError, storageError = false, nil, ""
	candyResults = map[string][]candyResult{}
	mu.Unlock()
	var lock sync.Mutex
	a := InspectionAccount{Target: Target{GroupID: 101, AccountID: 92}, GroupIDs: []int64{101}, Type: "oauth", Available: true, Enabled: true, Priority: 10}
	calls := 0
	confirmation := make(chan struct{})
	release := make(chan struct{})
	Initialize(t.TempDir(), func(context.Context, int64, []byte, http.Header) (*http.Response, error) {
		calls++
		if calls == 2 {
			close(confirmation)
			<-release
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("data: {\"type\":\"response.completed\",\"response\":{\"status\":\"completed\",\"output\":[{\"type\":\"message\",\"content\":[{\"type\":\"output_text\",\"text\":\"20\"}]}]}}\n\n"))}, nil
	}, func(context.Context, int64) (string, error) { return "participation", nil })
	i := newInspector(InspectionHooks{
		List: func(context.Context, InspectionSettings) ([]InspectionAccount, error) {
			lock.Lock()
			defer lock.Unlock()
			return []InspectionAccount{a}, nil
		},
		Get: func(context.Context, int64) (InspectionAccount, error) {
			lock.Lock()
			defer lock.Unlock()
			return a, nil
		},
		UpdatePriority: func(context.Context, int64, int) error { t.Error("opted-out account priority changed"); return nil },
		Notify: func(context.Context, InspectionSettings, string, string) error {
			t.Error("opted-out account notified")
			return nil
		},
	})
	defer i.cancel()
	i.poll = time.Millisecond
	s := DefaultInspectionSettings()
	s.GroupID, s.AdjustPriority = 101, true
	i.configure(s)
	if err := i.start(); err != nil {
		t.Fatal(err)
	}
	<-i.runDone
	if calls != 0 || len(i.snapshot().Results) != 0 {
		t.Fatal("opted-out account selected")
	}
	lock.Lock()
	a.AutoInspectionEnabled = true
	lock.Unlock()
	if err := i.start(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-confirmation:
	case <-time.After(5 * time.Second):
		t.Fatal("confirmation did not start")
	}
	lock.Lock()
	a.AutoInspectionEnabled = false
	lock.Unlock()
	close(release)
	<-i.runDone
	r := i.snapshot().Results[0]
	if calls != 2 || r.Tests != 2 || r.Outcome != "unchanged" || r.Error != "" || r.TargetPriority != nil || len(candyResults["92"]) != 2 {
		t.Fatalf("opt-out calls=%d result=%+v", calls, r)
	}
	for _, result := range candyResults["92"] {
		if result.Error != "" {
			t.Fatal("opt-out created a test error")
		}
	}
	// Manual fixed-count tests remain usable for this opted-out account.
	_, err := Start(context.Background(), StartRequest{Type: "candy", Targets: []Target{a.Target}, Model: "fixture", Runs: 3})
	if err != nil {
		t.Fatal(err)
	}
	waitIdle(t)
	if calls != 5 || len(candyResults["92"]) != 5 {
		t.Fatalf("manual count changed: %d", calls)
	}
}
