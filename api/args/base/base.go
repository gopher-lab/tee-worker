package base

import (
	"encoding/json"
	"fmt"

	"github.com/masa-finance/tee-worker/api/types"
)

// JobArgument defines the interface that all job arguments must implement
type JobArgument interface {
	UnmarshalJSON([]byte) error
	GetCapability() types.Capability
	ValidateCapability(jobType types.JobType) error
	SetDefaultValues()
	Validate() error
}

// Verify interface implementation
var _ JobArgument = (*Arguments)(nil)

type Arguments struct {
	Type types.Capability `json:"type"`
}

func (t *Arguments) UnmarshalJSON(data []byte) error {
	type Alias Arguments
	aux := &struct{ *Alias }{Alias: (*Alias)(t)}
	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("%v: %w", "failed to unmarshal arguments", err)
	}
	t.SetDefaultValues()
	return t.Validate()
}

func (a *Arguments) GetCapability() types.Capability {
	return a.Type
}

func (a *Arguments) ValidateCapability(jobType types.JobType) error {
	return jobType.ValidateCapability(&a.Type)
}

func (a *Arguments) SetDefaultValues() {
}

func (a *Arguments) Validate() error {
	return nil
}
