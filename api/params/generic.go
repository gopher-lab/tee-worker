package params

import (
	"github.com/masa-finance/tee-worker/api/args/base"
)

type GenericArgs struct {
	base.Arguments
	Data map[string]any `json:",inline"`
}

type Generic = Params[*GenericArgs]
