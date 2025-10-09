package linkedin

import (
	"github.com/masa-finance/tee-types/types/linkedin/experiences"
	"github.com/masa-finance/tee-types/types/linkedin/functions"
	"github.com/masa-finance/tee-types/types/linkedin/industries"
	"github.com/masa-finance/tee-types/types/linkedin/profile"
	"github.com/masa-finance/tee-types/types/linkedin/seniorities"
)

type LinkedInConfig struct {
	Experiences *experiences.ExperiencesConfig
	Seniorities *seniorities.SenioritiesConfig
	Functions   *functions.FunctionsConfig
	Industries  *industries.IndustriesConfig
}

var LinkedIn = LinkedInConfig{
	Experiences: &experiences.Experiences,
	Seniorities: &seniorities.Seniorities,
	Functions:   &functions.Functions,
	Industries:  &industries.Industries,
}

type Profile = *profile.Profile
