package telemetry

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/masa-finance/tee-worker/api/args/base"
	"github.com/masa-finance/tee-worker/api/types"
)

var (
	ErrUnmarshalling = errors.New("failed to unmarshal telemetry arguments")
)

type Telemetry = Arguments

// Verify interface implementation
var _ base.JobArgument = (*Arguments)(nil)

// Arguments defines args for Telemetry jobs
type Arguments struct {
	Type types.Capability `json:"type"`
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
}

func (t *Arguments) Validate() error {
	err := t.ValidateCapability(types.TelemetryJob)
	if err != nil {
		return err
	}
	return nil
}

// GetCapability returns the capability of the arguments
func (t *Arguments) GetCapability() types.Capability {
	return t.Type
}

// ValidateCapability validates the capability of the arguments
func (t *Arguments) ValidateCapability(jobType types.JobType) error {
	return jobType.ValidateCapability(&t.Type)
}

// NewArguments creates a new Arguments instance and applies default values immediately
func NewArguments() Arguments {
	args := Arguments{}
	args.SetDefaultValues()
	args.Validate()
	return args
}
