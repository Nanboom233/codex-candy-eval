package codexeval

import (
	"math"
	"testing"
)

func TestFingerprintNormalization(t *testing.T) {
	for _, tc := range []struct{ kind, raw, answer, category string }{
		{"int", "**４７**", "47", "valid"}, {"int", "٤٧", "47", "valid"}, {"int", "۴۷", "47", "valid"},
		{"int", "四十七。", "47", "valid"}, {"int", "一百", "100", "valid"}, {"int", "十", "10", "valid"},
		{"int", "forty-seven", "47", "valid"}, {"int", "one hundred", "1", "valid"},
		{"int", "101", "101", "invalid"}, {"int", "0", "0", "invalid"}, {"int", "A number: 47", "a", "invalid"},
		{"letter", "Zed!", "z", "valid"}, {"letter", "QUEUE", "q", "valid"}, {"letter", "abc", "abc", "invalid"},
		{"color", "Grey", "gray", "valid"}, {"color", "蓝色", "蓝", "valid"}, {"color", "Aqua", "cyan", "valid"},
		{"coin", "正面！", "heads", "valid"}, {"coin", "花", "tails", "valid"}, {"coin", "HEAD", "heads", "valid"},
		{"word", "Snow leopard", "snow", "valid"}, {"word", "“上海”", "上海", "valid"},
		{"word", "e\u0301", "é", "invalid"}, {"word", "🦊", "", "empty"}, {"int", "", "", "empty"},
		{"int", "As an AI, I cannot pick", "", "refusal"}, {"word", "抱歉，不能回答", "", "refusal"},
	} {
		t.Run(tc.kind+"/"+tc.raw, func(t *testing.T) {
			a, c := normalizeFingerprintAnswer(tc.raw, fingerprintProbe{Kind: tc.kind, Lo: 1, Hi: 100})
			if a != tc.answer || c != tc.category {
				t.Fatalf("got (%q, %q), want (%q, %q)", a, c, tc.answer, tc.category)
			}
		})
	}
}

func repeatedFingerprint(answer string, n int) map[string][]string {
	cells := map[string][]string{}
	for _, probe := range fingerprintProbes[:4] {
		for range n {
			cells[probe.ID] = append(cells[probe.ID], answer)
		}
	}
	return cells
}

func TestFingerprintStatistics(t *testing.T) {
	if got := fingerprintJSD([]string{"a", "a"}, []string{"b"}); got != 1 {
		t.Fatalf("disjoint JSD = %v", got)
	}
	if got := fingerprintJSD([]string{"a", "b"}, []string{"a", "a", "b", "b"}); got != 0 {
		t.Fatalf("equal distributions JSD = %v", got)
	}
	if got := fingerprintJSD([]string{"a"}, []string{"a", "b"}); math.Abs(got-.31127812445913283) > 1e-12 {
		t.Fatalf("JSD = %v", got)
	}
	a, b := repeatedFingerprint("a", 25), repeatedFingerprint("b", 25)
	entries, mean := compareFingerprintCells(a, b)
	if len(entries) != 4 || mean == nil || *mean != 1 {
		t.Fatalf("comparison = %v %v", entries, mean)
	}
	if p := fingerprintPermutation(a, b, "different"); p == nil || *p > .004 {
		t.Fatalf("disjoint p = %v", p)
	}
	if p := fingerprintPermutation(a, a, "equal"); p == nil || *p != 1 {
		t.Fatalf("identical p = %v", p)
	}
	if d := fingerprintSplitHalf(a); d == nil || *d != 0 {
		t.Fatalf("self JSD = %v", d)
	}
	delete(a, fingerprintProbes[0].ID)
	if _, mean := compareFingerprintCells(a, b); mean != nil {
		t.Fatal("three probes must not produce attribution")
	}
	if _, mean := compareFingerprintCells(repeatedFingerprint("a", 9), b); mean != nil {
		t.Fatal("nine samples must not qualify")
	}
}

func TestFingerprintAttribution(t *testing.T) {
	previous := fingerprintBaselines
	t.Cleanup(func() { fingerprintBaselines = previous })
	a, b := repeatedFingerprint("a", 25), repeatedFingerprint("b", 25)
	fingerprintBaselines = []fingerprintBaseline{{Model: "original", Cells: a}, {Model: "substitute", Cells: b}}
	for _, tc := range []struct {
		model   string
		samples map[string][]string
		status  string
	}{
		{"original", a, "consistent"}, {"original", b, "substitution"},
		{"original", repeatedFingerprint("c", 25), "different"}, {"unknown", a, "no_baseline"},
		{"original", repeatedFingerprint("a", 9), "insufficient"},
	} {
		if r := attributeFingerprint(tc.model, tc.samples); r.Status != tc.status {
			t.Fatalf("%s: got %+v, want %s", tc.model, r, tc.status)
		}
	}
	fingerprintBaselines = []fingerprintBaseline{{Model: "other", Cells: a}, {Model: "original", Cells: a}}
	if r := attributeFingerprint("original", a); r.Status != "ambiguous" {
		t.Fatalf("indistinguishable references: %+v", r)
	}
}
