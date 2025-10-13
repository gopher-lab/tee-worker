package capabilities_test

import (
	"os"
	"slices"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/masa-finance/tee-worker/internal/capabilities"
	"github.com/masa-finance/tee-worker/internal/config"
)

// MockJobServer implements JobServerInterface for testing
type MockJobServer struct {
	capabilities types.WorkerCapabilities
}

func (m *MockJobServer) GetWorkerCapabilities() types.WorkerCapabilities {
	return m.capabilities
}

var _ = Describe("DetectCapabilities", func() {
	DescribeTable("capability detection scenarios",
		func(jc config.JobConfiguration, jobServer JobServerInterface, expected types.WorkerCapabilities) {
			got := DetectCapabilities(jc, jobServer)

			// Extract job type keys and sort for consistent comparison
			gotKeys := make([]string, 0, len(got))
			for jobType := range got {
				gotKeys = append(gotKeys, jobType.String())
			}

			expectedKeys := make([]string, 0, len(expected))
			for jobType := range expected {
				expectedKeys = append(expectedKeys, jobType.String())
			}

			// Sort both slices for comparison
			slices.Sort(gotKeys)
			slices.Sort(expectedKeys)

			// Compare the sorted slices
			Expect(gotKeys).To(Equal(expectedKeys))
		},
		Entry("With JobServer - performs real detection (JobServer ignored)",
			config.JobConfiguration{},
			&MockJobServer{
				capabilities: types.WorkerCapabilities{
					types.WebJob:       {types.CapScraper},
					types.TelemetryJob: {types.CapTelemetry},
					types.TiktokJob:    {types.CapTranscription},
					types.TwitterJob:   {types.CapSearchByQuery, types.CapGetById, types.CapGetProfileById},
				},
			},
			types.WorkerCapabilities{
				types.TelemetryJob: {types.CapTelemetry},
				types.TiktokJob:    {types.CapTranscription},
			},
		),
		Entry("Without JobServer - basic capabilities only",
			config.JobConfiguration{},
			nil,
			types.WorkerCapabilities{
				types.TelemetryJob: {types.CapTelemetry},
				types.TiktokJob:    {types.CapTranscription},
			},
		),
		Entry("With Twitter accounts - adds credential capabilities",
			config.JobConfiguration{
				"twitter_accounts": []string{"account1", "account2"},
			},
			nil,
			types.WorkerCapabilities{
				types.TelemetryJob:         {types.CapTelemetry},
				types.TiktokJob:            {types.CapTranscription},
				types.TwitterCredentialJob: types.TwitterCredentialCaps,
				types.TwitterJob:           types.TwitterCredentialCaps,
			},
		),
		Entry("With Twitter API keys - adds API capabilities",
			config.JobConfiguration{
				"twitter_api_keys": []string{"key1", "key2"},
			},
			nil,
			types.WorkerCapabilities{
				types.TelemetryJob:  {types.CapTelemetry},
				types.TiktokJob:     {types.CapTranscription},
				types.TwitterApiJob: types.TwitterAPICaps,
				types.TwitterJob:    types.TwitterAPICaps,
			},
		),
		Entry("With mock elevated Twitter API keys - only basic capabilities detected",
			config.JobConfiguration{
				"twitter_api_keys": []string{"Bearer abcd1234-ELEVATED"},
			},
			nil,
			types.WorkerCapabilities{
				types.TelemetryJob: {types.CapTelemetry},
				types.TiktokJob:    {types.CapTranscription},
				// Note: Mock elevated keys will be detected as basic since we can't make real API calls in tests
				types.TwitterApiJob: types.TwitterAPICaps,
				types.TwitterJob:    types.TwitterAPICaps,
			},
		),
	)

	Context("Scraper Types", func() {
		DescribeTable("scraper type detection",
			func(jc config.JobConfiguration, expectedKeys []string) {
				caps := DetectCapabilities(jc, nil)

				jobNames := make([]string, 0, len(caps))
				for jobType := range caps {
					jobNames = append(jobNames, jobType.String())
				}

				// Sort both slices for comparison
				slices.Sort(jobNames)
				expectedSorted := make([]string, len(expectedKeys))
				copy(expectedSorted, expectedKeys)
				slices.Sort(expectedSorted)

				// Compare the sorted slices
				Expect(jobNames).To(Equal(expectedSorted))
			},
			Entry("Basic scrapers only",
				config.JobConfiguration{},
				[]string{"telemetry", "tiktok"},
			),
			Entry("With Twitter accounts",
				config.JobConfiguration{
					"twitter_accounts": []string{"user1:pass1"},
				},
				[]string{"telemetry", "tiktok", "twitter", "twitter-credential"},
			),
			Entry("With Twitter API keys",
				config.JobConfiguration{
					"twitter_api_keys": []string{"key1"},
				},
				[]string{"telemetry", "tiktok", "twitter", "twitter-api"},
			),
		)
	})

	Context("Apify Integration", func() {
		It("should add enhanced capabilities when valid Apify API key is provided", func() {
			apifyKey := os.Getenv("APIFY_API_KEY")
			if apifyKey == "" {
				Skip("APIFY_API_KEY is not set")
			}

			jc := config.JobConfiguration{
				"apify_api_key": apifyKey,
			}

			caps := DetectCapabilities(jc, nil)

			// TikTok should gain search capabilities with valid key
			tiktokCaps, ok := caps[types.TiktokJob]
			Expect(ok).To(BeTrue(), "expected tiktok capabilities to be present")
			Expect(tiktokCaps).To(ContainElement(types.CapSearchByQuery), "expected tiktok to include CapSearchByQuery capability")
			Expect(tiktokCaps).To(ContainElement(types.CapSearchByTrending), "expected tiktok to include CapSearchByTrending capability")

			// Twitter-Apify job should be present with follower/following capabilities
			twitterApifyCaps, ok := caps[types.TwitterApifyJob]
			Expect(ok).To(BeTrue(), "expected twitter-apify capabilities to be present")
			Expect(twitterApifyCaps).To(ContainElement(types.CapGetFollowers), "expected twitter-apify to include CapGetFollowers capability")
			Expect(twitterApifyCaps).To(ContainElement(types.CapGetFollowing), "expected twitter-apify to include CapGetFollowing capability")

			// Reddit should be present (only if rented!)
			redditCaps, hasReddit := caps[types.RedditJob]
			Expect(hasReddit).To(BeTrue(), "expected reddit capabilities to be present")
			Expect(redditCaps).To(ContainElement(types.CapScrapeUrls), "expected reddit to include CapScrapeUrls capability")
			Expect(redditCaps).To(ContainElement(types.CapSearchPosts), "expected reddit to include CapSearchPosts capability")
			Expect(redditCaps).To(ContainElement(types.CapSearchUsers), "expected reddit to include CapSearchUsers capability")
			Expect(redditCaps).To(ContainElement(types.CapSearchCommunities), "expected reddit to include CapSearchCommunities capability")
		})
		It("should add enhanced capabilities when valid Apify API key is provided alongside a Gemini API key", func() {
			apifyKey := os.Getenv("APIFY_API_KEY")
			if apifyKey == "" {
				Skip("APIFY_API_KEY is not set")
			}

			geminiKey := os.Getenv("GEMINI_API_KEY")
			if geminiKey == "" {
				Skip("GEMINI_API_KEY is not set")
			}

			jc := config.JobConfiguration{
				"apify_api_key":  apifyKey,
				"gemini_api_key": geminiKey,
			}
			caps := DetectCapabilities(jc, nil)

			// Web should be present
			webCaps, hasWeb := caps[types.WebJob]
			Expect(hasWeb).To(BeTrue(), "expected web capabilities to be present")
			Expect(webCaps).To(ContainElement(types.CapScraper), "expected web to include CapScraper capability")
		})
	})
})

// Helper function to check if a job type exists in capabilities
func hasJobType(capabilities types.WorkerCapabilities, jobName string) bool {
	_, exists := capabilities[types.JobType(jobName)]
	return exists
}
