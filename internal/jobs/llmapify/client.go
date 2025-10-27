package llmapify

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/masa-finance/tee-worker/v2/api/args/llm"
	"github.com/masa-finance/tee-worker/v2/api/types"
	"github.com/masa-finance/tee-worker/v2/internal/apify"
	"github.com/masa-finance/tee-worker/v2/internal/config"
	"github.com/masa-finance/tee-worker/v2/internal/jobs/stats"
	"github.com/masa-finance/tee-worker/v2/pkg/client"
	"github.com/sirupsen/logrus"
)

var (
	ErrFailedToCreateClient = errors.New("failed to create apify client")
)

type ApifyClient struct {
	client         client.Apify
	statsCollector *stats.StatsCollector
	llmConfig      config.LlmConfig
}

// NewInternalClient is a function variable that can be replaced in tests.
// It defaults to the actual implementation.
var NewInternalClient = func(apiKey string) (client.Apify, error) {
	return client.NewApifyClient(apiKey)
}

// NewClient creates a new LLM Apify client
func NewClient(apiToken string, llmConfig config.LlmConfig, statsCollector *stats.StatsCollector) (*ApifyClient, error) {
	client, err := NewInternalClient(apiToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedToCreateClient, err)
	}

	llmErr := llmConfig.HasValidKey()
	if llmErr != nil {
		return nil, llmErr
	}

	return &ApifyClient{
		client:         client,
		statsCollector: statsCollector,
		llmConfig:      llmConfig,
	}, nil
}

// ValidateApiKey tests if the Apify API token is valid
func (c *ApifyClient) ValidateApiKey() error {
	return c.client.ValidateApiKey()
}

func (c *ApifyClient) Process(workerID string, args llm.ProcessArguments, cursor client.Cursor) ([]*types.LLMProcessorResult, client.Cursor, error) {
	if c.statsCollector != nil {
		c.statsCollector.Add(workerID, stats.LLMQueries, 1)
	}

	model, key, err := c.llmConfig.GetModelAndKey()
	if err != nil {
		return nil, client.EmptyCursor, err
	}

	input, err := args.ToProcessorRequest(model, key)
	if err != nil {
		return nil, client.EmptyCursor, err
	}

	limit := uint(args.Items)
	dataset, nextCursor, err := c.client.RunActorAndGetResponse(apify.ActorIds.LLMDatasetProcessor, input, cursor, limit)
	if err != nil {
		if c.statsCollector != nil {
			c.statsCollector.Add(workerID, stats.LLMErrors, 1)
		}
		return nil, client.EmptyCursor, err
	}

	response := make([]*types.LLMProcessorResult, 0, len(dataset.Data.Items))

	for i, item := range dataset.Data.Items {
		var resp types.LLMProcessorResult
		if err := json.Unmarshal(item, &resp); err != nil {
			logrus.Warnf("Failed to unmarshal llm result at index %d: %v", i, err)
			continue
		}
		response = append(response, &resp)
	}

	if c.statsCollector != nil {
		c.statsCollector.Add(workerID, stats.LLMProcessedItems, uint(len(response)))
	}

	return response, nextCursor, nil
}
