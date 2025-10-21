package jobs

import (
	"encoding/json"
	"time"

	"github.com/masa-finance/tee-worker/api/args/linkedin"
	"github.com/masa-finance/tee-worker/api/types"
)

// Compile-time check to ensure LinkedInParams implements JobParameters
var _ JobParameters = (*LinkedInParams)(nil)

type LinkedInParams struct {
	JobType types.JobType             `json:"type"`
	Args    linkedin.ProfileArguments `json:"arguments"`
}

func (l LinkedInParams) Validate(cfg *SearchConfig) error {
	return l.Args.Validate()
}

func (l LinkedInParams) Type() types.JobType {
	return l.JobType
}

func (l LinkedInParams) Timeout() time.Duration {
	return 0
}

func (l LinkedInParams) PollInterval() time.Duration {
	return 0
}

func (l LinkedInParams) Arguments(cfg *SearchConfig) map[string]any {
	jsonData, _ := json.Marshal(l.Args)
	var result map[string]any
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil
	}
	return result
}
