package args_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-types/args"
)

var _ = Describe("LLMProcessorArguments", func() {
	Describe("Marshalling and unmarshalling", func() {
		It("should set default values", func() {
			llmArgs := args.LLMProcessorArguments{
				DatasetId: "ds1",
				Prompt:    "summarize: ${markdown}",
			}
			jsonData, err := json.Marshal(llmArgs)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &llmArgs)
			Expect(err).ToNot(HaveOccurred())
			Expect(llmArgs.Temperature).To(Equal(0.1))
			Expect(llmArgs.MaxTokens).To(Equal(uint(300)))
			Expect(llmArgs.Items).To(Equal(uint(1)))
		})

		It("should override default values", func() {
			llmArgs := args.LLMProcessorArguments{
				DatasetId:   "ds1",
				Prompt:      "summarize: ${markdown}",
				MaxTokens:   123,
				Temperature: 0.7,
				Items:       3,
			}
			jsonData, err := json.Marshal(llmArgs)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &llmArgs)
			Expect(err).ToNot(HaveOccurred())
			Expect(llmArgs.Temperature).To(Equal(0.7))
			Expect(llmArgs.MaxTokens).To(Equal(uint(123)))
			Expect(llmArgs.Items).To(Equal(uint(3)))
		})

		It("should fail unmarshal when dataset_id is missing", func() {
			var llmArgs args.LLMProcessorArguments
			jsonData := []byte(`{"type":"datasetprocessor","prompt":"p"}`)
			err := json.Unmarshal(jsonData, &llmArgs)
			Expect(errors.Is(err, args.ErrLLMDatasetIdRequired)).To(BeTrue())
		})

		It("should fail unmarshal when prompt is missing", func() {
			var llmArgs args.LLMProcessorArguments
			jsonData := []byte(`{"type":"datasetprocessor","dataset_id":"ds1"}`)
			err := json.Unmarshal(jsonData, &llmArgs)
			Expect(errors.Is(err, args.ErrLLMPromptRequired)).To(BeTrue())
		})
	})

	Describe("Validation", func() {
		It("should succeed with valid arguments", func() {
			llmArgs := &args.LLMProcessorArguments{
				DatasetId:   "ds1",
				Prompt:      "p",
				MaxTokens:   10,
				Temperature: 0.2,
				Items:       1,
			}
			err := llmArgs.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when dataset_id is missing", func() {
			llmArgs := &args.LLMProcessorArguments{
				Prompt:      "p",
				MaxTokens:   10,
				Temperature: 0.2,
			}
			err := llmArgs.Validate()
			Expect(errors.Is(err, args.ErrLLMDatasetIdRequired)).To(BeTrue())
		})

		It("should fail when prompt is missing", func() {
			llmArgs := &args.LLMProcessorArguments{
				DatasetId:   "ds1",
				MaxTokens:   10,
				Temperature: 0.2,
			}
			err := llmArgs.Validate()
			Expect(errors.Is(err, args.ErrLLMPromptRequired)).To(BeTrue())
		})
	})

	Describe("ToLLMProcessorRequest", func() {
		It("should map request fields to actor request fields", func() {
			llmArgs := args.LLMProcessorArguments{
				DatasetId:   "ds1",
				Prompt:      "p",
				MaxTokens:   42,
				Temperature: 0.7,
			}
			req, err := llmArgs.ToLLMProcessorRequest(args.LLMDefaultGeminiModel, "api-key")
			Expect(err).ToNot(HaveOccurred())
			Expect(req.InputDatasetId).To(Equal("ds1"))
			Expect(req.Prompt).To(Equal("p"))
			Expect(req.MaxTokens).To(Equal(uint(42)))
			Expect(req.Temperature).To(Equal("0.7"))
			Expect(req.MultipleColumns).To(BeFalse())
			Expect(req.Model).To(Equal(args.LLMDefaultGeminiModel))
			Expect(req.LLMProviderApiKey).To(Equal("api-key"))
		})
	})
})
