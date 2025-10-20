package jobs

import (
	"encoding/json"
	"maps"
	"time"

	"github.com/masa-finance/tee-worker/api/args"
	"github.com/masa-finance/tee-worker/api/args/tiktok"
	"github.com/masa-finance/tee-worker/api/types"
)

// Compile-time check to ensure TikTokParams implements JobParameters
var _ JobParameters = (*TikTokParams)(nil)

type TikTokTranscriptionParams struct {
	JobType types.JobType        `json:"type"`
	Args    tiktok.Transcription `json:"arguments"`
}

type TikTokSearchParams struct {
	JobType types.JobType `json:"type"`
	Args    tiktok.Query  `json:"arguments"`
}

type TikTokTrendingParams struct {
	JobType types.JobType   `json:"type"`
	Args    tiktok.Trending `json:"arguments"`
}

// TikTokArguments is a flexible map that supports multiple unique capabilities
type TikTokArguments map[string]any

type TikTokParams struct {
	JobType types.JobType   `json:"type"`
	Args    TikTokArguments `json:"arguments"`
}

func (t TikTokParams) Type() types.JobType {
	return t.JobType
}

func (t TikTokParams) Validate(cfg *SearchConfig) error {
	_, err := args.UnmarshalJobArguments(t.JobType, t.Args)
	return err
}

func (t TikTokParams) Timeout() time.Duration {
	return 0
}

func (t TikTokParams) PollInterval() time.Duration {
	return 0
}

func (t TikTokParams) Arguments(cfg *SearchConfig) map[string]any {
	// Use UnmarshalJobArguments to get properly typed arguments with correct type
	ja, err := args.UnmarshalJobArguments(types.TiktokJob, t.Args)
	if err != nil {
		// Fallback to original args if unmarshaling fails
		return maps.Clone(map[string]any(t.Args))
	}

	// Marshal the properly typed arguments back to map[string]any
	// This will include the correct type and all other fields
	jsonData, _ := json.Marshal(ja)
	var result map[string]any
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil
	}
	return result
}
