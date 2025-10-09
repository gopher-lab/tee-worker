package linkedinapify_test

import (
	"encoding/json"
	"errors"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/internal/apify"
	"github.com/masa-finance/tee-worker/internal/jobs/linkedinapify"
	"github.com/masa-finance/tee-worker/internal/jobs/stats"
	"github.com/masa-finance/tee-worker/pkg/client"

	profileArgs "github.com/masa-finance/tee-types/args/linkedin/profile"
	"github.com/masa-finance/tee-types/types"
	"github.com/masa-finance/tee-types/types/linkedin/profile"
)

// MockApifyClient is a mock implementation of the ApifyClient.
type MockApifyClient struct {
	RunActorAndGetResponseFunc func(actorID apify.ActorId, input any, cursor client.Cursor, limit uint) (*client.DatasetResponse, client.Cursor, error)
	ValidateApiKeyFunc         func() error
	ProbeActorAccessFunc       func(actorID apify.ActorId, input map[string]any) (bool, error)
}

func (m *MockApifyClient) RunActorAndGetResponse(actorID apify.ActorId, input any, cursor client.Cursor, limit uint) (*client.DatasetResponse, client.Cursor, error) {
	if m.RunActorAndGetResponseFunc != nil {
		return m.RunActorAndGetResponseFunc(actorID, input, cursor, limit)
	}
	return nil, "", errors.New("RunActorAndGetResponseFunc not defined")
}

func (m *MockApifyClient) ValidateApiKey() error {
	if m.ValidateApiKeyFunc != nil {
		return m.ValidateApiKeyFunc()
	}
	return errors.New("ValidateApiKeyFunc not defined")
}

func (m *MockApifyClient) ProbeActorAccess(actorID apify.ActorId, input map[string]any) (bool, error) {
	if m.ProbeActorAccessFunc != nil {
		return m.ProbeActorAccessFunc(actorID, input)
	}
	return false, errors.New("ProbeActorAccessFunc not defined")
}

var _ = Describe("LinkedInApifyClient", func() {
	var (
		mockClient     *MockApifyClient
		linkedinClient *linkedinapify.ApifyClient
		statsCollector *stats.StatsCollector
		apifyKey       string
	)

	BeforeEach(func() {
		apifyKey = os.Getenv("APIFY_API_KEY")
		mockClient = &MockApifyClient{}
		statsCollector = stats.StartCollector(100, nil)

		// Replace the client creation function with one that returns the mock
		linkedinapify.NewInternalClient = func(apiKey string) (client.Apify, error) {
			return mockClient, nil
		}
		var err error
		linkedinClient, err = linkedinapify.NewClient("test-token", statsCollector)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("SearchProfiles", func() {
		It("should construct the correct actor input", func() {
			args := profileArgs.Arguments{
				Query:    "software engineer",
				MaxItems: 10,
			}

			mockClient.RunActorAndGetResponseFunc = func(actorID apify.ActorId, input any, cursor client.Cursor, limit uint) (*client.DatasetResponse, client.Cursor, error) {
				Expect(actorID).To(Equal(apify.ActorIds.LinkedInSearchProfile))
				Expect(limit).To(Equal(uint(10)))

				// Verify the input is correctly converted to map
				inputMap, ok := input.(map[string]any)
				Expect(ok).To(BeTrue())
				Expect(inputMap["searchQuery"]).To(Equal("software engineer"))

				return &client.DatasetResponse{Data: client.ApifyDatasetData{Items: []json.RawMessage{}}}, "next", nil
			}

			_, _, _, err := linkedinClient.SearchProfiles("test-worker", &args, client.EmptyCursor)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle errors from the apify client", func() {
			expectedErr := errors.New("apify error")
			mockClient.RunActorAndGetResponseFunc = func(actorID apify.ActorId, input any, cursor client.Cursor, limit uint) (*client.DatasetResponse, client.Cursor, error) {
				return nil, "", expectedErr
			}

			args := profileArgs.Arguments{
				Query:    "test query",
				MaxItems: 5,
			}
			_, _, _, err := linkedinClient.SearchProfiles("test-worker", &args, client.EmptyCursor)
			Expect(err).To(MatchError(expectedErr))
		})

		It("should handle JSON unmarshalling errors gracefully", func() {
			invalidJSON := []byte(`{"id": "test", "firstName": 123}`) // firstName should be a string
			dataset := &client.DatasetResponse{
				Data: client.ApifyDatasetData{
					Items: []json.RawMessage{invalidJSON},
				},
			}
			mockClient.RunActorAndGetResponseFunc = func(actorID apify.ActorId, input any, cursor client.Cursor, limit uint) (*client.DatasetResponse, client.Cursor, error) {
				return dataset, "next", nil
			}

			args := profileArgs.Arguments{
				Query:    "test query",
				MaxItems: 1,
			}
			results, _, _, err := linkedinClient.SearchProfiles("test-worker", &args, client.EmptyCursor)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(BeEmpty()) // The invalid item should be skipped
		})

		It("should handle multiple valid profiles", func() {
			profile1, _ := json.Marshal(map[string]any{
				"id":        "profile-1",
				"firstName": "John",
				"lastName":  "Doe",
				"headline":  "Software Engineer",
			})
			profile2, _ := json.Marshal(map[string]any{
				"id":        "profile-2",
				"firstName": "Jane",
				"lastName":  "Smith",
				"headline":  "Product Manager",
			})
			dataset := &client.DatasetResponse{
				Data: client.ApifyDatasetData{
					Items: []json.RawMessage{profile1, profile2},
				},
			}
			mockClient.RunActorAndGetResponseFunc = func(actorID apify.ActorId, input any, cursor client.Cursor, limit uint) (*client.DatasetResponse, client.Cursor, error) {
				return dataset, "next", nil
			}

			args := profileArgs.Arguments{
				Query:    "test query",
				MaxItems: 2,
			}
			results, _, _, err := linkedinClient.SearchProfiles("test-worker", &args, client.EmptyCursor)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(2))
			Expect(results[0].FirstName).To(Equal("John"))
			Expect(results[1].FirstName).To(Equal("Jane"))
		})
	})

	Describe("ValidateApiKey", func() {
		It("should validate the API key", func() {
			mockClient.ValidateApiKeyFunc = func() error {
				return nil
			}
			Expect(linkedinClient.ValidateApiKey()).To(Succeed())
		})

		It("should return error when validation fails", func() {
			expectedErr := errors.New("invalid key")
			mockClient.ValidateApiKeyFunc = func() error {
				return expectedErr
			}
			Expect(linkedinClient.ValidateApiKey()).To(MatchError(expectedErr))
		})
	})

	// Integration tests that use the real client
	Context("Integration tests", func() {
		It("should validate API key with real client when APIFY_API_KEY is set", func() {
			if apifyKey == "" {
				Skip("APIFY_API_KEY required to run LinkedIn integration tests")
			}

			// Reset to use real client
			linkedinapify.NewInternalClient = func(apiKey string) (client.Apify, error) {
				return client.NewApifyClient(apiKey)
			}

			realClient, err := linkedinapify.NewClient(apifyKey, statsCollector)
			Expect(err).NotTo(HaveOccurred())
			Expect(realClient.ValidateApiKey()).To(Succeed())
		})

		It("should search profiles with real client when APIFY_API_KEY is set", func() {
			if apifyKey == "" {
				Skip("APIFY_API_KEY required to run LinkedIn integration tests")
			}

			// Reset to use real client
			linkedinapify.NewInternalClient = func(apiKey string) (client.Apify, error) {
				return client.NewApifyClient(apiKey)
			}

			realClient, err := linkedinapify.NewClient(apifyKey, statsCollector)
			Expect(err).NotTo(HaveOccurred())

			args := profileArgs.Arguments{
				QueryType:   types.CapSearchByProfile,
				Query:       "software engineer",
				MaxItems:    1,
				ScraperMode: profile.ScraperModeShort,
			}

			results, datasetId, cursor, err := realClient.SearchProfiles("test-worker", &args, client.EmptyCursor)
			Expect(err).NotTo(HaveOccurred())
			Expect(datasetId).NotTo(BeEmpty())
			Expect(results).NotTo(BeEmpty())
			Expect(results[0]).NotTo(BeNil())
			Expect(results[0].ID).NotTo(BeEmpty())
			Expect(cursor).NotTo(BeEmpty())
		})
	})
})
