package params

import (
	"encoding/json"
	"time"

	"github.com/masa-finance/tee-worker/api/args"
	"github.com/masa-finance/tee-worker/api/args/base"
	"github.com/masa-finance/tee-worker/api/types"
)

type JobParameters interface {
	// Validate returns an error if the arguments are invalid
	Validate(cfg *types.SearchConfig) error
	// Type returns the job type
	Type() types.JobType
	// Arguments converts the job parameter arguments to a Map
	Arguments(cfg *types.SearchConfig) map[string]any
	// Timeout() returns the timeout to wait when getting results from the tee-worker. Returning 0 means use the default.
	Timeout() time.Duration
	// PollInterval() returns how often to poll the tee-worker for a job's results. Returning 0 means use the default.
	PollInterval() time.Duration
}

// Compile-time check to ensure LinkedInParams implements JobParameters
var _ JobParameters = (*Params[base.JobArgument])(nil)

type Params[T base.JobArgument] struct {
	JobType types.JobType `json:"type"`
	Args    T             `json:"arguments"`
}

func (p Params[T]) Validate(_ *types.SearchConfig) error {
	return p.Args.Validate()
}

func (p Params[T]) Type() types.JobType {
	return p.JobType
}

func (l Params[T]) Timeout() time.Duration {
	return 0
}

func (l Params[T]) PollInterval() time.Duration {
	return 0
}

// TODO: revisit this...
// We marshal 3 times because:
// 1. Convert generic T to map[string]any (for UnmarshalJobArguments)
// 2. UnmarshalJobArguments validates/transforms the args based on job type
// 3. Convert the validated args back to map[string]any for the final result
func (l Params[T]) Arguments(cfg *types.SearchConfig) map[string]any {
	// Convert l.Args to map[string]any via JSON marshal/unmarshal
	jsonData, _ := json.Marshal(l.Args)
	var argsMap map[string]any
	json.Unmarshal(jsonData, &argsMap)

	// Use UnmarshalJobArguments to get properly typed arguments with correct type
	ja, err := args.UnmarshalJobArguments(l.JobType, argsMap)
	if err != nil {
		// Fallback to original args if unmarshaling fails
		return argsMap
	}

	// Convert ja to map[string]any
	resultData, _ := json.Marshal(ja)
	var result map[string]any
	if err := json.Unmarshal(resultData, &result); err != nil {
		return nil
	}
	return result
}
