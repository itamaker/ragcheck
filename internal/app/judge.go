package app

import (
	"sort"
	"strings"
	"unicode"
)

func judge(records []JudgeRecord) JudgeReport {
	report := JudgeReport{
		Queries: len(records),
	}
	if len(records) == 0 {
		return report
	}

	for _, record := range records {
		questionTokens := tokenize(record.Question)
		answerTokens := tokenize(record.Answer)
		referenceTokens := tokenize(record.Reference)
		contextTokens := tokenize(strings.Join(record.Contexts, " "))

		answerRelevance := overlapRatio(answerTokens, questionTokens)
		contextRelevance := overlapRatio(contextTokens, questionTokens)
		groundedness := overlapRatio(answerTokens, contextTokens)
		referenceCoverage := 0.0
		if len(referenceTokens) > 0 {
			referenceCoverage = overlapRatio(answerTokens, referenceTokens)
		}

		result := JudgeResult{
			QueryID:             record.QueryID,
			AnswerRelevance:     answerRelevance,
			ContextRelevance:    contextRelevance,
			Groundedness:        groundedness,
			ReferenceCoverage:   referenceCoverage,
			UnsupportedTokens:   topUnsupported(answerTokens, contextTokens, 5),
			MissingContextHints: topUnsupported(referenceTokens, contextTokens, 5),
		}

		report.Results = append(report.Results, result)
		report.AnswerRelevance += answerRelevance
		report.ContextRelevance += contextRelevance
		report.Groundedness += groundedness
		report.ReferenceCoverage += referenceCoverage
	}

	divisor := float64(len(records))
	report.AnswerRelevance /= divisor
	report.ContextRelevance /= divisor
	report.Groundedness /= divisor
	report.ReferenceCoverage /= divisor
	return report
}

func tokenize(text string) []string {
	normalized := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return ' '
	}, text)

	stopwords := map[string]struct{}{
		"the": {}, "and": {}, "for": {}, "with": {}, "that": {}, "this": {}, "from": {},
		"into": {}, "are": {}, "was": {}, "were": {}, "have": {}, "has": {}, "had": {},
		"your": {}, "about": {}, "then": {}, "than": {}, "when": {}, "what": {}, "why": {},
		"how": {}, "using": {}, "used": {}, "their": {}, "there": {}, "will": {},
	}

	var tokens []string
	for _, token := range strings.Fields(normalized) {
		if len(token) < 3 {
			continue
		}
		if _, ok := stopwords[token]; ok {
			continue
		}
		tokens = append(tokens, token)
	}
	return tokens
}

func overlapRatio(source []string, reference []string) float64 {
	if len(source) == 0 || len(reference) == 0 {
		return 0
	}
	refSet := map[string]struct{}{}
	for _, token := range reference {
		refSet[token] = struct{}{}
	}
	matched := 0
	seen := map[string]struct{}{}
	for _, token := range source {
		if _, counted := seen[token]; counted {
			continue
		}
		seen[token] = struct{}{}
		if _, ok := refSet[token]; ok {
			matched++
		}
	}
	return float64(matched) / float64(len(seen))
}

func topUnsupported(source []string, reference []string, limit int) []string {
	if len(source) == 0 || len(reference) == 0 {
		return nil
	}
	refSet := map[string]struct{}{}
	for _, token := range reference {
		refSet[token] = struct{}{}
	}
	var unsupported []string
	seen := map[string]struct{}{}
	for _, token := range source {
		if _, ok := refSet[token]; ok {
			continue
		}
		if _, ok := seen[token]; ok {
			continue
		}
		seen[token] = struct{}{}
		unsupported = append(unsupported, token)
	}
	sort.Strings(unsupported)
	if len(unsupported) > limit {
		unsupported = unsupported[:limit]
	}
	return unsupported
}
