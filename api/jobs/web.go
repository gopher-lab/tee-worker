package jobs

import (
	"encoding/json"
	"time"

	"github.com/masa-finance/tee-worker/api/args/web"
	"github.com/masa-finance/tee-worker/api/types"
)

// Compile-time check to ensure WebParams implements JobParameters
var _ JobParameters = (*WebParams)(nil)

type WebParams struct {
	JobType types.JobType        `json:"type"`
	Args    web.ScraperArguments `json:"arguments"`
}

func (w WebParams) Validate(cfg *SearchConfig) error {
	return w.Args.Validate()
}

func (w WebParams) Timeout() time.Duration {
	return 0
}

func (w WebParams) Type() types.JobType {
	return w.JobType
}

func (w WebParams) PollInterval() time.Duration {
	return 0
}

func (w WebParams) Arguments(cfg *SearchConfig) map[string]any {
	jsonData, _ := json.Marshal(w.Args)
	var result map[string]any
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil
	}
	return result
}
