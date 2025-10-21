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

// Verify interface implementation
var _ base.JobArgument = (*Arguments)(nil)

// Arguments defines args for Telemetry jobs
type Arguments struct {
	base.Arguments
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

func (t *Arguments) Validate() error {
	if err := types.TelemetryJob.ValidateCapability(&t.Type); err != nil {
		return err
	}
	return nil
}

// NewArguments creates a new Arguments instance and applies default values immediately
func NewArguments() Arguments {
	args := Arguments{}
	args.SetDefaultValues()
	args.Validate()
	return args
}
