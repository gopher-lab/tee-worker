package profile_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-types/args"
	"github.com/masa-finance/tee-types/args/linkedin/profile"
	"github.com/masa-finance/tee-types/types"
	"github.com/masa-finance/tee-types/types/linkedin/experiences"
	"github.com/masa-finance/tee-types/types/linkedin/functions"
	"github.com/masa-finance/tee-types/types/linkedin/industries"
	profiletypes "github.com/masa-finance/tee-types/types/linkedin/profile"
	"github.com/masa-finance/tee-types/types/linkedin/seniorities"
)

var _ = Describe("LinkedIn Profile Arguments", func() {
	Describe("Marshalling and unmarshalling", func() {
		It("should set default values", func() {
			args := args.LinkedInProfileArguments{
				QueryType: types.CapSearchByProfile,
				Query:     "software engineer",
			}
			jsonData, err := json.Marshal(args)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.MaxItems).To(Equal(uint(10)))
			Expect(args.ScraperMode).To(Equal(profiletypes.ScraperModeShort))
		})

		It("should override default values", func() {
			args := args.LinkedInProfileArguments{
				QueryType:   types.CapSearchByProfile,
				Query:       "software engineer",
				MaxItems:    50,
				ScraperMode: profiletypes.ScraperModeFull,
			}
			jsonData, err := json.Marshal(args)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.MaxItems).To(Equal(uint(50)))
			Expect(args.ScraperMode).To(Equal(profiletypes.ScraperModeFull))
		})
	})

	Describe("Validation", func() {
		It("should succeed with valid arguments", func() {
			args := args.LinkedInProfileArguments{
				QueryType:         types.CapSearchByProfile,
				Query:             "software engineer",
				ScraperMode:       profiletypes.ScraperModeShort,
				MaxItems:          10,
				YearsOfExperience: []experiences.Id{experiences.ThreeToFiveYears},
				SeniorityLevels:   []seniorities.Id{seniorities.Senior},
				Functions:         []functions.Id{functions.Engineering},
				Industries:        []industries.Id{industries.SoftwareDevelopment},
			}
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with max items too large", func() {
			args := args.LinkedInProfileArguments{
				QueryType:   types.CapSearchByProfile,
				Query:       "software engineer",
				ScraperMode: profiletypes.ScraperModeShort,
				MaxItems:    1500,
			}
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrMaxItemsTooLarge)).To(BeTrue())
		})

		It("should fail with invalid scraper mode", func() {
			args := args.LinkedInProfileArguments{
				QueryType:   types.CapSearchByProfile,
				Query:       "software engineer",
				ScraperMode: "InvalidMode",
				MaxItems:    10,
			}
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrScraperModeNotSupported)).To(BeTrue())
		})

		It("should fail with invalid years of experience", func() {
			args := args.LinkedInProfileArguments{
				QueryType:         types.CapSearchByProfile,
				Query:             "software engineer",
				ScraperMode:       profiletypes.ScraperModeShort,
				MaxItems:          10,
				YearsOfExperience: []experiences.Id{"invalid"},
			}
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrExperienceNotSupported)).To(BeTrue())

		})

		It("should fail with invalid years at current company", func() {
			args := args.LinkedInProfileArguments{
				QueryType:             types.CapSearchByProfile,
				Query:                 "software engineer",
				ScraperMode:           profiletypes.ScraperModeShort,
				MaxItems:              10,
				YearsAtCurrentCompany: []experiences.Id{"invalid"},
			}
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrExperienceNotSupported)).To(BeTrue())

		})

		It("should fail with invalid seniority level", func() {
			args := args.LinkedInProfileArguments{
				QueryType:       types.CapSearchByProfile,
				Query:           "software engineer",
				ScraperMode:     profiletypes.ScraperModeShort,
				MaxItems:        10,
				SeniorityLevels: []seniorities.Id{"invalid"},
			}
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrSeniorityNotSupported)).To(BeTrue())
		})

		It("should fail with invalid function", func() {
			args := args.LinkedInProfileArguments{
				QueryType:   types.CapSearchByProfile,
				Query:       "software engineer",
				ScraperMode: profiletypes.ScraperModeShort,
				MaxItems:    10,
				Functions:   []functions.Id{"invalid"},
			}
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrFunctionNotSupported)).To(BeTrue())

		})

		It("should fail with invalid industry", func() {
			args := args.LinkedInProfileArguments{
				QueryType:   types.CapSearchByProfile,
				Query:       "software engineer",
				ScraperMode: profiletypes.ScraperModeShort,
				MaxItems:    10,
				Industries:  []industries.Id{"invalid"},
			}
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrIndustryNotSupported)).To(BeTrue())

		})

		It("should handle multiple validation errors", func() {
			args := args.LinkedInProfileArguments{
				QueryType:         types.CapSearchByProfile,
				Query:             "software engineer",
				ScraperMode:       "InvalidMode",
				MaxItems:          1500,
				YearsOfExperience: []experiences.Id{"invalid"},
				SeniorityLevels:   []seniorities.Id{"invalid"},
			}
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			// Should contain multiple error messages
			Expect(errors.Is(err, profile.ErrMaxItemsTooLarge)).To(BeTrue())
			Expect(errors.Is(err, profile.ErrScraperModeNotSupported)).To(BeTrue())
			Expect(errors.Is(err, profile.ErrExperienceNotSupported)).To(BeTrue())
			Expect(errors.Is(err, profile.ErrSeniorityNotSupported)).To(BeTrue())
		})
	})

	Describe("GetCapability", func() {
		It("should return the query type", func() {
			args := args.LinkedInProfileArguments{
				QueryType: types.CapSearchByProfile,
			}
			Expect(args.GetCapability()).To(Equal(types.CapSearchByProfile))
		})
	})

	Describe("ValidateForJobType", func() {
		It("should succeed with valid job type and capability", func() {
			args := args.LinkedInProfileArguments{
				QueryType:   types.CapSearchByProfile,
				Query:       "software engineer",
				ScraperMode: profiletypes.ScraperModeShort,
				MaxItems:    10,
			}
			err := args.ValidateForJobType(types.LinkedInJob)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with invalid job type", func() {
			args := args.LinkedInProfileArguments{
				QueryType:   types.CapSearchByQuery,
				Query:       "software engineer",
				ScraperMode: profiletypes.ScraperModeShort,
				MaxItems:    10,
			}
			err := args.ValidateForJobType(types.LinkedInJob)
			Expect(err).To(HaveOccurred())
		})

		It("should fail if base validation fails", func() {
			args := args.LinkedInProfileArguments{
				QueryType:   types.CapSearchByProfile,
				Query:       "software engineer",
				ScraperMode: "InvalidMode",
				MaxItems:    10,
			}
			err := args.ValidateForJobType(types.LinkedInJob)
			Expect(err).To(HaveOccurred())
		})
	})
})
