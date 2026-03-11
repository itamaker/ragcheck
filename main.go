package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type QRel struct {
	QueryID  string   `json:"query_id"`
	Relevant []string `json:"relevant"`
}

type RunResult struct {
	QueryID string   `json:"query_id"`
	Results []string `json:"results"`
}

type Metrics struct {
	Queries      int     `json:"queries"`
	K            int     `json:"k"`
	PrecisionAtK float64 `json:"precision_at_k"`
	RecallAtK    float64 `json:"recall_at_k"`
	HitRateAtK   float64 `json:"hit_rate_at_k"`
	MRRAtK       float64 `json:"mrr_at_k"`
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		usage()
		return 2
	}

	switch args[0] {
	case "score":
		return runScore(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", args[0])
		usage()
		return 2
	}
}

func runScore(args []string) int {
	fs := flag.NewFlagSet("score", flag.ContinueOnError)
	qrelsPath := fs.String("qrels", "", "path to qrels JSON")
	runPath := fs.String("run", "", "path to retrieval run JSON")
	k := fs.Int("k", 5, "top-k cutoff")
	jsonOutput := fs.Bool("json", false, "emit machine-readable JSON")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *qrelsPath == "" || *runPath == "" {
		fmt.Fprintln(os.Stderr, "both -qrels and -run are required")
		return 2
	}

	qrels, err := loadQrels(*qrelsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	runResults, err := loadRun(*runPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	metrics := score(qrels, runResults, *k)
	if *jsonOutput {
		body, err := json.MarshalIndent(metrics, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Println(string(body))
		return 0
	}

	fmt.Printf("Queries: %d\n", metrics.Queries)
	fmt.Printf("Precision@%d: %.3f\n", metrics.K, metrics.PrecisionAtK)
	fmt.Printf("Recall@%d: %.3f\n", metrics.K, metrics.RecallAtK)
	fmt.Printf("HitRate@%d: %.3f\n", metrics.K, metrics.HitRateAtK)
	fmt.Printf("MRR@%d: %.3f\n", metrics.K, metrics.MRRAtK)
	return 0
}

func usage() {
	fmt.Println("ragcheck scores retrieval runs.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  ragcheck score -qrels examples/qrels.json -run examples/run.json -k 3")
}

func loadQrels(path string) ([]QRel, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read qrels: %w", err)
	}

	var qrels []QRel
	if err := json.Unmarshal(body, &qrels); err != nil {
		return nil, fmt.Errorf("decode qrels: %w", err)
	}
	return qrels, nil
}

func loadRun(path string) ([]RunResult, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read run: %w", err)
	}

	var runResults []RunResult
	if err := json.Unmarshal(body, &runResults); err != nil {
		return nil, fmt.Errorf("decode run: %w", err)
	}
	return runResults, nil
}

func score(qrels []QRel, runResults []RunResult, k int) Metrics {
	if len(qrels) == 0 || k <= 0 {
		return Metrics{K: k}
	}

	runIndex := map[string][]string{}
	for _, result := range runResults {
		runIndex[result.QueryID] = result.Results
	}

	var precisionTotal float64
	var recallTotal float64
	var hitTotal float64
	var mrrTotal float64

	for _, qrel := range qrels {
		results := runIndex[qrel.QueryID]
		relevant := make(map[string]struct{}, len(qrel.Relevant))
		for _, doc := range qrel.Relevant {
			relevant[doc] = struct{}{}
		}

		hits := 0
		firstRelevantRank := 0
		for i := 0; i < len(results) && i < k; i++ {
			if _, ok := relevant[results[i]]; ok {
				hits++
				if firstRelevantRank == 0 {
					firstRelevantRank = i + 1
				}
			}
		}

		precisionTotal += float64(hits) / float64(k)
		if len(relevant) > 0 {
			recallTotal += float64(hits) / float64(len(relevant))
		}
		if hits > 0 {
			hitTotal++
		}
		if firstRelevantRank > 0 {
			mrrTotal += 1 / float64(firstRelevantRank)
		}
	}

	queries := float64(len(qrels))
	return Metrics{
		Queries:      len(qrels),
		K:            k,
		PrecisionAtK: precisionTotal / queries,
		RecallAtK:    recallTotal / queries,
		HitRateAtK:   hitTotal / queries,
		MRRAtK:       mrrTotal / queries,
	}
}
