package args

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	teetypes "github.com/masa-finance/tee-types/types"
)

var (
	ErrWebURLRequired      = errors.New("url is required")
	ErrWebURLInvalid       = errors.New("invalid URL format")
	ErrWebURLSchemeMissing = errors.New("url must include a scheme (http:// or https://)")
	ErrWebMaxDepth         = errors.New("max depth must be non-negative")
	ErrWebMaxPages         = errors.New("max pages must be at least 1")
)

const (
	WebDefaultMaxPages             = 1
	WebDefaultMethod               = "GET"
	WebDefaultRespectRobotsTxtFile = false
	WebDefaultSaveMarkdown         = true
)

type WebArguments struct {
	QueryType teetypes.WebQueryType `json:"type"`
	URL       string                `json:"url"`
	MaxDepth  int                   `json:"max_depth"`
	MaxPages  int                   `json:"max_pages"`
}

// UnmarshalJSON implements custom JSON unmarshaling with validation
func (w *WebArguments) UnmarshalJSON(data []byte) error {
	// Prevent infinite recursion (you call json.Unmarshal which then calls `UnmarshalJSON`, which then calls `json.Unmarshal`...)
	type Alias WebArguments
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(w),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("failed to unmarshal web arguments: %w", err)
	}

	w.setDefaultValues()

	return w.Validate()
}

func (w *WebArguments) setDefaultValues() {
	if w.MaxPages == 0 {
		w.MaxPages = WebDefaultMaxPages
	}
}

// Validate validates the Web arguments
func (w *WebArguments) Validate() error {
	if w.URL == "" {
		return ErrWebURLRequired
	}

	// Validate URL format
	parsedURL, err := url.Parse(w.URL)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWebURLInvalid, err)
	}

	// Ensure URL has a scheme
	if parsedURL.Scheme == "" {
		return ErrWebURLSchemeMissing
	}

	if w.MaxDepth < 0 {
		return fmt.Errorf("%w: got %v", ErrWebMaxDepth, w.MaxDepth)
	}

	if w.MaxPages < 1 {
		return fmt.Errorf("%w: got %v", ErrWebMaxPages, w.MaxPages)
	}

	return nil
}

// ValidateForJobType validates Web arguments for a specific job type
func (w *WebArguments) ValidateForJobType(jobType teetypes.JobType) error {
	if err := w.Validate(); err != nil {
		return err
	}

	// Validate capability against job-specific capabilities
	return jobType.ValidateCapability(w.GetCapability())
}

// GetCapability returns the capability for web operations (always scraper)
func (w *WebArguments) GetCapability() teetypes.Capability {
	return teetypes.CapScraper
}

func (w WebArguments) ToWebScraperRequest() teetypes.WebScraperRequest {
	return teetypes.WebScraperRequest{
		StartUrls: []teetypes.WebStartURL{
			{URL: w.URL, Method: WebDefaultMethod},
		},
		MaxCrawlDepth:        w.MaxDepth,
		MaxCrawlPages:        w.MaxPages,
		RespectRobotsTxtFile: WebDefaultRespectRobotsTxtFile,
		SaveMarkdown:         WebDefaultSaveMarkdown,
	}
}
