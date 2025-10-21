package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/masa-finance/tee-worker/api/args/base"
	"github.com/masa-finance/tee-worker/api/types"
)

var (
	ErrURLRequired      = errors.New("url is required")
	ErrURLInvalid       = errors.New("invalid URL format")
	ErrURLSchemeMissing = errors.New("url must include a scheme (http:// or https://)")
	ErrMaxDepth         = errors.New("max depth must be non-negative")
	ErrMaxPages         = errors.New("max pages must be at least 1")
	ErrUnmarshalling    = errors.New("failed to unmarshal web page arguments")
)

const (
	DefaultMaxPages             = 1
	DefaultMethod               = "GET"
	DefaultRespectRobotsTxtFile = false
	DefaultSaveMarkdown         = true
)

// Verify interface implementation
var _ base.JobArgument = (*Arguments)(nil)

type Arguments struct {
	Type     types.Capability `json:"type"`
	URL      string           `json:"url"`
	MaxDepth int              `json:"max_depth"`
	MaxPages int              `json:"max_pages"`
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

func (w *Arguments) SetDefaultValues() {
	if w.MaxPages == 0 {
		w.MaxPages = DefaultMaxPages
	}
}

// Validate validates the  arguments
// TODO: use a validation library
func (w *Arguments) Validate() error {
	err := w.ValidateCapability(types.WebJob)
	if err != nil {
		return err
	}

	if w.URL == "" {
		return ErrURLRequired
	}

	// Validate URL format
	parsedURL, err := url.Parse(w.URL)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrURLInvalid, err)
	}

	// Ensure URL has a scheme
	if parsedURL.Scheme == "" {
		return ErrURLSchemeMissing
	}

	if w.MaxDepth < 0 {
		return fmt.Errorf("%w: got %v", ErrMaxDepth, w.MaxDepth)
	}

	if w.MaxPages < 1 {
		return fmt.Errorf("%w: got %v", ErrMaxPages, w.MaxPages)
	}

	return nil
}

func (w *Arguments) GetCapability() types.Capability {
	return w.Type
}

func (w *Arguments) ValidateCapability(jobType types.JobType) error {
	return jobType.ValidateCapability(&w.Type)
}

func (w Arguments) ToScraperRequest() types.WebScraperRequest {
	return types.WebScraperRequest{
		StartUrls: []types.WebStartURL{
			{URL: w.URL, Method: DefaultMethod},
		},
		MaxCrawlDepth:        w.MaxDepth,
		MaxCrawlPages:        w.MaxPages,
		RespectRobotsTxtFile: DefaultRespectRobotsTxtFile,
		SaveMarkdown:         DefaultSaveMarkdown,
	}
}

func NewArguments() Arguments {
	args := Arguments{}
	args.SetDefaultValues()
	args.Validate()
	return args
}
