package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/masa-finance/tee-worker/api/args/base"
	"github.com/masa-finance/tee-worker/api/types"
	"github.com/masa-finance/tee-worker/pkg/util"
)

var (
	ErrDatasetIdRequired = errors.New("dataset id is required")
	ErrPromptRequired    = errors.New("prompt is required")
	ErrUnmarshalling     = errors.New("failed to unmarshal  arguments")
)

const (
	DefaultMaxTokens       uint    = 300
	DefaultTemperature     float64 = 0.1
	DefaultMultipleColumns bool    = false
	DefaultGeminiModel     string  = "gemini-1.5-flash-8b"
	DefaultClaudeModel     string  = "claude-3-5-haiku-latest"
	DefaultItems           uint    = 1
)

var SupportedModels = util.NewSet(DefaultGeminiModel, DefaultClaudeModel)

// Verify interface implementation
var _ base.JobArgument = (*Arguments)(nil)

type Arguments struct {
	Type        types.Capability `json:"type"`
	DatasetId   string           `json:"dataset_id"`
	Prompt      string           `json:"prompt"`
	MaxTokens   uint             `json:"max_tokens"`
	Temperature float64          `json:"temperature"`
	Items       uint             `json:"items"`
}

func (t *Arguments) UnmarshalJSON(data []byte) error {
	type Alias Arguments
	aux := &struct{ *Alias }{Alias: (*Alias)(t)}
	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("%w: %w", ErrUnmarshalling, err)
	}
	t.SetDefaultValues()
	return t.Validate()
}

func (l *Arguments) SetDefaultValues() {
	if l.Temperature == 0 {
		l.Temperature = DefaultTemperature
	}
	if l.MaxTokens == 0 {
		l.MaxTokens = DefaultMaxTokens
	}
	if l.Items == 0 {
		l.Items = DefaultItems
	}
}

func (l *Arguments) Validate() error {
	if l.DatasetId == "" {
		return ErrDatasetIdRequired
	}
	if l.Prompt == "" {
		return ErrPromptRequired
	}
	return nil
}

func (l *Arguments) GetCapability() types.Capability {
	return l.Type
}

func (l *Arguments) ValidateCapability(jobType types.JobType) error {
	return nil //  is not yet a standalone job type
}

// NewArguments creates a new Arguments instance and applies default values immediately
func NewArguments() Arguments {
	args := Arguments{}
	args.SetDefaultValues()
	args.Validate() // This will set the default capability via ValidateCapability
	return args
}

func (l Arguments) ToProcessorRequest(model string, key string) (types.LLMProcessorRequest, error) {
	if !SupportedModels.Contains(model) {
		return types.LLMProcessorRequest{}, fmt.Errorf("model %s is not supported", model)
	}
	if key == "" {
		return types.LLMProcessorRequest{}, fmt.Errorf("key is required")
	}

	return types.LLMProcessorRequest{
		InputDatasetId:    l.DatasetId,
		LLMProviderApiKey: key,
		Prompt:            l.Prompt,
		MaxTokens:         l.MaxTokens,
		Temperature:       strconv.FormatFloat(l.Temperature, 'f', -1, 64),
		MultipleColumns:   DefaultMultipleColumns, // overrides default in actor API
		Model:             model,                  // overrides default in actor API
	}, nil
}
