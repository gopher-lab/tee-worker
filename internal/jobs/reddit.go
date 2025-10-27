package jobs

import (
	"encoding/json"
	"errors"
	"fmt"

	"time"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/tee-worker/v2/api/args"
	"github.com/masa-finance/tee-worker/v2/api/args/reddit/search"
	"github.com/masa-finance/tee-worker/v2/api/types"
	"github.com/masa-finance/tee-worker/v2/internal/config"
	"github.com/masa-finance/tee-worker/v2/internal/jobs/redditapify"
	"github.com/masa-finance/tee-worker/v2/internal/jobs/stats"
	"github.com/masa-finance/tee-worker/v2/pkg/client"
)

// RedditApifyClient defines the interface for the Reddit Apify client.
// This allows for mocking in tests.
type RedditApifyClient interface {
	ScrapeUrls(workerID string, urls []types.RedditStartURL, after time.Time, args redditapify.CommonArgs, cursor client.Cursor, maxResults uint) ([]*types.RedditResponse, client.Cursor, error)
	SearchPosts(workerID string, queries []string, after time.Time, args redditapify.CommonArgs, cursor client.Cursor, maxResults uint) ([]*types.RedditResponse, client.Cursor, error)
	SearchCommunities(workerID string, queries []string, args redditapify.CommonArgs, cursor client.Cursor, maxResults uint) ([]*types.RedditResponse, client.Cursor, error)
	SearchUsers(workerID string, queries []string, skipPosts bool, args redditapify.CommonArgs, cursor client.Cursor, maxResults uint) ([]*types.RedditResponse, client.Cursor, error)
}

// NewRedditApifyClient is a function variable that can be replaced in tests.
// It defaults to the actual implementation.
var NewRedditApifyClient = func(apiKey string, statsCollector *stats.StatsCollector) (RedditApifyClient, error) {
	return redditapify.NewClient(apiKey, statsCollector)
}

type RedditScraper struct {
	configuration  config.RedditConfig
	statsCollector *stats.StatsCollector
	capabilities   []types.Capability
}

func NewRedditScraper(jc config.JobConfiguration, statsCollector *stats.StatsCollector) *RedditScraper {
	config := jc.GetRedditConfig()
	logrus.Info("Reddit scraper via Apify initialized")
	return &RedditScraper{
		configuration:  config,
		statsCollector: statsCollector,
		capabilities:   types.RedditCaps,
	}
}

func (r *RedditScraper) ExecuteJob(j types.Job) (types.JobResult, error) {
	logrus.WithField("job_uuid", j.UUID).Info("Starting ExecuteJob for Reddit scrape")

	jobArgs, err := args.UnmarshalJobArguments(types.JobType(j.Type), map[string]any(j.Arguments))
	if err != nil {
		msg := fmt.Errorf("failed to unmarshal job arguments: %w", err)
		return types.JobResult{Error: msg.Error()}, msg
	}

	// Type assert to Reddit arguments
	redditArgs, ok := jobArgs.(*search.Arguments)
	if !ok {
		return types.JobResult{Error: "invalid argument type for Reddit job"}, errors.New("invalid argument type")
	}
	logrus.Debugf("reddit job args: %+v", *redditArgs)

	redditClient, err := NewRedditApifyClient(r.configuration.ApifyApiKey, r.statsCollector)
	if err != nil {
		return types.JobResult{Error: "error while scraping Reddit"}, fmt.Errorf("error creating Reddit Apify client: %w", err)
	}

	commonArgs := redditapify.CommonArgs{}
	commonArgs.CopyFromArgs(redditArgs)

	switch redditArgs.Type {
	case types.CapScrapeUrls:
		urls := make([]types.RedditStartURL, 0, len(redditArgs.URLs))
		for _, u := range redditArgs.URLs {
			urls = append(urls, types.RedditStartURL{
				URL:    u,
				Method: "GET",
			})
		}

		resp, cursor, err := redditClient.ScrapeUrls(j.WorkerID, urls, redditArgs.After, commonArgs, client.Cursor(redditArgs.NextCursor), redditArgs.MaxResults)
		return processRedditResponse(j, resp, cursor, err)

	case types.CapSearchUsers:
		resp, cursor, err := redditClient.SearchUsers(j.WorkerID, redditArgs.Queries, redditArgs.SkipPosts, commonArgs, client.Cursor(redditArgs.NextCursor), redditArgs.MaxResults)
		return processRedditResponse(j, resp, cursor, err)

	case types.CapSearchPosts:
		resp, cursor, err := redditClient.SearchPosts(j.WorkerID, redditArgs.Queries, redditArgs.After, commonArgs, client.Cursor(redditArgs.NextCursor), redditArgs.MaxResults)
		return processRedditResponse(j, resp, cursor, err)

	case types.CapSearchCommunities:
		resp, cursor, err := redditClient.SearchCommunities(j.WorkerID, redditArgs.Queries, commonArgs, client.Cursor(redditArgs.NextCursor), redditArgs.MaxResults)
		return processRedditResponse(j, resp, cursor, err)

	default:
		return types.JobResult{Error: "invalid type for Reddit job"}, fmt.Errorf("invalid type for Reddit job: %s", redditArgs.Type)
	}
}

func processRedditResponse(j types.Job, resp []*types.RedditResponse, cursor client.Cursor, err error) (types.JobResult, error) {
	if err != nil {
		return types.JobResult{Error: fmt.Sprintf("error while scraping Reddit: %s", err.Error())}, fmt.Errorf("error scraping Reddit: %w", err)
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return types.JobResult{Error: "error marshalling Reddit response"}, fmt.Errorf("error marshalling Reddit response: %w", err)
	}
	return types.JobResult{
		Data:       data,
		Job:        j,
		NextCursor: cursor.String(),
	}, nil
}
