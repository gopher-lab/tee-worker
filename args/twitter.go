package args

import (
	"encoding/json"
	"errors"
	"fmt"

	teetypes "github.com/masa-finance/tee-types/types"
)

var (
	ErrTwitterCountNegative      = errors.New("count must be non-negative")
	ErrTwitterCountTooLarge      = errors.New("count must be less than or equal to 1000")
	ErrTwitterMaxResultsTooLarge = errors.New("max_results must be less than or equal to 1000")
	ErrTwitterMaxResultsNegative = errors.New("max_results must be non-negative")
)

const (
	TwitterMaxResults = 1000
)

// TwitterSearchArguments defines args for Twitter searches
type TwitterSearchArguments struct {
	QueryType  teetypes.Capability `json:"type"`  // Optional, type of search
	Query      string              `json:"query"` // Username or search query
	Count      int                 `json:"count"`
	StartTime  string              `json:"start_time"`  // Optional ISO timestamp
	EndTime    string              `json:"end_time"`    // Optional ISO timestamp
	MaxResults int                 `json:"max_results"` // Optional, max number of results
	NextCursor string              `json:"next_cursor"`
}

// UnmarshalJSON implements custom JSON unmarshaling with validation
func (t *TwitterSearchArguments) UnmarshalJSON(data []byte) error {
	// Prevent infinite recursion (you call json.Unmarshal which then calls `UnmarshalJSON`, which then calls `json.Unmarshal`...)
	type Alias TwitterSearchArguments
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("failed to unmarshal Twitter arguments: %w", err)
	}

	return t.Validate()
}

// Validate validates the Twitter arguments (general validation)
func (t *TwitterSearchArguments) Validate() error {
	// note, query is not required for all capabilities
	if t.Count < 0 {
		return fmt.Errorf("%w, got: %d", ErrTwitterCountNegative, t.Count)
	}
	if t.Count > TwitterMaxResults {
		return fmt.Errorf("%w, got: %d", ErrTwitterCountTooLarge, t.Count)
	}
	if t.MaxResults < 0 {
		return fmt.Errorf("%w, got: %d", ErrTwitterMaxResultsNegative, t.MaxResults)
	}
	if t.MaxResults > TwitterMaxResults {
		return fmt.Errorf("%w, got: %d", ErrTwitterMaxResultsTooLarge, t.MaxResults)
	}

	return nil
}

// ValidateForJobType validates Twitter arguments for a specific job type
func (t *TwitterSearchArguments) ValidateForJobType(jobType teetypes.JobType) error {
	if err := t.Validate(); err != nil {
		return err
	}

	// Validate QueryType against job-specific capabilities
	return jobType.ValidateCapability(teetypes.Capability(t.QueryType))
}

// GetCapability returns the QueryType as a typed Capability
func (t *TwitterSearchArguments) GetCapability() teetypes.Capability {
	return teetypes.Capability(t.QueryType)
}

func (t *TwitterSearchArguments) IsSingleTweetOperation() bool {
	capability := t.GetCapability()
	return capability == teetypes.CapGetById
}

func (t *TwitterSearchArguments) IsMultipleTweetOperation() bool {
	capability := t.GetCapability()
	return capability == teetypes.CapSearchByQuery ||
		capability == teetypes.CapSearchByFullArchive ||
		capability == teetypes.CapGetHomeTweets ||
		capability == teetypes.CapGetForYouTweets ||
		capability == teetypes.CapGetTweets ||
		capability == teetypes.CapGetReplies ||
		capability == teetypes.CapGetMedia
}

func (t *TwitterSearchArguments) IsSingleProfileOperation() bool {
	capability := t.GetCapability()
	return capability == teetypes.CapGetProfileById ||
		capability == teetypes.CapSearchByProfile
}

func (t *TwitterSearchArguments) IsMultipleProfileOperation() bool {
	capability := t.GetCapability()
	return capability == teetypes.CapGetFollowing ||
		capability == teetypes.CapGetFollowers ||
		capability == teetypes.CapGetRetweeters
}

func (t *TwitterSearchArguments) IsSingleSpaceOperation() bool {
	capability := t.GetCapability()
	return capability == teetypes.CapGetSpace
}

func (t *TwitterSearchArguments) IsTrendsOperation() bool {
	capability := t.GetCapability()
	return capability == teetypes.CapGetTrends
}
