package linkedinapify

import (
	"encoding/json"
	"fmt"

	profileArgs "github.com/masa-finance/tee-worker/api/args/linkedin/profile"
	profileTypes "github.com/masa-finance/tee-worker/api/types/linkedin/profile"
	"github.com/masa-finance/tee-worker/internal/apify"
	"github.com/masa-finance/tee-worker/internal/jobs/stats"
	"github.com/masa-finance/tee-worker/pkg/client"
	"github.com/sirupsen/logrus"
)

type ApifyClient struct {
	client         client.Apify
	statsCollector *stats.StatsCollector
}

// NewInternalClient is a function variable that can be replaced in tests.
// It defaults to the actual implementation.
var NewInternalClient = func(apiKey string) (client.Apify, error) {
	return client.NewApifyClient(apiKey)
}

// NewClient creates a new LinkedIn Apify client
func NewClient(apiToken string, statsCollector *stats.StatsCollector) (*ApifyClient, error) {
	client, err := NewInternalClient(apiToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create apify client: %w", err)
	}

	return &ApifyClient{
		client:         client,
		statsCollector: statsCollector,
	}, nil
}

// ValidateApiKey tests if the Apify API token is valid
func (c *ApifyClient) ValidateApiKey() error {
	return c.client.ValidateApiKey()
}

func (c *ApifyClient) SearchProfiles(workerID string, args *profileArgs.Arguments, cursor client.Cursor) ([]*profileTypes.Profile, string, client.Cursor, error) {
	if c.statsCollector != nil {
		c.statsCollector.Add(workerID, stats.LinkedInQueries, 1)
	}

	requestBytes, err := json.Marshal(args)
	if err != nil {
		return nil, "", client.EmptyCursor, fmt.Errorf("failed to marshal request: %w", err)
	}

	var input map[string]any
	if err := json.Unmarshal(requestBytes, &input); err != nil {
		return nil, "", client.EmptyCursor, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	dataset, nextCursor, err := c.client.RunActorAndGetResponse(apify.ActorIds.LinkedInSearchProfile, input, cursor, args.MaxItems)
	if err != nil {
		if c.statsCollector != nil {
			c.statsCollector.Add(workerID, stats.LinkedInErrors, 1)
		}
		return nil, "", client.EmptyCursor, err
	}

	response := make([]*profileTypes.Profile, 0, len(dataset.Data.Items))

	for i, item := range dataset.Data.Items {
		var resp profileTypes.Profile
		if err := json.Unmarshal(item, &resp); err != nil {
			logrus.Warnf("Failed to unmarshal scrape result at index %d: %v", i, err)
			continue
		}
		response = append(response, &resp)
	}

	if c.statsCollector != nil {
		c.statsCollector.Add(workerID, stats.LinkedInProfiles, uint(len(response)))
	}

	return response, dataset.DatasetId, nextCursor, nil
}
