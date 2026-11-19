// sub2api-package: backend/internal/codexeval
package codexeval

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestFingerprintUpstreamParity(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(os.Getenv("SUB2API_ASSETS"), "lm-detector-golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Name    string
		Outputs []struct {
			Text          string
			ExpectedCount int `json:"expected_count"`
		}
		Expected fingerprintAttribution
	}
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			replies := make([]fingerprintReply, len(tc.Outputs))
			for i, output := range tc.Outputs {
				replies[i] = prepareFingerprint(output.Text, output.ExpectedCount)
			}
			got, want := analyzeFingerprint(replies), tc.Expected
			if got == nil || got.Decision != want.Decision || got.Prediction != want.Prediction ||
				got.PredictionName != want.PredictionName || got.FamilyPrediction != want.FamilyPrediction ||
				got.Used != want.Used || !slices.Equal(got.Diagnostics, want.Diagnostics) || len(got.Results) != len(want.Results) {
				t.Fatalf("attribution mismatch: got %+v, want %q with diagnostics %v", got, want.Prediction, want.Diagnostics)
			}
			if got.Sequence == nil || got.Sequence.ModelLike != want.Sequence.ModelLike {
				t.Fatalf("sequence = %+v, want %+v", got.Sequence, want.Sequence)
			}
			checkNumber(t, "lag-1 autocorrelation", &got.Sequence.Autocorrelation, &want.Sequence.Autocorrelation)
			checkNumber(t, "small-step share", &got.Sequence.SmallStepShare, &want.Sequence.SmallStepShare)
			checkNumber(t, "top probability", got.Probability, want.Probability)
			for i, candidate := range got.Results {
				expected := want.Results[i]
				checkNumber(t, candidate.Model+" probability", candidate.Probability, expected.Probability)
				checkNumber(t, candidate.Model+" score", &candidate.Score, &expected.Score)
				candidate.Probability, expected.Probability = nil, nil
				candidate.Score, expected.Score = 0, 0
				if candidate != expected {
					t.Fatalf("rank %d = %+v, want %+v", i+1, candidate, expected)
				}
			}
		})
	}
}

func TestFingerprintThreeLagAcceptance(t *testing.T) {
	numbers := make([]int, 300)
	pattern := []int{1, 81, 161, 241, 321, 1, 81, 161, 241, 321, 321}
	for i := range numbers {
		numbers[i] = pattern[i%len(pattern)]
	}
	shape := sequenceShape([][]int{numbers})
	if shape.Autocorrelation <= -0.23 || shape.SmallStepShare <= 0.05 || !shape.ModelLike {
		t.Fatalf("three-lag acceptance = %+v", shape)
	}
}

func TestFingerprintParticipationStopsWithoutHistory(t *testing.T) {
	oldDelay := retryDelay
	retryDelay = time.Millisecond
	t.Cleanup(func() { retryDelay = oldDelay })
	for _, sentBeforeStop := range []int{0, 1} {
		resetInspectionTestState(t)
		calls := 0
		Initialize(t.TempDir(), func(context.Context, int64, []byte, http.Header) (*http.Response, error) {
			calls++
			return &http.Response{StatusCode: 503, Body: io.NopCloser(strings.NewReader("temporary"))}, nil
		}, nil)
		p := &fingerprintProgress{Total: 3, Concurrency: 1}
		fingerprintRunning["7"] = p
		auth := credential{ID: "7", Target: Target{AccountID: 7}, ctx: context.Background()}
		runFingerprint(auth.ctx, auth, "fixture", p, func() error {
			if calls >= sentBeforeStop {
				return errInspectionStopped
			}
			return nil
		})
		if calls != sentBeforeStop || len(fingerprintResults["7"]) != 0 || fingerprintRunning["7"] != nil || p.Done != 0 {
			t.Fatalf("stopped run: calls=%d progress=%+v history=%+v", calls, p, fingerprintResults["7"])
		}
	}
}

func checkNumber(t *testing.T, label string, got, want *float64) {
	t.Helper()
	if (got == nil) != (want == nil) {
		t.Fatalf("%s: got probability=%t, want probability=%t", label, got != nil, want != nil)
	}
	if got != nil && !(math.Abs(*got-*want) <= 1e-9) {
		t.Fatalf("%s = %.12g, want %.12g", label, *got, *want)
	}
}

func TestFingerprintReplyPreparation(t *testing.T) {
	if got := parseNumbers("生成 300 个：1, 22, 355, 0, 356, 99。完成 4"); !slices.Equal(got, []int{1, 22, 355, 99}) {
		t.Fatalf("parsed %v", got)
	}
	if prepareFingerprint(strings.Repeat("42,", 164), 300).Valid() || !prepareFingerprint(strings.Repeat("42,", 165), 300).Valid() {
		t.Fatal("expected a minimum of 165 numbers for a challenge of 300")
	}
}

func TestFingerprintParseAndChallenges(t *testing.T) {
	for text, want := range map[string][]int{
		"生成 300 个：1, 22, 355, 0, 356, 99。完成 4": {1, 22, 355, 99},
		"1,2 words 3,4": {1, 2}, "1,2 中文 3,4,5": {3, 4, 5}, "1,2 😀 3,4": {1, 2, 3, 4},
	} {
		if got := parseNumbers(text); !reflect.DeepEqual(got, want) {
			t.Fatalf("parse %q = %v, want %v", text, got, want)
		}
	}
	instructed, prefixed := 0, 0
	for _, c := range fingerprintSuite {
		if c.Instructions != "" {
			instructed++
		}
		if strings.Contains(c.Prompt, "\n\nFinal task:\n") {
			prefixed++
		}
	}
	if len(fingerprintSuite) != 36 || instructed != 12 || prefixed != 9 {
		t.Fatalf("suite: %d challenges, %d with instructions, %d with a user prefix", len(fingerprintSuite), instructed, prefixed)
	}
	seen := map[string]bool{}
	for _, c := range fingerprintChallenges() {
		if seen[c.Prompt] || !slices.Contains(fingerprintSuite, c) {
			t.Fatalf("bad challenge: %+v", c)
		}
		seen[c.Prompt] = true
	}
	if FingerprintModelID("openai/GPT-6.1-sol") != "gpt-6-astra" || FingerprintModelID("deepseek/deepseek-flash") != "deepseek-v4.1-flash-1017" || FingerprintModelID("gpt-6-sol") != "gpt-6-sol" {
		t.Fatal("aliased models must compare as their candidates")
	}
}

func TestFingerprintRevisionAndPreparedReply(t *testing.T) {
	if want := fmt.Sprintf("sha256:%x", sha256.Sum256(detectorFile)); detectorRevision != want {
		t.Fatalf("revision = %s, want %s", detectorRevision, want)
	}
	if got := analyzeFingerprint([]fingerprintReply{prepareFingerprint("1, 2", 300)}); got != nil {
		t.Fatalf("short output was accepted: %+v", got)
	}
	reply := prepareFingerprint(strings.Repeat("42,", 301), 300)
	if !reply.Valid() || reply.Count() != 300 || strings.HasSuffix(reply.Text, ",") {
		t.Fatalf("reply truncation = %+v", reply)
	}
	if got := analyzeFingerprint([]fingerprintReply{reply}); got == nil || got.Used != 1 || got.Decision != "not_model" || got.Results != nil {
		t.Fatalf("repeated numbers = %+v", got)
	}
	if len(detectorFile) != 5232052 || sharedDetector.BankBuiltAt != "2026-11-17T19:37:59.839365+00:00" {
		t.Fatal("detector assets do not match v0.4.7")
	}
}
