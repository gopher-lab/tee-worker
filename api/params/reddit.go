package params

import (
	"time"

	"github.com/masa-finance/tee-worker/v2/api/args/reddit"
	"github.com/masa-finance/tee-worker/v2/api/types"
)

type Reddit struct {
	Params[*reddit.SearchArguments]
}

func (r Reddit) Timeout() time.Duration {
	if r.Args.Type == types.CapSearchCommunities {
		// Apify communities search takes 3-4 minutes
		return 5 * time.Minute
	}
	return 0
}

func (r Reddit) PollInterval() time.Duration {
	if r.Args.Type == types.CapSearchCommunities {
		// Apify communities search takes 3-4 minutes, so don't poll as often
		return 5 * time.Second
	}
	return 0
}
