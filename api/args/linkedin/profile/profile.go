package profile

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/masa-finance/tee-worker/api/args/base"
	"github.com/masa-finance/tee-worker/api/types"
	"github.com/masa-finance/tee-worker/api/types/linkedin/experiences"
	"github.com/masa-finance/tee-worker/api/types/linkedin/functions"
	"github.com/masa-finance/tee-worker/api/types/linkedin/industries"
	"github.com/masa-finance/tee-worker/api/types/linkedin/profile"
	"github.com/masa-finance/tee-worker/api/types/linkedin/seniorities"
)

var (
	ErrScraperModeNotSupported = errors.New("scraper mode not supported")
	ErrMaxItemsTooLarge        = errors.New("max items must be less than or equal to 100")
	ErrExperienceNotSupported  = errors.New("years of experience not supported")
	ErrSeniorityNotSupported   = errors.New("seniority level not supported")
	ErrFunctionNotSupported    = errors.New("function not supported")
	ErrIndustryNotSupported    = errors.New("industry not supported")
	ErrUnmarshalling           = errors.New("failed to unmarshal LinkedIn profile arguments")
)

const (
	DefaultMaxItems    = 10
	DefaultScraperMode = profile.ScraperModeShort
	MaxItems           = 1000 // 2500 on the actor, but we will run over 1MB memory limit on responses
)

// Verify interface implementation
var _ base.JobArgument = (*Arguments)(nil)

// Arguments defines args for LinkedIn profile operations
type Arguments struct {
	Type                  types.Capability    `json:"type"`
	ScraperMode           profile.ScraperMode `json:"profileScraperMode"`
	Query                 string              `json:"searchQuery"`
	MaxItems              uint                `json:"maxItems"`
	Locations             []string            `json:"locations,omitempty"`
	CurrentCompanies      []string            `json:"currentCompanies,omitempty"`
	PastCompanies         []string            `json:"pastCompanies,omitempty"`
	CurrentJobTitles      []string            `json:"currentJobTitles,omitempty"`
	PastJobTitles         []string            `json:"pastJobTitles,omitempty"`
	Schools               []string            `json:"schools,omitempty"`
	YearsOfExperience     []experiences.Id    `json:"yearsOfExperienceIds,omitempty"`
	YearsAtCurrentCompany []experiences.Id    `json:"yearsAtCurrentCompanyIds,omitempty"`
	SeniorityLevels       []seniorities.Id    `json:"seniorityLevelIds,omitempty"`
	Functions             []functions.Id      `json:"functionIds,omitempty"`
	Industries            []industries.Id     `json:"industryIds,omitempty"`
	FirstNames            []string            `json:"firstNames,omitempty"`
	LastNames             []string            `json:"lastNames,omitempty"`
	RecentlyChangedJobs   bool                `json:"recentlyChangedJobs,omitempty"`
	StartPage             uint                `json:"startPage,omitempty"`
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
	if a.MaxItems == 0 {
		a.MaxItems = DefaultMaxItems
	}
	if a.ScraperMode == "" {
		a.ScraperMode = DefaultScraperMode
	}
}

// TODO: use a validation library
func (a *Arguments) Validate() error {
	var errs []error

	if a.MaxItems > MaxItems {
		errs = append(errs, ErrMaxItemsTooLarge)
	}

	err := a.ValidateCapability(types.LinkedInJob)
	if err != nil {
		errs = append(errs, err)
	}

	if !profile.AllScraperModes.Contains(a.ScraperMode) {
		errs = append(errs, ErrScraperModeNotSupported)
	}

	for _, yoe := range a.YearsOfExperience {
		if !experiences.All.Contains(yoe) {
			errs = append(errs, fmt.Errorf("%w: %v", ErrExperienceNotSupported, yoe))
		}
	}
	for _, yac := range a.YearsAtCurrentCompany {
		if !experiences.All.Contains(yac) {
			errs = append(errs, fmt.Errorf("%w: %v", ErrExperienceNotSupported, yac))
		}
	}
	for _, sl := range a.SeniorityLevels {
		if !seniorities.All.Contains(sl) {
			errs = append(errs, fmt.Errorf("%w: %v", ErrSeniorityNotSupported, sl))
		}
	}
	for _, f := range a.Functions {
		if !functions.All.Contains(f) {
			errs = append(errs, fmt.Errorf("%w: %v", ErrFunctionNotSupported, f))
		}
	}
	for _, i := range a.Industries {
		if !industries.All.Contains(i) {
			errs = append(errs, fmt.Errorf("%w: %v", ErrIndustryNotSupported, i))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (a *Arguments) GetCapability() types.Capability {
	return a.Type
}

func (a *Arguments) ValidateCapability(jobType types.JobType) error {
	return jobType.ValidateCapability(&a.Type)
}

// NewArguments creates a new Arguments instance and applies default values immediately
func NewArguments() Arguments {
	args := Arguments{}
	args.SetDefaultValues()
	args.Validate()
	return args
}
