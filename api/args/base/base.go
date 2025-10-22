package base

import (
	"github.com/masa-finance/tee-worker/api/types"
)

// JobArgument defines the interface that all job arguments must implement
type JobArgument interface {
	GetCapability() types.Capability
	SetDefaultValues()
	Validate() error
}

// Verify interface implementation
var _ JobArgument = (*Arguments)(nil)

type Arguments struct {
	Type types.Capability `json:"type"`
}

func (a *Arguments) GetCapability() types.Capability {
	return a.Type
}

func (a *Arguments) SetDefaultValues() {
}

func (a *Arguments) Validate() error {
	return nil
}
