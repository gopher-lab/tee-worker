package linkedin_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-types/types"
	"github.com/masa-finance/tee-types/types/linkedin/experiences"
	"github.com/masa-finance/tee-types/types/linkedin/functions"
	"github.com/masa-finance/tee-types/types/linkedin/industries"
	"github.com/masa-finance/tee-types/types/linkedin/seniorities"
)

var _ = Describe("LinkedIn Types", func() {
	Describe("LinkedIn Package", func() {
		It("should have all required fields", func() {
			linkedin := types.LinkedIn

			Expect(linkedin.Seniorities).ToNot(BeNil())
			Expect(linkedin.Experiences).ToNot(BeNil())
			Expect(linkedin.Functions).ToNot(BeNil())
			Expect(linkedin.Industries).ToNot(BeNil())
		})
	})

	Describe("Seniorities", func() {
		It("should have all seniority levels", func() {
			s := types.LinkedIn.Seniorities

			Expect(s.InTraining).To(Equal(seniorities.InTraining))
			Expect(s.EntryLevel).To(Equal(seniorities.EntryLevel))
			Expect(s.Senior).To(Equal(seniorities.Senior))
			Expect(s.Strategic).To(Equal(seniorities.Strategic))
			Expect(s.EntryLevelManager).To(Equal(seniorities.EntryLevelManager))
			Expect(s.ExperiencedManager).To(Equal(seniorities.ExperiencedManager))
			Expect(s.Director).To(Equal(seniorities.Director))
			Expect(s.VicePresident).To(Equal(seniorities.VicePresident))
			Expect(s.CXO).To(Equal(seniorities.CXO))
			Expect(s.Partner).To(Equal(seniorities.Partner))
		})

		It("should have All set containing all seniorities", func() {
			all := types.LinkedIn.Seniorities.All

			Expect(all.Contains(seniorities.InTraining)).To(BeTrue())
			Expect(all.Contains(seniorities.EntryLevel)).To(BeTrue())
			Expect(all.Contains(seniorities.Senior)).To(BeTrue())
			Expect(all.Contains(seniorities.Strategic)).To(BeTrue())
			Expect(all.Contains(seniorities.EntryLevelManager)).To(BeTrue())
			Expect(all.Contains(seniorities.ExperiencedManager)).To(BeTrue())
			Expect(all.Contains(seniorities.Director)).To(BeTrue())
			Expect(all.Contains(seniorities.VicePresident)).To(BeTrue())
			Expect(all.Contains(seniorities.CXO)).To(BeTrue())
			Expect(all.Contains(seniorities.Partner)).To(BeTrue())

			Expect(all.Length()).To(Equal(10))
		})
	})

	Describe("Experiences", func() {
		It("should have all experience levels", func() {
			e := types.LinkedIn.Experiences

			Expect(e.LessThanAYear).To(Equal(experiences.LessThanAYear))
			Expect(e.OneToTwoYears).To(Equal(experiences.OneToTwoYears))
			Expect(e.ThreeToFiveYears).To(Equal(experiences.ThreeToFiveYears))
			Expect(e.SixToTenYears).To(Equal(experiences.SixToTenYears))
			Expect(e.MoreThanTenYears).To(Equal(experiences.MoreThanTenYears))
		})

		It("should have All set containing all experiences", func() {
			all := types.LinkedIn.Experiences.All

			Expect(all.Contains(experiences.LessThanAYear)).To(BeTrue())
			Expect(all.Contains(experiences.OneToTwoYears)).To(BeTrue())
			Expect(all.Contains(experiences.ThreeToFiveYears)).To(BeTrue())
			Expect(all.Contains(experiences.SixToTenYears)).To(BeTrue())
			Expect(all.Contains(experiences.MoreThanTenYears)).To(BeTrue())

			Expect(all.Length()).To(Equal(5))
		})
	})

	Describe("Functions", func() {
		It("should have all function types", func() {
			f := types.LinkedIn.Functions

			Expect(f.Accounting).To(Equal(functions.Accounting))
			Expect(f.Engineering).To(Equal(functions.Engineering))
			Expect(f.Marketing).To(Equal(functions.Marketing))
			Expect(f.Sales).To(Equal(functions.Sales))
			Expect(f.HumanResources).To(Equal(functions.HumanResources))
		})

		It("should have All set containing all functions", func() {
			all := types.LinkedIn.Functions.All

			Expect(all.Contains(functions.Accounting)).To(BeTrue())
			Expect(all.Contains(functions.Engineering)).To(BeTrue())
			Expect(all.Contains(functions.Marketing)).To(BeTrue())
			Expect(all.Contains(functions.Sales)).To(BeTrue())
			Expect(all.Contains(functions.HumanResources)).To(BeTrue())
			Expect(all.Contains(functions.InformationTechnology)).To(BeTrue())
			Expect(all.Contains(functions.Finance)).To(BeTrue())

			Expect(all.Length()).To(Equal(25))
		})
	})

	Describe("Industries", func() {
		It("should have all industry types", func() {
			i := types.LinkedIn.Industries

			Expect(i.SoftwareDevelopment).To(Equal(industries.SoftwareDevelopment))
			Expect(i.FinancialServices).To(Equal(industries.FinancialServices))
			Expect(i.Manufacturing).To(Equal(industries.Manufacturing))
			Expect(i.Retail).To(Equal(industries.Retail))
			Expect(i.Education).To(Equal(industries.Education))
		})

		It("should have All set containing all industries", func() {
			all := types.LinkedIn.Industries.All

			Expect(all.Contains(industries.SoftwareDevelopment)).To(BeTrue())
			Expect(all.Contains(industries.FinancialServices)).To(BeTrue())
			Expect(all.Contains(industries.Manufacturing)).To(BeTrue())
			Expect(all.Contains(industries.Retail)).To(BeTrue())
			Expect(all.Contains(industries.Education)).To(BeTrue())
			Expect(all.Contains(industries.Hospitals)).To(BeTrue())
			Expect(all.Contains(industries.ProfessionalServices)).To(BeTrue())

			Expect(all.Length()).To(BeNumerically(">=", 100)) // Should have many industries
		})
	})
})
