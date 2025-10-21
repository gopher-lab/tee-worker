package trending

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/masa-finance/tee-worker/api/args/base"
	"github.com/masa-finance/tee-worker/api/types"
)

var (
	ErrTrendingCountryCodeRequired = errors.New("country_code is required")
	ErrTrendingSortByRequired      = errors.New("sort_by is required")
	ErrTrendingPeriodRequired      = errors.New("period is required")
	ErrTrendingMaxItemsNegative    = errors.New("max_items must be non-negative")
	ErrUnmarshalling               = errors.New("failed to unmarshal TikTok searchbytrending arguments")
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

// Verify interface implementation
var _ base.JobArgument = (*Arguments)(nil)

// Arguments defines args for lexis-solutions/tiktok-trending-videos-scraper
type Arguments struct {
	base.Arguments
	CountryCode string `json:"country_code,omitempty"`
	SortBy      string `json:"sort_by,omitempty"`
	MaxItems    int    `json:"max_items,omitempty"`
	Period      string `json:"period,omitempty"`
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

func (a *Arguments) SetDefaultValues() {
	if a.CountryCode == "" {
		a.CountryCode = "US"
	}
	if a.SortBy == "" {
		a.SortBy = sortTrending
	}
	if a.Period == "" {
		a.Period = periodWeek
	}
}

// TODO: use a validation library
func (t *Arguments) Validate() error {
	if err := types.TiktokJob.ValidateCapability(&t.Type); err != nil {
		return err
	}
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
		return fmt.Errorf("%w: '%s'", ErrTrendingCountryCodeRequired, t.CountryCode)
	}
	if _, ok := allowedSorts[strings.ToLower(t.SortBy)]; !ok {
		return fmt.Errorf("%w: '%s'", ErrTrendingSortByRequired, t.SortBy)
	}
	if _, ok := allowedPeriods[t.Period]; !ok {
		// Extract keys for error message
		var validKeys []string
		for key := range allowedPeriods {
			validKeys = append(validKeys, key)
		}
		return fmt.Errorf("%w: '%s' (allowed: %s)", ErrTrendingPeriodRequired, t.Period, strings.Join(validKeys, ", "))
	}
	if t.MaxItems < 0 {
		return fmt.Errorf("%w, got: %d", ErrTrendingMaxItemsNegative, t.MaxItems)
	}
	return nil
}

func NewArguments() Arguments {
	args := Arguments{}
	args.Type = types.CapSearchByTrending
	args.SetDefaultValues()
	return args
}
