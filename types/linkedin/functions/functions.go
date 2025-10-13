package functions

import "github.com/masa-finance/tee-types/pkg/util"

// id represents a LinkedIn function identifier
type Id string

// Function constants
const (
	Accounting                    Id = "1"
	Administrative                Id = "2"
	ArtsAndDesign                 Id = "3"
	BusinessDevelopment           Id = "4"
	CommunityAndSocialServices    Id = "5"
	Consulting                    Id = "6"
	Education                     Id = "7"
	Engineering                   Id = "8"
	Entrepreneurship              Id = "9"
	Finance                       Id = "10"
	HealthcareServices            Id = "11"
	HumanResources                Id = "12"
	InformationTechnology         Id = "13"
	Legal                         Id = "14"
	Marketing                     Id = "15"
	MediaAndCommunication         Id = "16"
	MilitaryAndProtectiveServices Id = "17"
	Operations                    Id = "18"
	ProductManagement             Id = "19"
	ProgramAndProjectManagement   Id = "20"
	Purchasing                    Id = "21"
	QualityAssurance              Id = "22"
	RealEstate                    Id = "23"
	Research                      Id = "24"
	Sales                         Id = "25"
)

var All = util.NewSet(
	Accounting,
	Administrative,
	ArtsAndDesign,
	BusinessDevelopment,
	CommunityAndSocialServices,
	Consulting,
	Education,
	Engineering,
	Entrepreneurship,
	Finance,
	HealthcareServices,
	HumanResources,
	InformationTechnology,
	Legal,
	Marketing,
	MediaAndCommunication,
	MilitaryAndProtectiveServices,
	Operations,
	ProductManagement,
	ProgramAndProjectManagement,
	Purchasing,
	QualityAssurance,
	RealEstate,
	Research,
	Sales,
)

type FunctionsConfig struct {
	All                           util.Set[Id]
	Accounting                    Id
	Administrative                Id
	ArtsAndDesign                 Id
	BusinessDevelopment           Id
	CommunityAndSocialServices    Id
	Consulting                    Id
	Education                     Id
	Engineering                   Id
	Entrepreneurship              Id
	Finance                       Id
	HealthcareServices            Id
	HumanResources                Id
	InformationTechnology         Id
	Legal                         Id
	Marketing                     Id
	MediaAndCommunication         Id
	MilitaryAndProtectiveServices Id
	Operations                    Id
	ProductManagement             Id
	ProgramAndProjectManagement   Id
	Purchasing                    Id
	QualityAssurance              Id
	RealEstate                    Id
	Research                      Id
	Sales                         Id
}

var Functions = FunctionsConfig{
	All:                           *All,
	Accounting:                    Accounting,
	Administrative:                Administrative,
	ArtsAndDesign:                 ArtsAndDesign,
	BusinessDevelopment:           BusinessDevelopment,
	CommunityAndSocialServices:    CommunityAndSocialServices,
	Consulting:                    Consulting,
	Education:                     Education,
	Engineering:                   Engineering,
	Entrepreneurship:              Entrepreneurship,
	Finance:                       Finance,
	HealthcareServices:            HealthcareServices,
	HumanResources:                HumanResources,
	InformationTechnology:         InformationTechnology,
	Legal:                         Legal,
	Marketing:                     Marketing,
	MediaAndCommunication:         MediaAndCommunication,
	MilitaryAndProtectiveServices: MilitaryAndProtectiveServices,
	Operations:                    Operations,
	ProductManagement:             ProductManagement,
	ProgramAndProjectManagement:   ProgramAndProjectManagement,
	Purchasing:                    Purchasing,
	QualityAssurance:              QualityAssurance,
	RealEstate:                    RealEstate,
	Research:                      Research,
	Sales:                         Sales,
}
