package app

import (
	"math"
	"sort"
)

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
	var mapTotal float64
	var ndcgTotal float64
	missingRuns := 0

	for _, qrel := range qrels {
		results, ok := runIndex[qrel.QueryID]
		if !ok {
			missingRuns++
		}

		grades := qrel.gradeIndex()
		relevantCount := 0
		for _, grade := range grades {
			if grade > 0 {
				relevantCount++
			}
		}

		hits := 0
		firstRelevantRank := 0
		precisionSum := 0.0
		dcg := 0.0

		for i := 0; i < len(results) && i < k; i++ {
			grade := grades[results[i]]
			if grade > 0 {
				hits++
				precisionSum += float64(hits) / float64(i+1)
				if firstRelevantRank == 0 {
					firstRelevantRank = i + 1
				}
			}
			dcg += (math.Pow(2, grade) - 1) / math.Log2(float64(i+2))
		}

		precisionTotal += float64(hits) / float64(k)
		if relevantCount > 0 {
			recallTotal += float64(hits) / float64(relevantCount)
			mapTotal += precisionSum / float64(min(relevantCount, k))
		}
		if hits > 0 {
			hitTotal++
		}
		if firstRelevantRank > 0 {
			mrrTotal += 1 / float64(firstRelevantRank)
		}

		idealGrades := sortedGrades(grades)
		idcg := 0.0
		for i := 0; i < len(idealGrades) && i < k; i++ {
			idcg += (math.Pow(2, idealGrades[i]) - 1) / math.Log2(float64(i+2))
		}
		if idcg > 0 {
			ndcgTotal += dcg / idcg
		}
	}

	queries := float64(len(qrels))
	return Metrics{
		Queries:      len(qrels),
		K:            k,
		MissingRuns:  missingRuns,
		PrecisionAtK: precisionTotal / queries,
		RecallAtK:    recallTotal / queries,
		HitRateAtK:   hitTotal / queries,
		MRRAtK:       mrrTotal / queries,
		MAPAtK:       mapTotal / queries,
		NDCGAtK:      ndcgTotal / queries,
	}
}

func (q QRel) gradeIndex() map[string]float64 {
	if len(q.Grades) > 0 {
		return q.Grades
	}
	grades := make(map[string]float64, len(q.Relevant))
	for _, doc := range q.Relevant {
		grades[doc] = 1
	}
	return grades
}

func sortedGrades(grades map[string]float64) []float64 {
	values := make([]float64, 0, len(grades))
	for _, grade := range grades {
		if grade > 0 {
			values = append(values, grade)
		}
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i] > values[j]
	})
	return values
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
