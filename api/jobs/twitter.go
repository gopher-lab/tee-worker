package jobs

import (
	"encoding/json"
	"time"

	"github.com/masa-finance/tee-worker/api/args/twitter"
	"github.com/masa-finance/tee-worker/api/types"
)

// Compile-time check to ensure TwitterParams implements JobParameters
var _ JobParameters = (*TwitterParams)(nil)

type TwitterParams struct {
	JobType types.JobType           `json:"type"`      // Any of the Twitter* job types
	Args    twitter.SearchArguments `json:"arguments"` // Search arguments
}

func (t TwitterParams) Validate(cfg *SearchConfig) error {
	return t.Args.Validate()
}

func (t TwitterParams) Type() types.JobType {
	return t.JobType
}

func (t TwitterParams) Timeout() time.Duration {
	return 0
}

func (t TwitterParams) PollInterval() time.Duration {
	return 0
}

func (t TwitterParams) Arguments(cfg *SearchConfig) map[string]any {
	jsonData, _ := json.Marshal(t.Args)
	var result map[string]any
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil
	}
	return result
}
