package types

type LLMProcessorRequest struct {
	InputDatasetId    string `json:"inputDatasetId"`
	LLMProviderApiKey string `json:"llmProviderApiKey"` // encrypted api key by miner
	Model             string `json:"model"`
	MultipleColumns   bool   `json:"multipleColumns"`
	Prompt            string `json:"prompt"`      // example: summarize the content of this webpage: ${markdown}
	Temperature       string `json:"temperature"` // the actor expects a string
	MaxTokens         uint   `json:"maxTokens"`
}

type LLMProcessorResult struct {
	LLMResponse string `json:"llmresponse"`
}
