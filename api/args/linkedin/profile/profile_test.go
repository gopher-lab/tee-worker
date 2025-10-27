package profile_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/v2/api/args/linkedin/profile"
	"github.com/masa-finance/tee-worker/v2/api/types"
	"github.com/masa-finance/tee-worker/v2/api/types/linkedin/experiences"
	"github.com/masa-finance/tee-worker/v2/api/types/linkedin/functions"
	"github.com/masa-finance/tee-worker/v2/api/types/linkedin/industries"
	ptypes "github.com/masa-finance/tee-worker/v2/api/types/linkedin/profile"
	"github.com/masa-finance/tee-worker/v2/api/types/linkedin/seniorities"
)

var _ = Describe("LinkedIn Profile Arguments", func() {
	Describe("Marshalling and unmarshalling", func() {
		It("should set default values", func() {
			args := profile.NewArguments()
			args.Query = "software engineer"
			jsonData, err := json.Marshal(args)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.MaxItems).To(Equal(uint(10)))
			Expect(args.ScraperMode).To(Equal(ptypes.ScraperModeShort))
		})

		It("should override default values", func() {
			args := profile.NewArguments()
			args.Query = "software engineer"
			args.MaxItems = 50
			args.ScraperMode = ptypes.ScraperModeFull
			jsonData, err := json.Marshal(args)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.MaxItems).To(Equal(uint(50)))
			Expect(args.ScraperMode).To(Equal(ptypes.ScraperModeFull))
		})
	})

	Describe("Validation", func() {
		It("should succeed with valid arguments", func() {
			args := profile.NewArguments()
			args.Query = "software engineer"
			args.ScraperMode = ptypes.ScraperModeShort
			args.MaxItems = 10
			args.YearsOfExperience = []experiences.Id{experiences.ThreeToFiveYears}
			args.SeniorityLevels = []seniorities.Id{seniorities.Senior}
			args.Functions = []functions.Id{functions.Engineering}
			args.Industries = []industries.Id{industries.SoftwareDevelopment}
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with max items too large", func() {
			args := profile.NewArguments()
			args.Query = "software engineer"
			args.ScraperMode = ptypes.ScraperModeShort
			args.MaxItems = 1500
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrMaxItemsTooLarge)).To(BeTrue())
		})

		It("should fail with invalid scraper mode", func() {
			args := profile.NewArguments()
			args.Query = "software engineer"
			args.ScraperMode = "InvalidMode"
			args.MaxItems = 10
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrScraperModeNotSupported)).To(BeTrue())
		})

		It("should fail with invalid years of experience", func() {
			args := profile.NewArguments()
			args.Query = "software engineer"
			args.ScraperMode = ptypes.ScraperModeShort
			args.MaxItems = 10
			args.YearsOfExperience = []experiences.Id{"invalid"}
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrExperienceNotSupported)).To(BeTrue())

		})

		It("should fail with invalid years at current company", func() {
			args := profile.NewArguments()
			args.Query = "software engineer"
			args.ScraperMode = ptypes.ScraperModeShort
			args.MaxItems = 10
			args.YearsAtCurrentCompany = []experiences.Id{"invalid"}
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrExperienceNotSupported)).To(BeTrue())

		})

		It("should fail with invalid seniority level", func() {
			args := profile.NewArguments()
			args.Query = "software engineer"
			args.ScraperMode = ptypes.ScraperModeShort
			args.MaxItems = 10
			args.SeniorityLevels = []seniorities.Id{"invalid"}
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrSeniorityNotSupported)).To(BeTrue())
		})

		It("should fail with invalid function", func() {
			args := profile.NewArguments()
			args.Query = "software engineer"
			args.ScraperMode = ptypes.ScraperModeShort
			args.MaxItems = 10
			args.Functions = []functions.Id{"invalid"}
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrFunctionNotSupported)).To(BeTrue())

		})

		It("should fail with invalid industry", func() {
			args := profile.NewArguments()
			args.Query = "software engineer"
			args.ScraperMode = ptypes.ScraperModeShort
			args.MaxItems = 10
			args.Industries = []industries.Id{"invalid"}
			err := args.Validate()
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, profile.ErrIndustryNotSupported)).To(BeTrue())

		})

		It("should handle multiple validation errors", func() {
			args := profile.NewArguments()
			args.Query = "software engineer"
			args.ScraperMode = "InvalidMode"
			args.MaxItems = 1500
			args.YearsOfExperience = []experiences.Id{"invalid"}
			args.SeniorityLevels = []seniorities.Id{"invalid"}
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
			args := profile.NewArguments()
			Expect(args.GetCapability()).To(Equal(types.CapSearchByProfile))
		})
	})
})
