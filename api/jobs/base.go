package jobs

import (
	"time"

	"github.com/masa-finance/tee-worker/api/types"
)

type JobParameters interface {
	// Validate returns an error if the arguments are invalid
	Validate(cfg *SearchConfig) error
	// Type returns the job type
	Type() types.JobType
	// Arguments converts the job parameter arguments to a Map
	Arguments(cfg *SearchConfig) map[string]any
	// Timeout() returns the timeout to wait when getting results from the tee-worker. Returning 0 means use the default.
	Timeout() time.Duration
	// PollInterval() returns how often to poll the tee-worker for a job's results. Returning 0 means use the default.
	PollInterval() time.Duration
}
