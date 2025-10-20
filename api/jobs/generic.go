package jobs

import (
	"encoding/json"
	"maps"
	"time"

	"github.com/masa-finance/tee-worker/api/args"
	"github.com/masa-finance/tee-worker/api/types"
)

// Compile-time check to ensure GenericParams implements JobParameters
var _ JobParameters = (*GenericParams)(nil)

// This is a generic Params struct that assumes the args are just a Map, and the validation will be done by the appropriate JSON unmarshaller in tee-types. Note that we're unmarshalling at least twice which will have a (probably heavy) runtime cost.
type GenericParams struct {
	JobType types.JobType  `json:"type"`
	Args    map[string]any `json:"arguments"`
}

func (p GenericParams) Validate(_ *SearchConfig) error {
	_, err := args.UnmarshalJobArguments(p.JobType, p.Args)
	return err
}

func (p GenericParams) Arguments(_ *SearchConfig) map[string]any {
	// Use UnmarshalJobArguments to get properly typed arguments with correct type
	ja, err := args.UnmarshalJobArguments(p.JobType, p.Args)
	if err != nil {
		// Fallback to original args if unmarshaling fails
		return maps.Clone(map[string]any(p.Args))
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

func (p GenericParams) Type() types.JobType {
	return p.JobType
}

func (p GenericParams) Timeout() time.Duration {
	return 0
}

func (p GenericParams) PollInterval() time.Duration {
	return 0
}
