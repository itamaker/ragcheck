package main

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
}
