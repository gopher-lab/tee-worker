package args

import (
	"github.com/masa-finance/tee-worker/api/types"
)

// TelemetryJobArguments for telemetry jobs (simple case)
type TelemetryJobArguments struct{}

func (t *TelemetryJobArguments) Validate() error {
	return nil
}

func (t *TelemetryJobArguments) GetCapability() types.Capability {
	return types.CapTelemetry
}
