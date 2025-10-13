package args

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	teetypes "github.com/masa-finance/tee-types/types"
)

// Period constants for TikTok trending search
const (
	periodWeek  string = "7"
	periodMonth string = "30"
)

const (
	sortTrending string = "vv"
	sortLike     string = "like"
	sortComment  string = "comment"
	sortRepost   string = "repost"
)

// TikTokTranscriptionArguments defines args for TikTok transcriptions
type TikTokTranscriptionArguments struct {
	VideoURL string `json:"video_url"`
	Language string `json:"language,omitempty"` // e.g., "eng-US"
}

// UnmarshalJSON implements custom JSON unmarshaling with validation
func (t *TikTokTranscriptionArguments) UnmarshalJSON(data []byte) error {
	// Prevent infinite recursion (you call json.Unmarshal which then calls `UnmarshalJSON`, which then calls `json.Unmarshal`...)
	type Alias TikTokTranscriptionArguments
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("failed to unmarshal TikTok arguments: %w", err)
	}

	return t.Validate()
}

// Validate validates the TikTok arguments
func (t *TikTokTranscriptionArguments) Validate() error {
	if t.VideoURL == "" {
		return errors.New("video_url is required")
	}

	// Validate URL format
	parsedURL, err := url.Parse(t.VideoURL)
	if err != nil {
		return fmt.Errorf("invalid video_url format: %w", err)
	}

	// Basic TikTok URL validation
	if !t.IsTikTokURL(parsedURL) {
		return errors.New("URL must be a valid TikTok video URL")
	}

	// Validate language format if provided
	if t.Language != "" {
		if err := t.validateLanguageCode(); err != nil {
			return err
		}
	}

	return nil
}

// GetCapability returns the capability for TikTok operations (always transcription)
func (t *TikTokTranscriptionArguments) GetCapability() teetypes.Capability {
	return teetypes.CapTranscription
}

// IsTikTokURL validates if the URL is a TikTok URL
func (t *TikTokTranscriptionArguments) IsTikTokURL(parsedURL *url.URL) bool {
	host := strings.ToLower(parsedURL.Host)
	return host == "tiktok.com" || strings.HasSuffix(host, ".tiktok.com")
}

// HasLanguagePreference returns true if a language preference is specified
func (t *TikTokTranscriptionArguments) HasLanguagePreference() bool {
	return t.Language != ""
}

// GetVideoURL returns the source video URL
func (t *TikTokTranscriptionArguments) GetVideoURL() string {
	return t.VideoURL
}

// GetLanguageCode returns the language code, defaulting to "en-us" if not specified
func (t *TikTokTranscriptionArguments) GetLanguageCode() string {
	if t.Language == "" {
		return "eng-US"
	}
	return t.Language
}

// ValidateForJobType validates TikTok arguments for a specific job type
func (t *TikTokTranscriptionArguments) ValidateForJobType(jobType teetypes.JobType) error {
	if err := t.Validate(); err != nil {
		return err
	}

	// Validate capability against job-specific capabilities
	return jobType.ValidateCapability(t.GetCapability())
}

// validateLanguageCode validates the language code format
func (t *TikTokTranscriptionArguments) validateLanguageCode() error {
	// Basic validation for language codes like "en-us", "eng-us", "es-es", etc.
	parts := strings.Split(t.Language, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid language format '%s', expected format: 'lang-region' (e.g., 'en-us' or 'eng-us')", t.Language)
	}

	// Language code can be 2 or 3 letters, region must be 2 letters
	if (len(parts[0]) != 2 && len(parts[0]) != 3) || len(parts[1]) != 2 {
		return fmt.Errorf("invalid language format '%s', expected 2-3 letter language code and 2-letter region code", t.Language)
	}

	return nil
}

// TikTokSearchByQueryArguments defines args for epctex/tiktok-search-scraper
type TikTokSearchByQueryArguments struct {
	QueryType string   `json:"type"`
	Search    []string `json:"search,omitempty"`
	StartUrls []string `json:"start_urls,omitempty"`
	MaxItems  uint     `json:"max_items,omitempty"`
	EndPage   uint     `json:"end_page,omitempty"`
}

func (t *TikTokSearchByQueryArguments) UnmarshalJSON(data []byte) error {
	// Prevent infinite recursion (you call json.Unmarshal which then calls `UnmarshalJSON`, which then calls `json.Unmarshal`...)
	type Alias TikTokSearchByQueryArguments
	aux := &struct{ *Alias }{Alias: (*Alias)(t)}
	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("failed to unmarshal TikTok searchbyquery arguments: %w", err)
	}
	t.QueryType = strings.ToLower(t.QueryType)
	return t.Validate()
}

func (t *TikTokSearchByQueryArguments) Validate() error {
	if len(t.Search) == 0 && len(t.StartUrls) == 0 {
		return errors.New("either 'search' or 'start_urls' is required for searchbyquery")
	}
	return nil
}

func (t *TikTokSearchByQueryArguments) ValidateForJobType(jobType teetypes.JobType) error {
	if err := jobType.ValidateCapability(teetypes.CapSearchByQuery); err != nil {
		return err
	}
	return t.Validate()
}

func (t *TikTokSearchByQueryArguments) GetCapability() teetypes.Capability {
	return teetypes.CapSearchByQuery
}

// TikTokSearchByTrendingArguments defines args for lexis-solutions/tiktok-trending-videos-scraper
type TikTokSearchByTrendingArguments struct {
	QueryType   string `json:"type"`
	CountryCode string `json:"country_code,omitempty"`
	SortBy      string `json:"sort_by,omitempty"`
	MaxItems    int    `json:"max_items,omitempty"`
	Period      string `json:"period,omitempty"`
}

func (t *TikTokSearchByTrendingArguments) UnmarshalJSON(data []byte) error {
	// Prevent infinite recursion (you call json.Unmarshal which then calls `UnmarshalJSON`, which then calls `json.Unmarshal`...)
	type Alias TikTokSearchByTrendingArguments
	aux := &struct{ *Alias }{Alias: (*Alias)(t)}
	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("failed to unmarshal TikTok searchbytrending arguments: %w", err)
	}
	t.QueryType = strings.ToLower(t.QueryType)
	if t.CountryCode == "" {
		t.CountryCode = "US"
	}
	if t.SortBy == "" {
		t.SortBy = sortTrending
	}
	if t.Period == "" {
		t.Period = periodWeek
	}
	return t.Validate()
}

func (t *TikTokSearchByTrendingArguments) Validate() error {
	allowedSorts := map[string]struct{}{
		sortTrending: {}, sortLike: {}, sortComment: {}, sortRepost: {},
	}

	allowedPeriods := map[string]struct{}{
		periodWeek:  {},
		periodMonth: {},
	}

	allowedCountries := map[string]struct{}{
		"AU": {}, "BR": {}, "CA": {}, "EG": {}, "FR": {}, "DE": {}, "ID": {}, "IL": {}, "IT": {}, "JP": {},
		"MY": {}, "PH": {}, "RU": {}, "SA": {}, "SG": {}, "KR": {}, "ES": {}, "TW": {}, "TH": {}, "TR": {},
		"AE": {}, "GB": {}, "US": {}, "VN": {},
	}

	if _, ok := allowedCountries[strings.ToUpper(t.CountryCode)]; !ok {
		return fmt.Errorf("invalid country_code '%s'", t.CountryCode)
	}
	if _, ok := allowedSorts[strings.ToLower(t.SortBy)]; !ok {
		return fmt.Errorf("invalid sort_by '%s'", t.SortBy)
	}
	if _, ok := allowedPeriods[t.Period]; !ok {
		// Extract keys for error message
		var validKeys []string
		for key := range allowedPeriods {
			validKeys = append(validKeys, key)
		}
		return fmt.Errorf("invalid period '%s' (allowed: %s)", t.Period, strings.Join(validKeys, ", "))
	}
	if t.MaxItems < 0 {
		return fmt.Errorf("max_items must be non-negative, got: %d", t.MaxItems)
	}
	return nil
}

func (t *TikTokSearchByTrendingArguments) ValidateForJobType(jobType teetypes.JobType) error {
	if err := jobType.ValidateCapability(teetypes.CapSearchByTrending); err != nil {
		return err
	}
	return t.Validate()
}

func (t *TikTokSearchByTrendingArguments) GetCapability() teetypes.Capability {
	return teetypes.CapSearchByTrending
}
