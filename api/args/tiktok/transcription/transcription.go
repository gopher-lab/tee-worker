package transcription

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/masa-finance/tee-worker/api/args/base"
	"github.com/masa-finance/tee-worker/api/types"
)

var (
	ErrVideoURLRequired    = errors.New("video_url is required")
	ErrInvalidVideoURL     = errors.New("invalid video_url format")
	ErrInvalidTikTokURL    = errors.New("url must be a valid TikTok video URL")
	ErrInvalidLanguageCode = errors.New("invalid language code")
	ErrUnmarshalling       = errors.New("failed to unmarshal TikTok transcription arguments")
)

const (
	DefaultLanguage = "eng-US"
)

// Verify interface implementation
var _ base.JobArgument = (*Arguments)(nil)

// Arguments defines args for TikTok transcriptions
type Arguments struct {
	Type     types.Capability `json:"type"`
	VideoURL string           `json:"video_url"`
	Language string           `json:"language,omitempty"`
}

func (a *Arguments) UnmarshalJSON(data []byte) error {
	type Alias Arguments
	aux := &struct{ *Alias }{Alias: (*Alias)(a)}
	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("%w: %w", ErrUnmarshalling, err)
	}
	a.SetDefaultValues()
	return a.Validate()
}

func (a *Arguments) SetDefaultValues() {
	if a.Language == "" {
		a.Language = DefaultLanguage
	}
}

// Validate validates the TikTok arguments
func (t *Arguments) Validate() error {
	err := t.ValidateCapability(types.TiktokJob)
	if err != nil {
		return err
	}
	if t.VideoURL == "" {
		return ErrVideoURLRequired
	}

	// Validate URL format
	parsedURL, err := url.Parse(t.VideoURL)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidVideoURL, err)
	}

	// Basic TikTok URL validation
	if !t.IsTikTokURL(parsedURL) {
		return ErrInvalidTikTokURL
	}

	// Validate language format if provided
	if t.Language != "" {
		if err := t.validateLanguageCode(); err != nil {
			return err
		}
	}

	return nil
}

func (t *Arguments) GetCapability() types.Capability {
	return t.Type
}

func (t *Arguments) ValidateCapability(jobType types.JobType) error {
	return jobType.ValidateCapability(&t.Type)
}

// IsTikTokURL validates if the URL is a TikTok URL
func (t *Arguments) IsTikTokURL(parsedURL *url.URL) bool {
	host := strings.ToLower(parsedURL.Host)
	return host == "tiktok.com" || strings.HasSuffix(host, ".tiktok.com")
}

// HasLanguagePreference returns true if a language preference is specified
func (t *Arguments) HasLanguagePreference() bool {
	return t.Language != ""
}

// GetVideoURL returns the source video URL
func (t *Arguments) GetVideoURL() string {
	return t.VideoURL
}

// GetLanguageCode returns the language code, defaulting to "en-us" if not specified
func (t *Arguments) GetLanguageCode() string {
	return t.Language
}

// validateLanguageCode validates the language code format
func (t *Arguments) validateLanguageCode() error {
	parts := strings.Split(t.Language, "-")
	if len(parts) != 2 || (len(parts[0]) != 2 && len(parts[0]) != 3) || len(parts[1]) != 2 {
		return fmt.Errorf("%w: %s", ErrInvalidLanguageCode, t.Language)
	}
	return nil
}

func NewArguments() Arguments {
	args := Arguments{
		Type: types.CapTranscription,
	}
	args.SetDefaultValues()
	return args
}
