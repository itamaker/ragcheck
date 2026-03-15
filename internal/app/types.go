package app

type QRel struct {
	QueryID  string             `json:"query_id"`
	Relevant []string           `json:"relevant,omitempty"`
	Grades   map[string]float64 `json:"grades,omitempty"`
}

type RunResult struct {
	QueryID string   `json:"query_id"`
	Results []string `json:"results"`
}

type Metrics struct {
	Queries      int     `json:"queries"`
	K            int     `json:"k"`
	MissingRuns  int     `json:"missing_runs"`
	PrecisionAtK float64 `json:"precision_at_k"`
	RecallAtK    float64 `json:"recall_at_k"`
	HitRateAtK   float64 `json:"hit_rate_at_k"`
	MRRAtK       float64 `json:"mrr_at_k"`
	MAPAtK       float64 `json:"map_at_k"`
	NDCGAtK      float64 `json:"ndcg_at_k"`
}

type JudgeRecord struct {
	QueryID   string   `json:"query_id"`
	Question  string   `json:"question"`
	Answer    string   `json:"answer"`
	Reference string   `json:"reference,omitempty"`
	Contexts  []string `json:"contexts"`
}

type JudgeResult struct {
	QueryID             string   `json:"query_id"`
	AnswerRelevance     float64  `json:"answer_relevance"`
	ContextRelevance    float64  `json:"context_relevance"`
	Groundedness        float64  `json:"groundedness"`
	ReferenceCoverage   float64  `json:"reference_coverage"`
	UnsupportedTokens   []string `json:"unsupported_tokens,omitempty"`
	MissingContextHints []string `json:"missing_context_hints,omitempty"`
}

type JudgeReport struct {
	Queries           int           `json:"queries"`
	AnswerRelevance   float64       `json:"answer_relevance"`
	ContextRelevance  float64       `json:"context_relevance"`
	Groundedness      float64       `json:"groundedness"`
	ReferenceCoverage float64       `json:"reference_coverage"`
	Results           []JudgeResult `json:"results,omitempty"`
}
