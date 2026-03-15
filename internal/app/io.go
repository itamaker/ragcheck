package app

import (
	"encoding/json"
	"fmt"
	"os"
)

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

func loadJudgeRecords(path string) ([]JudgeRecord, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read judge input: %w", err)
	}

	var records []JudgeRecord
	if err := json.Unmarshal(body, &records); err != nil {
		return nil, fmt.Errorf("decode judge input: %w", err)
	}
	return records, nil
}
