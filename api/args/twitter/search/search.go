package search

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/masa-finance/tee-worker/api/args/base"
	"github.com/masa-finance/tee-worker/api/types"
)

var (
	ErrCountNegative      = errors.New("count must be non-negative")
	ErrCountTooLarge      = errors.New("count must be less than or equal to 1000")
	ErrMaxResultsTooLarge = errors.New("max_results must be less than or equal to 1000")
	ErrMaxResultsNegative = errors.New("max_results must be non-negative")
	ErrUnmarshalling      = errors.New("failed to unmarshal twitter search arguments")
)

const (
	MaxResults        = 1000
	DefaultMaxResults = 10
)

// Verify interface implementation
var _ base.JobArgument = (*Arguments)(nil)

// Arguments defines args for Twitter searches
type Arguments struct {
	Type       types.Capability `json:"type"`
	Query      string           `json:"query"` // Username or search query
	Count      int              `json:"count"`
	StartTime  string           `json:"start_time"`  // Optional ISO timestamp
	EndTime    string           `json:"end_time"`    // Optional ISO timestamp
	MaxResults int              `json:"max_results"` // Optional, max number of results
	NextCursor string           `json:"next_cursor"`
}

func (t *Arguments) UnmarshalJSON(data []byte) error {
	type Alias Arguments
	aux := &struct{ *Alias }{Alias: (*Alias)(t)}
	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("%w: %w", ErrUnmarshalling, err)
	}
	t.SetDefaultValues()
	return t.Validate()
}

// SetDefaultValues sets default values for the arguments
func (t *Arguments) SetDefaultValues() {
	if t.MaxResults == 0 {
		t.MaxResults = DefaultMaxResults
	}
}

// Validate validates the  arguments (general validation)
func (t *Arguments) Validate() error {
	// note, query is not required for all capabilities
	err := t.ValidateCapability(types.TwitterJob)
	if err != nil {
		return err
	}
	if t.Count < 0 {
		return fmt.Errorf("%w, got: %d", ErrCountNegative, t.Count)
	}
	if t.Count > MaxResults {
		return fmt.Errorf("%w, got: %d", ErrCountTooLarge, t.Count)
	}
	if t.MaxResults < 0 {
		return fmt.Errorf("%w, got: %d", ErrMaxResultsNegative, t.MaxResults)
	}
	if t.MaxResults > MaxResults {
		return fmt.Errorf("%w, got: %d", ErrMaxResultsTooLarge, t.MaxResults)
	}

	return nil
}

func (t *Arguments) GetCapability() types.Capability {
	return t.Type
}

func (t *Arguments) ValidateCapability(jobType types.JobType) error {
	return jobType.ValidateCapability(&t.Type)
}

func (t *Arguments) IsSingleTweetOperation() bool {
	return t.GetCapability() == types.CapGetById
}

func (t *Arguments) IsMultipleTweetOperation() bool {
	c := t.GetCapability()
	return c == types.CapSearchByQuery ||
		c == types.CapSearchByFullArchive ||
		c == types.CapGetTweets ||
		c == types.CapGetReplies ||
		c == types.CapGetMedia
}

func (t *Arguments) IsSingleProfileOperation() bool {
	c := t.GetCapability()
	return c == types.CapGetProfileById ||
		c == types.CapSearchByProfile
}

func (t *Arguments) IsMultipleProfileOperation() bool {
	c := t.GetCapability()
	return c == types.CapGetRetweeters
}

func (t *Arguments) IsSingleSpaceOperation() bool {
	return t.GetCapability() == types.CapGetSpace
}

func (t *Arguments) IsTrendsOperation() bool {
	return t.GetCapability() == types.CapGetTrends
}

func NewArguments() Arguments {
	args := Arguments{}
	args.SetDefaultValues()
	args.Validate()
	return args
}
