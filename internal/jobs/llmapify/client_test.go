package llmapify_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/v2/api/args/llm/process"
	"github.com/masa-finance/tee-worker/v2/api/types"
	"github.com/masa-finance/tee-worker/v2/internal/apify"
	"github.com/masa-finance/tee-worker/v2/internal/config"
	"github.com/masa-finance/tee-worker/v2/internal/jobs/llmapify"
	"github.com/masa-finance/tee-worker/v2/pkg/client"
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

var _ = Describe("LLMApifyClient", func() {
	var (
		mockClient *MockApifyClient
		llmClient  *llmapify.ApifyClient
		apifyKey   string
	)

	BeforeEach(func() {
		apifyKey = os.Getenv("APIFY_API_KEY")
		mockClient = &MockApifyClient{}
		// Replace the client creation function with one that returns the mock
		llmapify.NewInternalClient = func(apiKey string) (client.Apify, error) {
			return mockClient, nil
		}
		var err error
		llmClient, err = llmapify.NewClient("test-token", config.LlmConfig{ClaudeApiKey: "test-claude-llm-key"}, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Process", func() {
		It("should construct the correct actor input", func() {
			llmArgs := process.NewArguments()
			llmArgs.DatasetId = "test-dataset-id"
			llmArgs.Prompt = "test-prompt"

			// Marshal and unmarshal to apply defaults
			jsonData, err := json.Marshal(llmArgs)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal(jsonData, &llmArgs)
			Expect(err).ToNot(HaveOccurred())

			mockClient.RunActorAndGetResponseFunc = func(actorID apify.ActorId, input any, cursor client.Cursor, limit uint) (*client.DatasetResponse, client.Cursor, error) {
				Expect(actorID).To(Equal(apify.ActorIds.LLMDatasetProcessor))
				Expect(limit).To(Equal(uint(1)))

				// Verify the input is correctly converted to LLMProcessorRequest
				request, ok := input.(types.LLMProcessorRequest)
				Expect(ok).To(BeTrue())
				Expect(request.InputDatasetId).To(Equal("test-dataset-id"))
				Expect(request.Prompt).To(Equal("test-prompt"))
				Expect(request.LLMProviderApiKey).To(Equal("test-claude-llm-key"))                                  // should be set from constructor
				Expect(request.Model).To(Equal(process.DefaultClaudeModel))                                         // default model
				Expect(request.MultipleColumns).To(Equal(process.DefaultMultipleColumns))                           // default value
				Expect(request.MaxTokens).To(Equal(process.DefaultMaxTokens))                                       // default value
				Expect(request.Temperature).To(Equal(strconv.FormatFloat(process.DefaultTemperature, 'f', -1, 64))) // default value

				return &client.DatasetResponse{Data: client.ApifyDatasetData{Items: []json.RawMessage{}}}, "next", nil
			}

			_, _, processErr := llmClient.Process("test-worker", llmArgs, client.EmptyCursor)
			Expect(processErr).NotTo(HaveOccurred())
		})

		It("should handle errors from the apify client", func() {
			expectedErr := errors.New("apify error")
			mockClient.RunActorAndGetResponseFunc = func(actorID apify.ActorId, input any, cursor client.Cursor, limit uint) (*client.DatasetResponse, client.Cursor, error) {
				return nil, "", expectedErr
			}

			llmArgs := process.NewArguments()
			llmArgs.DatasetId = "test-dataset-id"
			llmArgs.Prompt = "test-prompt"
			_, _, err := llmClient.Process("test-worker", llmArgs, client.EmptyCursor)
			Expect(err).To(MatchError(expectedErr))
		})

		It("should handle JSON unmarshalling errors gracefully", func() {
			invalidJSON := []byte(`{"llmresponse": 123}`) // llmresponse should be a string
			dataset := &client.DatasetResponse{
				Data: client.ApifyDatasetData{
					Items: []json.RawMessage{invalidJSON},
				},
			}
			mockClient.RunActorAndGetResponseFunc = func(actorID apify.ActorId, input any, cursor client.Cursor, limit uint) (*client.DatasetResponse, client.Cursor, error) {
				return dataset, "next", nil
			}

			llmArgs := process.NewArguments()
			llmArgs.DatasetId = "test-dataset-id"
			llmArgs.Prompt = "test-prompt"
			results, _, err := llmClient.Process("test-worker", llmArgs, client.EmptyCursor)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(BeEmpty()) // The invalid item should be skipped
		})

		It("should correctly unmarshal valid items", func() {
			llmResultJSON, _ := json.Marshal(map[string]any{
				"llmresponse": "This is a summary of the webpage content.",
			})
			dataset := &client.DatasetResponse{
				Data: client.ApifyDatasetData{
					Items: []json.RawMessage{llmResultJSON},
				},
			}
			mockClient.RunActorAndGetResponseFunc = func(actorID apify.ActorId, input any, cursor client.Cursor, limit uint) (*client.DatasetResponse, client.Cursor, error) {
				return dataset, "next", nil
			}

			llmArgs := process.NewArguments()
			llmArgs.DatasetId = "test-dataset-id"
			llmArgs.Prompt = "test-prompt"
			results, cursor, err := llmClient.Process("test-worker", llmArgs, client.EmptyCursor)
			Expect(err).NotTo(HaveOccurred())
			Expect(cursor).To(Equal(client.Cursor("next")))
			Expect(results).To(HaveLen(1))
			Expect(results[0].LLMResponse).To(Equal("This is a summary of the webpage content."))
		})

		It("should handle multiple valid results", func() {
			llmResult1, _ := json.Marshal(map[string]any{
				"llmresponse": "First summary.",
			})
			llmResult2, _ := json.Marshal(map[string]any{
				"llmresponse": "Second summary.",
			})
			dataset := &client.DatasetResponse{
				Data: client.ApifyDatasetData{
					Items: []json.RawMessage{llmResult1, llmResult2},
				},
			}
			mockClient.RunActorAndGetResponseFunc = func(actorID apify.ActorId, input any, cursor client.Cursor, limit uint) (*client.DatasetResponse, client.Cursor, error) {
				return dataset, "next", nil
			}

			llmArgs := process.NewArguments()
			llmArgs.DatasetId = "test-dataset-id"
			llmArgs.Prompt = "test-prompt"
			results, _, err := llmClient.Process("test-worker", llmArgs, client.EmptyCursor)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(2))
			Expect(results[0].LLMResponse).To(Equal("First summary."))
			Expect(results[1].LLMResponse).To(Equal("Second summary."))
		})

		It("should use custom values when provided", func() {
			llmArgs := process.NewArguments()
			llmArgs.DatasetId = "test-dataset-id"
			llmArgs.Prompt = "test-prompt"
			llmArgs.MaxTokens = 500
			llmArgs.Temperature = 0.5

			mockClient.RunActorAndGetResponseFunc = func(actorID apify.ActorId, input any, cursor client.Cursor, limit uint) (*client.DatasetResponse, client.Cursor, error) {
				request, ok := input.(types.LLMProcessorRequest)
				Expect(ok).To(BeTrue())
				Expect(request.MaxTokens).To(Equal(uint(500)))
				Expect(request.Temperature).To(Equal("0.5"))
				Expect(request.LLMProviderApiKey).To(Equal("test-claude-llm-key")) // should be set from constructor

				return &client.DatasetResponse{Data: client.ApifyDatasetData{Items: []json.RawMessage{}}}, "next", nil
			}

			_, _, err := llmClient.Process("test-worker", llmArgs, client.EmptyCursor)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("ValidateApiKey", func() {
		It("should validate the API key", func() {
			mockClient.ValidateApiKeyFunc = func() error {
				return nil
			}
			Expect(llmClient.ValidateApiKey()).To(Succeed())
		})

		It("should return error when validation fails", func() {
			expectedErr := errors.New("invalid key")
			mockClient.ValidateApiKeyFunc = func() error {
				return expectedErr
			}
			Expect(llmClient.ValidateApiKey()).To(MatchError(expectedErr))
		})
	})

	// Integration tests that use the real client
	Context("Integration tests", func() {
		It("should validate API key with real client when both APIFY_API_KEY and GEMINI_API_KEY are set", func() {
			geminiKey := os.Getenv("GEMINI_API_KEY")
			if apifyKey == "" || geminiKey == "" {
				Skip("Both APIFY_API_KEY and GEMINI_API_KEY must be set for integration tests")
			}

			// Reset to use real client
			llmapify.NewInternalClient = func(apiKey string) (client.Apify, error) {
				return client.NewApifyClient(apiKey)
			}

			realClient, err := llmapify.NewClient(apifyKey, config.LlmConfig{GeminiApiKey: config.LlmApiKey(geminiKey)}, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(realClient.ValidateApiKey()).To(Succeed())
		})

		It("should process a real dataset when both APIFY_API_KEY and GEMINI_API_KEY are set", func() {
			geminiKey := os.Getenv("GEMINI_API_KEY")
			if apifyKey == "" || geminiKey == "" {
				Skip("Both APIFY_API_KEY and GEMINI_API_KEY must be set for integration tests")
			}

			// Reset to use real client
			llmapify.NewInternalClient = func(apiKey string) (client.Apify, error) {
				return client.NewApifyClient(apiKey)
			}

			realClient, err := llmapify.NewClient(apifyKey, config.LlmConfig{GeminiApiKey: config.LlmApiKey(geminiKey)}, nil)
			Expect(err).NotTo(HaveOccurred())

			llmArgs := process.NewArguments()
			llmArgs.DatasetId = "V6tyuuZIgfiETl1cl"
			llmArgs.Prompt = "summarize the content of this webpage ${markdown}"
			// Marshal and unmarshal to apply defaults
			jsonData, err := json.Marshal(llmArgs)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal(jsonData, &llmArgs)
			Expect(err).ToNot(HaveOccurred())

			results, cursor, err := realClient.Process("test-worker", llmArgs, client.EmptyCursor)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			Expect(results[0]).NotTo(BeNil())
			Expect(results[0].LLMResponse).NotTo(BeEmpty())
			Expect(cursor).NotTo(BeEmpty())

			prettyJSON, err := json.MarshalIndent(results, "", "  ")
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(prettyJSON))
		})
	})
})
