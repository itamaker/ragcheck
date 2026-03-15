package app

import "testing"

func TestScore(t *testing.T) {
	t.Parallel()

	qrels := []QRel{
		{QueryID: "q1", Relevant: []string{"d1", "d3"}},
		{QueryID: "q2", Relevant: []string{"d4"}},
	}
	run := []RunResult{
		{QueryID: "q1", Results: []string{"d1", "d2", "d3"}},
		{QueryID: "q2", Results: []string{"d5", "d4", "d6"}},
	}

	metrics := score(qrels, run, 3)
	if metrics.Queries != 2 {
		t.Fatalf("Queries = %d, want 2", metrics.Queries)
	}
	if metrics.HitRateAtK != 1 {
		t.Fatalf("HitRateAtK = %.3f, want 1.000", metrics.HitRateAtK)
	}
	if metrics.MRRAtK <= 0.6 {
		t.Fatalf("MRRAtK = %.3f, want > 0.6", metrics.MRRAtK)
	}
	if metrics.MAPAtK <= 0.6 {
		t.Fatalf("MAPAtK = %.3f, want > 0.6", metrics.MAPAtK)
	}
	if metrics.NDCGAtK <= 0.7 {
		t.Fatalf("NDCGAtK = %.3f, want > 0.7", metrics.NDCGAtK)
	}
}

func TestJudge(t *testing.T) {
	t.Parallel()

	report := judge([]JudgeRecord{
		{
			QueryID:   "q1",
			Question:  "What causes retrieval latency spikes?",
			Answer:    "Latency spikes come from cache misses and slow vector search.",
			Reference: "Cache misses and slow vector search cause retrieval latency spikes.",
			Contexts: []string{
				"The incident report says cache misses increased tail latency.",
				"Slow vector search also contributed to retrieval latency spikes.",
			},
		},
	})

	if report.Queries != 1 {
		t.Fatalf("Queries = %d, want 1", report.Queries)
	}
	if report.Groundedness <= 0.4 {
		t.Fatalf("Groundedness = %.3f, want > 0.4", report.Groundedness)
	}
	if report.ReferenceCoverage <= 0.4 {
		t.Fatalf("ReferenceCoverage = %.3f, want > 0.4", report.ReferenceCoverage)
	}
	if len(report.Results) != 1 {
		t.Fatalf("expected one per-record result")
	}
}
