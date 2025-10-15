package linkedin

import (
	"github.com/masa-finance/tee-worker/api/types/linkedin/experiences"
	"github.com/masa-finance/tee-worker/api/types/linkedin/functions"
	"github.com/masa-finance/tee-worker/api/types/linkedin/industries"
	"github.com/masa-finance/tee-worker/api/types/linkedin/profile"
	"github.com/masa-finance/tee-worker/api/types/linkedin/seniorities"
)

type LinkedInConfig struct {
	Experiences *experiences.ExperiencesConfig
	Seniorities *seniorities.SenioritiesConfig
	Functions   *functions.FunctionsConfig
	Industries  *industries.IndustriesConfig
	Profile     *profile.Profile
}

var LinkedIn = LinkedInConfig{
	Experiences: &experiences.Experiences,
	Seniorities: &seniorities.Seniorities,
	Functions:   &functions.Functions,
	Industries:  &industries.Industries,
	Profile:     &profile.Profile{},
}

type Profile = profile.Profile
