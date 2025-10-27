package jobs

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/tee-worker/v2/api/args"
	"github.com/masa-finance/tee-worker/v2/api/types"
	"github.com/masa-finance/tee-worker/v2/internal/config"
	"github.com/masa-finance/tee-worker/v2/internal/jobs/linkedinapify"
	"github.com/masa-finance/tee-worker/v2/internal/jobs/stats"
	"github.com/masa-finance/tee-worker/v2/pkg/client"

	pArgs "github.com/masa-finance/tee-worker/v2/api/args/linkedin/profile"
	pTypes "github.com/masa-finance/tee-worker/v2/api/types/linkedin/profile"
)

// LinkedInApifyClient defines the interface for the LinkedIn Apify client to allow mocking in tests
type LinkedInApifyClient interface {
	SearchProfiles(workerID string, args *pArgs.Arguments, cursor client.Cursor) ([]*pTypes.Profile, string, client.Cursor, error)
	ValidateApiKey() error
}

// NewLinkedInApifyClient is a function variable that can be replaced in tests.
// It defaults to the actual implementation.
var NewLinkedInApifyClient = func(apiKey string, statsCollector *stats.StatsCollector) (LinkedInApifyClient, error) {
	return linkedinapify.NewClient(apiKey, statsCollector)
}

type LinkedInScraper struct {
	configuration  config.JobConfiguration
	statsCollector *stats.StatsCollector
	capabilities   []types.Capability
}

func NewLinkedInScraper(jc config.JobConfiguration, statsCollector *stats.StatsCollector) *LinkedInScraper {
	logrus.Info("LinkedIn scraper via Apify initialized")
	return &LinkedInScraper{
		configuration:  jc,
		statsCollector: statsCollector,
		capabilities:   types.LinkedInCaps,
	}
}

func (ls *LinkedInScraper) ExecuteJob(j types.Job) (types.JobResult, error) {
	logrus.WithField("job_uuid", j.UUID).Info("Starting ExecuteJob for LinkedIn profile search")

	// Require Apify key for LinkedIn scraping
	apifyApiKey := ls.configuration.GetString("apify_api_key", "")
	if apifyApiKey == "" {
		msg := errors.New("apify API key is required for LinkedIn job")
		return types.JobResult{Error: msg.Error()}, msg
	}

	jobArgs, err := args.UnmarshalJobArguments(types.JobType(j.Type), map[string]any(j.Arguments))
	if err != nil {
		msg := fmt.Errorf("failed to unmarshal job arguments: %w", err)
		return types.JobResult{Error: msg.Error()}, msg
	}

	linkedinArgs, ok := jobArgs.(*pArgs.Arguments)
	if !ok {
		return types.JobResult{Error: "invalid argument type for LinkedIn job"}, errors.New("invalid argument type")
	}
	logrus.Debugf("LinkedIn job args: %+v", *linkedinArgs)

	linkedinClient, err := NewLinkedInApifyClient(apifyApiKey, ls.statsCollector)
	if err != nil {
		return types.JobResult{Error: "error while creating LinkedIn Apify client"}, fmt.Errorf("error creating LinkedIn Apify client: %w", err)
	}

	profiles, datasetId, cursor, err := linkedinClient.SearchProfiles(j.WorkerID, linkedinArgs, client.EmptyCursor)
	if err != nil {
		return types.JobResult{Error: fmt.Sprintf("error while searching LinkedIn profiles: %s", err.Error())}, fmt.Errorf("error searching LinkedIn profiles: %w", err)
	}

	if datasetId == "" {
		return types.JobResult{Error: "missing dataset id from LinkedIn profile search"}, errors.New("missing dataset id from LinkedIn profile search")
	}

	data, err := json.Marshal(profiles)
	if err != nil {
		return types.JobResult{Error: "error marshalling LinkedIn response"}, fmt.Errorf("error marshalling LinkedIn response: %w", err)
	}

	return types.JobResult{
		Data:       data,
		Job:        j,
		NextCursor: cursor.String(),
	}, nil
}
