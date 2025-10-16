package query

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/masa-finance/tee-worker/api/args/base"
	"github.com/masa-finance/tee-worker/api/types"
)

var (
	ErrSearchOrUrlsRequired = errors.New("either 'search' or 'start_urls' are required")
	ErrUnmarshalling        = errors.New("failed to unmarshal TikTok searchbyquery arguments")
)

const (
	DefaultMaxItems = 10
	DefaultType     = types.CapSearchByQuery
)

// Verify interface implementation
var _ base.JobArgument = (*Arguments)(nil)

type Arguments struct {
	Type      types.Capability `json:"type"`
	Search    []string         `json:"search,omitempty"`
	StartUrls []string         `json:"start_urls,omitempty"`
	MaxItems  uint             `json:"max_items,omitempty"`
	EndPage   uint             `json:"end_page,omitempty"`
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

func (t *Arguments) SetDefaultValues() {
	if t.MaxItems == 0 {
		t.MaxItems = DefaultMaxItems
	}
}

func (t *Arguments) GetCapability() types.Capability {
	return t.Type
}

func (t *Arguments) ValidateCapability(jobType types.JobType) error {
	return jobType.ValidateCapability(&t.Type)
}

func (t *Arguments) Validate() error {
	err := t.ValidateCapability(types.TiktokJob)
	if err != nil {
		return err
	}
	if len(t.Search) == 0 && len(t.StartUrls) == 0 {
		return ErrSearchOrUrlsRequired
	}
	return nil
}

func NewArguments() Arguments {
	args := Arguments{
		Type: types.CapSearchByQuery,
	}
	args.SetDefaultValues()
	return args
}
