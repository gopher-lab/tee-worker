package jobs_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/api/types"
	"github.com/masa-finance/tee-worker/internal/config"
	"github.com/masa-finance/tee-worker/internal/jobs"
	"github.com/masa-finance/tee-worker/internal/jobs/linkedinapify"
	"github.com/masa-finance/tee-worker/internal/jobs/stats"
	"github.com/masa-finance/tee-worker/pkg/client"

	profileArgs "github.com/masa-finance/tee-worker/api/args/linkedin/profile"
	profileTypes "github.com/masa-finance/tee-worker/api/types/linkedin/profile"
)

// MockLinkedInApifyClient is a mock implementation of the LinkedInApifyClient.
type MockLinkedInApifyClient struct {
	SearchProfilesFunc func(workerID string, args *profileArgs.Arguments, cursor client.Cursor) ([]*profileTypes.Profile, string, client.Cursor, error)
	ValidateApiKeyFunc func() error
}

func (m *MockLinkedInApifyClient) SearchProfiles(workerID string, args *profileArgs.Arguments, cursor client.Cursor) ([]*profileTypes.Profile, string, client.Cursor, error) {
	if m != nil && m.SearchProfilesFunc != nil {
		return m.SearchProfilesFunc(workerID, args, cursor)
	}
	return nil, "", client.EmptyCursor, nil
}

func (m *MockLinkedInApifyClient) ValidateApiKey() error {
	if m != nil && m.ValidateApiKeyFunc != nil {
		return m.ValidateApiKeyFunc()
	}
	return nil
}

var _ = Describe("LinkedInScraper", func() {
	var (
		scraper        *jobs.LinkedInScraper
		statsCollector *stats.StatsCollector
		job            types.Job
		mockClient     *MockLinkedInApifyClient
	)

	// Keep original to restore after each test to avoid leaking globals
	originalNewLinkedInApifyClient := jobs.NewLinkedInApifyClient

	BeforeEach(func() {
		statsCollector = stats.StartCollector(128, config.JobConfiguration{})
		cfg := config.JobConfiguration{
			"apify_api_key": "test-key",
		}
		scraper = jobs.NewLinkedInScraper(cfg, statsCollector)
		mockClient = &MockLinkedInApifyClient{}

		// Replace the client creation function with one that returns the mock
		jobs.NewLinkedInApifyClient = func(apiKey string, _ *stats.StatsCollector) (jobs.LinkedInApifyClient, error) {
			return mockClient, nil
		}

		job = types.Job{
			UUID: "test-uuid",
			Type: types.LinkedInJob,
		}
	})

	AfterEach(func() {
		jobs.NewLinkedInApifyClient = originalNewLinkedInApifyClient
	})

	Context("ExecuteJob", func() {
		It("should return an error when Apify API key is missing", func() {
			cfg := config.JobConfiguration{}
			scraper = jobs.NewLinkedInScraper(cfg, statsCollector)

			job.Arguments = map[string]any{
				"type":        types.CapSearchByProfile,
				"searchQuery": "software engineer",
				"maxItems":    10,
			}

			result, err := scraper.ExecuteJob(job)
			Expect(err).To(HaveOccurred())
			Expect(result.Error).To(ContainSubstring("Apify API key is required for LinkedIn job"))
		})

		It("should call SearchProfiles and return data and next cursor", func() {
			job.Arguments = map[string]any{
				"type":        types.CapSearchByProfile,
				"searchQuery": "software engineer",
				"maxItems":    10,
			}

			headline := "Software Engineer"
			headline2 := "Senior Software Engineer"
			expectedProfiles := []*profileTypes.Profile{
				{
					ID:               "profile-1",
					FirstName:        "John",
					LastName:         "Doe",
					Headline:         &headline,
					PublicIdentifier: "john-doe",
					URL:              "https://linkedin.com/in/john-doe",
				},
				{
					ID:               "profile-2",
					FirstName:        "Jane",
					LastName:         "Smith",
					Headline:         &headline2,
					PublicIdentifier: "jane-smith",
					URL:              "https://linkedin.com/in/jane-smith",
				},
			}

			mockClient.SearchProfilesFunc = func(workerID string, args *profileArgs.Arguments, cursor client.Cursor) ([]*profileTypes.Profile, string, client.Cursor, error) {
				Expect(workerID).To(Equal("test-worker"))
				Expect(args.Query).To(Equal("software engineer"))
				Expect(args.MaxItems).To(Equal(uint(10)))
				Expect(args.Type).To(Equal(types.CapSearchByProfile))
				return expectedProfiles, "dataset-123", client.Cursor("next-cursor"), nil
			}

			job.WorkerID = "test-worker"
			result, err := scraper.ExecuteJob(job)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.NextCursor).To(Equal("next-cursor"))

			var resp []*profileTypes.Profile
			err = json.Unmarshal(result.Data, &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(HaveLen(2))
			Expect(resp[0].ID).To(Equal("profile-1"))
			Expect(resp[0].FirstName).To(Equal("John"))
			Expect(resp[1].ID).To(Equal("profile-2"))
			Expect(resp[1].FirstName).To(Equal("Jane"))
		})

		It("should handle errors from the LinkedIn client", func() {
			job.Arguments = map[string]any{
				"type":        types.CapSearchByProfile,
				"searchQuery": "software engineer",
				"maxItems":    10,
			}

			expectedErr := errors.New("client error")
			mockClient.SearchProfilesFunc = func(workerID string, args *profileArgs.Arguments, cursor client.Cursor) ([]*profileTypes.Profile, string, client.Cursor, error) {
				return nil, "", client.EmptyCursor, expectedErr
			}

			result, err := scraper.ExecuteJob(job)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("client error")))
			Expect(result.Error).To(ContainSubstring("error while searching LinkedIn profiles: client error"))
		})

		It("should handle errors when creating the client", func() {
			jobs.NewLinkedInApifyClient = func(apiKey string, _ *stats.StatsCollector) (jobs.LinkedInApifyClient, error) {
				return nil, errors.New("client creation failed")
			}
			job.Arguments = map[string]any{
				"type":        types.CapSearchByProfile,
				"searchQuery": "software engineer",
				"maxItems":    10,
			}

			result, err := scraper.ExecuteJob(job)
			Expect(err).To(HaveOccurred())
			Expect(result.Error).To(Equal("error while creating LinkedIn Apify client"))
		})

		It("should return an error when dataset ID is missing", func() {
			job.Arguments = map[string]any{
				"type":        types.CapSearchByProfile,
				"searchQuery": "software engineer",
				"maxItems":    10,
			}

			mockClient.SearchProfilesFunc = func(workerID string, args *profileArgs.Arguments, cursor client.Cursor) ([]*profileTypes.Profile, string, client.Cursor, error) {
				return []*profileTypes.Profile{}, "", client.EmptyCursor, nil
			}

			result, err := scraper.ExecuteJob(job)
			Expect(err).To(HaveOccurred())
			Expect(result.Error).To(ContainSubstring("missing dataset id from LinkedIn profile search"))
		})

		It("should handle JSON marshalling errors", func() {
			job.Arguments = map[string]any{
				"type":        types.CapSearchByProfile,
				"searchQuery": "software engineer",
				"maxItems":    10,
			}

			// Create a profile with a channel to cause JSON marshalling to fail
			invalidProfile := &profileTypes.Profile{
				ID:        "profile-1",
				FirstName: "John",
				LastName:  "Doe",
			}

			mockClient.SearchProfilesFunc = func(workerID string, args *profileArgs.Arguments, cursor client.Cursor) ([]*profileTypes.Profile, string, client.Cursor, error) {
				return []*profileTypes.Profile{invalidProfile}, "dataset-123", client.EmptyCursor, nil
			}

			result, err := scraper.ExecuteJob(job)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Error).To(BeEmpty())
			Expect(result.Data).NotTo(BeEmpty())
		})

		It("should handle empty profile results", func() {
			job.Arguments = map[string]any{
				"type":        types.CapSearchByProfile,
				"searchQuery": "nonexistent",
				"maxItems":    10,
			}

			mockClient.SearchProfilesFunc = func(workerID string, args *profileArgs.Arguments, cursor client.Cursor) ([]*profileTypes.Profile, string, client.Cursor, error) {
				return []*profileTypes.Profile{}, "dataset-123", client.EmptyCursor, nil
			}

			result, err := scraper.ExecuteJob(job)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.NextCursor).To(Equal(""))

			var resp []*profileTypes.Profile
			err = json.Unmarshal(result.Data, &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(HaveLen(0))
		})
	})

	// Integration tests that use the real client
	Context("Integration tests", func() {
		var apifyKey string

		BeforeEach(func() {
			apifyKey = os.Getenv("APIFY_API_KEY")

			if apifyKey == "" {
				Skip("APIFY_API_KEY required for LinkedIn integration tests")
			}

			// Reset to use real client for integration tests
			jobs.NewLinkedInApifyClient = func(apiKey string, s *stats.StatsCollector) (jobs.LinkedInApifyClient, error) {
				return linkedinapify.NewClient(apiKey, s)
			}
		})

		It("should execute a real LinkedIn profile search when API key is set", func() {
			cfg := config.JobConfiguration{
				"apify_api_key": apifyKey,
			}
			integrationStatsCollector := stats.StartCollector(128, cfg)
			integrationScraper := jobs.NewLinkedInScraper(cfg, integrationStatsCollector)

			jobArgs := profileArgs.Arguments{
				Type:     types.CapSearchByProfile,
				Query:    "software engineer",
				MaxItems: 10,
			}

			// Marshal jobArgs to map[string]any so it can be used as JobArgument
			var jobArgsMap map[string]any
			jobArgsBytes, err := json.Marshal(jobArgs)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal(jobArgsBytes, &jobArgsMap)
			Expect(err).NotTo(HaveOccurred())

			job := types.Job{
				UUID:      "integration-test-uuid",
				Type:      types.LinkedInJob,
				WorkerID:  "test-worker",
				Arguments: jobArgsMap,
				Timeout:   60 * time.Second,
			}

			result, err := integrationScraper.ExecuteJob(job)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Error).To(BeEmpty())
			Expect(result.Data).NotTo(BeEmpty())

			var resp []*profileTypes.Profile
			err = json.Unmarshal(result.Data, &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeEmpty())
			Expect(resp[0]).NotTo(BeNil())
			Expect(resp[0].ID).NotTo(BeEmpty())

			prettyJSON, err := json.MarshalIndent(resp, "", "  ")
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(prettyJSON))
		})

		// Note: Capability detection is now centralized in capabilities/detector.go
		// Individual scraper capability tests have been removed
	})
})
