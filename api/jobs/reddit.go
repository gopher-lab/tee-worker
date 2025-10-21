package jobs

import (
	"encoding/json"
	"time"

	"github.com/masa-finance/tee-worker/api/args/reddit"
	"github.com/masa-finance/tee-worker/api/types"
)

// Compile-time check to ensure RedditParams implements JobParameters
var _ JobParameters = (*RedditParams)(nil)

type RedditParams struct {
	JobType types.JobType          `json:"type"`      // Type of search: 'reddit'
	Args    reddit.SearchArguments `json:"arguments"` // Scrape arguments
}

func (r RedditParams) Validate(cfg *SearchConfig) error {
	return r.Args.Validate()
}

func (r RedditParams) Type() types.JobType {
	return r.JobType
}

func (r RedditParams) Timeout() time.Duration {
	if r.Args.Type == types.CapSearchCommunities {
		// Apify communities search takes 3-4 minutes
		return 5 * time.Minute
	}
	return 0
}

func (r RedditParams) PollInterval() time.Duration {
	if r.Args.Type == types.CapSearchCommunities {
		// Apify communities search takes 3-4 minutes, so don't poll as often
		return 5 * time.Second
	}
	return 0
}

func (r RedditParams) Arguments(cfg *SearchConfig) map[string]any {
	jsonData, _ := json.Marshal(r.Args)
	var result map[string]any
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil
	}
	return result
}
