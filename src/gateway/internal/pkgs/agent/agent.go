package agent

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type checkingAgent struct {
	factRunnable  compose.Runnable[map[string]any, *schema.Message]
	checkRunnable compose.Runnable[map[string]any, *schema.Message]
}

type Facts struct {
	FactChecks []FactCheck `json:"fact_checks"`
}

type FactCheck struct {
	Type      string `json:"type"`
	Excerpt   string `json:"excerpt"`
	Highlight string `json:"highlight"`
	Note      string `json:"note"`
}

type CheckingResult struct {
	Result  string   `json:"result"`
	Sources []string `json:"sources"`
	Reason  string   `json:"reason"`
}

func NewCheckingAgent(ctx context.Context) (*checkingAgent, error) {
	ftpl := newFactTemplate(ctx)
	factCm, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	fg, err := buildFactGraph(ftpl, factCm)
	if err != nil {
		return nil, err
	}
	fr, err := fg.Compile(ctx)
	if err != nil {
		return nil, err
	}

	ctpl := newCheckTemplate(ctx)
	checkCm, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	cg, err := buildCheckGraph(ctx, ctpl, checkCm)
	if err != nil {
		return nil, err
	}
	cr, err := cg.Compile(ctx)
	if err != nil {
		return nil, err
	}

	return &checkingAgent{factRunnable: fr, checkRunnable: cr}, nil
}

func (a *checkingAgent) Run(ctx context.Context, text string) ([]CheckingResult, error) {
	factOutput, err := a.factRunnable.Invoke(ctx, map[string]any{
		"text": text,
	})
	if err != nil {
		return nil, err
	}
	facts := Facts{}
	if err := json.Unmarshal([]byte(factOutput.Content), &facts); err != nil {
		return nil, err
	}
	results := []CheckingResult{}
	resultChan := make(chan CheckingResult, len(facts.FactChecks))
	var wg sync.WaitGroup
	for _, factCheck := range facts.FactChecks {
		wg.Add(1)
		go func(factCheck FactCheck) {
			defer wg.Done()
			checkOutput, err := a.checkRunnable.Invoke(ctx, map[string]any{
				"input": factCheck,
			})
			if err != nil {
				resultChan <- CheckingResult{Reason: err.Error()}
				return
			}
			checkingResult := CheckingResult{}
			if err := json.Unmarshal([]byte(checkOutput.Content), &checkingResult); err != nil {
				resultChan <- CheckingResult{Reason: err.Error()}
				return
			}
			resultChan <- checkingResult
		}(factCheck)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	for result := range resultChan {
		if result.Result == "" {
			result.Result = "FAILED"
		}
		results = append(results, result)
	}
	return results, nil
}
