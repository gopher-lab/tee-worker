package industries

import "github.com/masa-finance/tee-worker/v2/pkg/util"

// Id represents a LinkedIn industry identifier
type Id string

// Industry constants
const (
	// Technology & Software
	SoftwareDevelopment              Id = "4"
	ComputerHardwareManufacturing    Id = "3"
	ComputerNetworkingProducts       Id = "5"
	ItServicesAndItConsulting        Id = "96"
	ComputerAndNetworkSecurity       Id = "118"
	Telecommunications               Id = "8"
	WirelessServices                 Id = "119"
	TechnologyInformationAndInternet Id = "6"
	DataInfrastructureAndAnalytics   Id = "2458"
	InformationServices              Id = "84"
	InternetPublishing               Id = "3132"
	SocialNetworkingPlatforms        Id = "3127"
	ComputerGames                    Id = "109"
	MobileGamingApps                 Id = "3131"
	BlockchainServices               Id = "3134"
	BusinessIntelligencePlatforms    Id = "3128"

	// Financial Services
	FinancialServices                        Id = "43"
	Banking                                  Id = "41"
	Insurance                                Id = "42"
	InvestmentBanking                        Id = "45"
	CapitalMarkets                           Id = "129"
	VentureCapitalAndPrivateEquityPrincipals Id = "106"
	SecuritiesAndCommodityExchanges          Id = "1713"
	FundsAndTrusts                           Id = "1742"

	// Healthcare & Medical
	Hospitals                     Id = "2081"
	MedicalPractices              Id = "13"
	MedicalEquipmentManufacturing Id = "17"
	PublicHealth                  Id = "2358"
	VeterinaryServices            Id = "16"
	BiotechnologyResearch         Id = "12"

	// Manufacturing
	Manufacturing                              Id = "25"
	ComputersAndElectronicsManufacturing       Id = "24"
	SemiconductorManufacturing                 Id = "7"
	MachineryManufacturing                     Id = "55"
	IndustrialMachineryManufacturing           Id = "135"
	FoodAndBeverageManufacturing               Id = "23"
	TextileManufacturing                       Id = "60"
	MotorVehicleManufacturing                  Id = "53"
	MotorVehiclePartsManufacturing             Id = "1042"
	AviationAndAerospaceComponentManufacturing Id = "52"
	DefenseAndSpaceManufacturing               Id = "1"
	PlasticsManufacturing                      Id = "117"
	RubberProductsManufacturing                Id = "763"
	PaperAndForestProductManufacturing         Id = "61"
	WoodProductManufacturing                   Id = "784"
	FurnitureAndHomeFurnishingsManufacturing   Id = "26"
	SportingGoodsManufacturing                 Id = "20"
	PrintingServices                           Id = "83"

	// Retail & Consumer Goods
	Retail                                           Id = "27"
	RetailGroceries                                  Id = "22"
	OnlineAndMailOrderRetail                         Id = "1445"
	RetailApparelAndFashion                          Id = "19"
	RetailAppliancesElectricalAndElectronicEquipment Id = "1319"
	RetailBooksAndPrintedNews                        Id = "1409"
	RetailBuildingMaterialsAndGardenEquipment        Id = "1324"
	RetailFurnitureAndHomeFurnishings                Id = "1309"
	RetailHealthAndPersonalCareProducts              Id = "1359"
	RetailLuxuryGoodsAndJewelry                      Id = "143"
	RetailMotorVehicles                              Id = "1292"
	RetailOfficeEquipment                            Id = "138"
	RetailOfficeSuppliesAndGifts                     Id = "1424"

	// Professional Services
	ProfessionalServices                     Id = "1810"
	Accounting                               Id = "47"
	LegalServices                            Id = "10"
	LawPractice                              Id = "9"
	BusinessConsultingAndServices            Id = "11"
	StrategicManagementServices              Id = "102"
	HumanResourcesServices                   Id = "137"
	MarketingServices                        Id = "1862"
	AdvertisingServices                      Id = "80"
	PublicRelationsAndCommunicationsServices Id = "98"
	MarketResearch                           Id = "97"
	ArchitectureAndPlanning                  Id = "50"
	DesignServices                           Id = "99"
	GraphicDesign                            Id = "140"
	InteriorDesign                           Id = "3126"
	EngineeringServices                      Id = "3242"
	EnvironmentalServices                    Id = "86"
	ResearchServices                         Id = "70"
	ThinkTanks                               Id = "130"
	Photography                              Id = "136"
	TranslationAndLocalization               Id = "108"
	WritingAndEditing                        Id = "103"

	// Education
	Education                       Id = "1999"
	HigherEducation                 Id = "68"
	ProfessionalTrainingAndCoaching Id = "105"
	SportsAndRecreationInstruction  Id = "2027"

	// Transportation & Logistics
	TransportationLogisticsSupplyChainAndStorage Id = "116"
	AirlinesAndAviation                          Id = "94"
	FreightAndPackageTransportation              Id = "87"
	MaritimeTransportation                       Id = "95"
	RailTransportation                           Id = "1481"
	TruckTransportation                          Id = "92"
	WarehousingAndStorage                        Id = "93"
	PostalServices                               Id = "1573"

	// Energy & Utilities
	Utilities                      Id = "59"
	ElectricPowerGeneration        Id = "383"
	RenewableEnergyPowerGeneration Id = "3240"
	OilAndGas                      Id = "57"
	Mining                         Id = "56"
	OilGasAndMining                Id = "332"

	// Media & Entertainment
	TechnologyInformationAndMedia           Id = "1594"
	BroadcastMediaProductionAndDistribution Id = "36"
	RadioAndTelevisionBroadcasting          Id = "1633"
	MoviesVideosAndSound                    Id = "35"
	MediaProduction                         Id = "126"
	SoundRecording                          Id = "1623"
	BookAndPeriodicalPublishing             Id = "82"
	NewspaperPublishing                     Id = "81"
	PeriodicalPublishing                    Id = "1600"
	EntertainmentProviders                  Id = "28"
	ArtistsAndWriters                       Id = "38"
	Musicians                               Id = "115"

	// Construction & Real Estate
	Construction               Id = "48"
	CivilEngineering           Id = "51"
	RealEstate                 Id = "44"
	RealEstateAgentsAndBrokers Id = "1770"

	// Hospitality & Services
	Hospitality                Id = "31"
	HotelsAndMotels            Id = "2194"
	Restaurants                Id = "32"
	FoodAndBeverageServices    Id = "34"
	TravelArrangements         Id = "30"
	EventsServices             Id = "110"
	WellnessAndFitnessServices Id = "124"
	ConsumerServices           Id = "91"

	// Government & Non-Profit
	ArmedForces                 Id = "71"
	GovernmentRelationsServices Id = "148"
	NonProfitOrganizations      Id = "100"
	CivicAndSocialOrganizations Id = "90"
	PoliticalOrganizations      Id = "107"
	ProfessionalOrganizations   Id = "1911"
	Fundraising                 Id = "101"

	// Wholesale & Distribution
	Wholesale                      Id = "133"
	WholesaleImportAndExport       Id = "134"
	WholesaleComputerEquipment     Id = "1157"
	WholesaleFoodAndBeverage       Id = "1231"
	WholesaleBuildingMaterials     Id = "49"
	WholesaleMachinery             Id = "1187"
	WholesaleMotorVehiclesAndParts Id = "1128"

	// Other Services
	StaffingAndRecruiting     Id = "104"
	ExecutiveSearchServices   Id = "1923"
	OfficeAdministration      Id = "1916"
	SecurityAndInvestigations Id = "121"
	EquipmentRentalServices   Id = "1779"
	Libraries                 Id = "85"
)

var All = util.NewSet(
	// Technology & Software
	SoftwareDevelopment,
	ComputerHardwareManufacturing,
	ComputerNetworkingProducts,
	ItServicesAndItConsulting,
	ComputerAndNetworkSecurity,
	Telecommunications,
	WirelessServices,
	TechnologyInformationAndInternet,
	DataInfrastructureAndAnalytics,
	InformationServices,
	InternetPublishing,
	SocialNetworkingPlatforms,
	ComputerGames,
	MobileGamingApps,
	BlockchainServices,
	BusinessIntelligencePlatforms,

	// Financial Services
	FinancialServices,
	Banking,
	Insurance,
	InvestmentBanking,
	CapitalMarkets,
	VentureCapitalAndPrivateEquityPrincipals,
	SecuritiesAndCommodityExchanges,
	FundsAndTrusts,

	// Healthcare & Medical
	Hospitals,
	MedicalPractices,
	MedicalEquipmentManufacturing,
	PublicHealth,
	VeterinaryServices,
	BiotechnologyResearch,

	// Manufacturing
	Manufacturing,
	ComputersAndElectronicsManufacturing,
	SemiconductorManufacturing,
	MachineryManufacturing,
	IndustrialMachineryManufacturing,
	FoodAndBeverageManufacturing,
	TextileManufacturing,
	MotorVehicleManufacturing,
	MotorVehiclePartsManufacturing,
	AviationAndAerospaceComponentManufacturing,
	DefenseAndSpaceManufacturing,
	PlasticsManufacturing,
	RubberProductsManufacturing,
	PaperAndForestProductManufacturing,
	WoodProductManufacturing,
	FurnitureAndHomeFurnishingsManufacturing,
	SportingGoodsManufacturing,
	PrintingServices,

	// Retail & Consumer Goods
	Retail,
	RetailGroceries,
	OnlineAndMailOrderRetail,
	RetailApparelAndFashion,
	RetailAppliancesElectricalAndElectronicEquipment,
	RetailBooksAndPrintedNews,
	RetailBuildingMaterialsAndGardenEquipment,
	RetailFurnitureAndHomeFurnishings,
	RetailHealthAndPersonalCareProducts,
	RetailLuxuryGoodsAndJewelry,
	RetailMotorVehicles,
	RetailOfficeEquipment,
	RetailOfficeSuppliesAndGifts,

	// Professional Services
	ProfessionalServices,
	Accounting,
	LegalServices,
	LawPractice,
	BusinessConsultingAndServices,
	StrategicManagementServices,
	HumanResourcesServices,
	MarketingServices,
	AdvertisingServices,
	PublicRelationsAndCommunicationsServices,
	MarketResearch,
	ArchitectureAndPlanning,
	DesignServices,
	GraphicDesign,
	InteriorDesign,
	EngineeringServices,
	EnvironmentalServices,
	ResearchServices,
	ThinkTanks,
	Photography,
	TranslationAndLocalization,
	WritingAndEditing,

	// Education
	Education,
	HigherEducation,
	ProfessionalTrainingAndCoaching,
	SportsAndRecreationInstruction,

	// Transportation & Logistics
	TransportationLogisticsSupplyChainAndStorage,
	AirlinesAndAviation,
	FreightAndPackageTransportation,
	MaritimeTransportation,
	RailTransportation,
	TruckTransportation,
	WarehousingAndStorage,
	PostalServices,

	// Energy & Utilities
	Utilities,
	ElectricPowerGeneration,
	RenewableEnergyPowerGeneration,
	OilAndGas,
	Mining,
	OilGasAndMining,

	// Media & Entertainment
	TechnologyInformationAndMedia,
	BroadcastMediaProductionAndDistribution,
	RadioAndTelevisionBroadcasting,
	MoviesVideosAndSound,
	MediaProduction,
	SoundRecording,
	BookAndPeriodicalPublishing,
	NewspaperPublishing,
	PeriodicalPublishing,
	EntertainmentProviders,
	ArtistsAndWriters,
	Musicians,

	// Construction & Real Estate
	Construction,
	CivilEngineering,
	RealEstate,
	RealEstateAgentsAndBrokers,

	// Hospitality & Services
	Hospitality,
	HotelsAndMotels,
	Restaurants,
	FoodAndBeverageServices,
	TravelArrangements,
	EventsServices,
	WellnessAndFitnessServices,
	ConsumerServices,

	// Government & Non-Profit
	ArmedForces,
	GovernmentRelationsServices,
	NonProfitOrganizations,
	CivicAndSocialOrganizations,
	PoliticalOrganizations,
	ProfessionalOrganizations,
	Fundraising,

	// Wholesale & Distribution
	Wholesale,
	WholesaleImportAndExport,
	WholesaleComputerEquipment,
	WholesaleFoodAndBeverage,
	WholesaleBuildingMaterials,
	WholesaleMachinery,
	WholesaleMotorVehiclesAndParts,

	// Other Services
	StaffingAndRecruiting,
	ExecutiveSearchServices,
	OfficeAdministration,
	SecurityAndInvestigations,
	EquipmentRentalServices,
	Libraries,
)

type IndustriesConfig struct {
	All util.Set[Id]
	// Technology & Software
	SoftwareDevelopment              Id
	ComputerHardwareManufacturing    Id
	ComputerNetworkingProducts       Id
	ItServicesAndItConsulting        Id
	ComputerAndNetworkSecurity       Id
	Telecommunications               Id
	WirelessServices                 Id
	TechnologyInformationAndInternet Id
	DataInfrastructureAndAnalytics   Id
	InformationServices              Id
	InternetPublishing               Id
	SocialNetworkingPlatforms        Id
	ComputerGames                    Id
	MobileGamingApps                 Id
	BlockchainServices               Id
	BusinessIntelligencePlatforms    Id

	// Financial Services
	FinancialServices                        Id
	Banking                                  Id
	Insurance                                Id
	InvestmentBanking                        Id
	CapitalMarkets                           Id
	VentureCapitalAndPrivateEquityPrincipals Id
	SecuritiesAndCommodityExchanges          Id
	FundsAndTrusts                           Id

	// Healthcare & Medical
	Hospitals                     Id
	MedicalPractices              Id
	MedicalEquipmentManufacturing Id
	PublicHealth                  Id
	VeterinaryServices            Id
	BiotechnologyResearch         Id

	// Manufacturing
	Manufacturing                              Id
	ComputersAndElectronicsManufacturing       Id
	SemiconductorManufacturing                 Id
	MachineryManufacturing                     Id
	IndustrialMachineryManufacturing           Id
	FoodAndBeverageManufacturing               Id
	TextileManufacturing                       Id
	MotorVehicleManufacturing                  Id
	MotorVehiclePartsManufacturing             Id
	AviationAndAerospaceComponentManufacturing Id
	DefenseAndSpaceManufacturing               Id
	PlasticsManufacturing                      Id
	RubberProductsManufacturing                Id
	PaperAndForestProductManufacturing         Id
	WoodProductManufacturing                   Id
	FurnitureAndHomeFurnishingsManufacturing   Id
	SportingGoodsManufacturing                 Id
	PrintingServices                           Id

	// Retail & Consumer Goods
	Retail                                           Id
	RetailGroceries                                  Id
	OnlineAndMailOrderRetail                         Id
	RetailApparelAndFashion                          Id
	RetailAppliancesElectricalAndElectronicEquipment Id
	RetailBooksAndPrintedNews                        Id
	RetailBuildingMaterialsAndGardenEquipment        Id
	RetailFurnitureAndHomeFurnishings                Id
	RetailHealthAndPersonalCareProducts              Id
	RetailLuxuryGoodsAndJewelry                      Id
	RetailMotorVehicles                              Id
	RetailOfficeEquipment                            Id
	RetailOfficeSuppliesAndGifts                     Id

	// Professional Services
	ProfessionalServices                     Id
	Accounting                               Id
	LegalServices                            Id
	LawPractice                              Id
	BusinessConsultingAndServices            Id
	StrategicManagementServices              Id
	HumanResourcesServices                   Id
	MarketingServices                        Id
	AdvertisingServices                      Id
	PublicRelationsAndCommunicationsServices Id
	MarketResearch                           Id
	ArchitectureAndPlanning                  Id
	DesignServices                           Id
	GraphicDesign                            Id
	InteriorDesign                           Id
	EngineeringServices                      Id
	EnvironmentalServices                    Id
	ResearchServices                         Id
	ThinkTanks                               Id
	Photography                              Id
	TranslationAndLocalization               Id
	WritingAndEditing                        Id

	// Education
	Education                       Id
	HigherEducation                 Id
	ProfessionalTrainingAndCoaching Id
	SportsAndRecreationInstruction  Id

	// Transportation & Logistics
	TransportationLogisticsSupplyChainAndStorage Id
	AirlinesAndAviation                          Id
	FreightAndPackageTransportation              Id
	MaritimeTransportation                       Id
	RailTransportation                           Id
	TruckTransportation                          Id
	WarehousingAndStorage                        Id
	PostalServices                               Id

	// Energy & Utilities
	Utilities                      Id
	ElectricPowerGeneration        Id
	RenewableEnergyPowerGeneration Id
	OilAndGas                      Id
	Mining                         Id
	OilGasAndMining                Id

	// Media & Entertainment
	TechnologyInformationAndMedia           Id
	BroadcastMediaProductionAndDistribution Id
	RadioAndTelevisionBroadcasting          Id
	MoviesVideosAndSound                    Id
	MediaProduction                         Id
	SoundRecording                          Id
	BookAndPeriodicalPublishing             Id
	NewspaperPublishing                     Id
	PeriodicalPublishing                    Id
	EntertainmentProviders                  Id
	ArtistsAndWriters                       Id
	Musicians                               Id

	// Construction & Real Estate
	Construction               Id
	CivilEngineering           Id
	RealEstate                 Id
	RealEstateAgentsAndBrokers Id

	// Hospitality & Services
	Hospitality                Id
	HotelsAndMotels            Id
	Restaurants                Id
	FoodAndBeverageServices    Id
	TravelArrangements         Id
	EventsServices             Id
	WellnessAndFitnessServices Id
	ConsumerServices           Id

	// Government & Non-Profit
	ArmedForces                 Id
	GovernmentRelationsServices Id
	NonProfitOrganizations      Id
	CivicAndSocialOrganizations Id
	PoliticalOrganizations      Id
	ProfessionalOrganizations   Id
	Fundraising                 Id

	// Wholesale & Distribution
	Wholesale                      Id
	WholesaleImportAndExport       Id
	WholesaleComputerEquipment     Id
	WholesaleFoodAndBeverage       Id
	WholesaleBuildingMaterials     Id
	WholesaleMachinery             Id
	WholesaleMotorVehiclesAndParts Id

	// Other Services
	StaffingAndRecruiting     Id
	ExecutiveSearchServices   Id
	OfficeAdministration      Id
	SecurityAndInvestigations Id
	EquipmentRentalServices   Id
	Libraries                 Id
}

var Industries = IndustriesConfig{
	All: *All,
	// Technology & Software
	SoftwareDevelopment:              SoftwareDevelopment,
	ComputerHardwareManufacturing:    ComputerHardwareManufacturing,
	ComputerNetworkingProducts:       ComputerNetworkingProducts,
	ItServicesAndItConsulting:        ItServicesAndItConsulting,
	ComputerAndNetworkSecurity:       ComputerAndNetworkSecurity,
	Telecommunications:               Telecommunications,
	WirelessServices:                 WirelessServices,
	TechnologyInformationAndInternet: TechnologyInformationAndInternet,
	DataInfrastructureAndAnalytics:   DataInfrastructureAndAnalytics,
	InformationServices:              InformationServices,
	InternetPublishing:               InternetPublishing,
	SocialNetworkingPlatforms:        SocialNetworkingPlatforms,
	ComputerGames:                    ComputerGames,
	MobileGamingApps:                 MobileGamingApps,
	BlockchainServices:               BlockchainServices,
	BusinessIntelligencePlatforms:    BusinessIntelligencePlatforms,

	// Financial Services
	FinancialServices:                        FinancialServices,
	Banking:                                  Banking,
	Insurance:                                Insurance,
	InvestmentBanking:                        InvestmentBanking,
	CapitalMarkets:                           CapitalMarkets,
	VentureCapitalAndPrivateEquityPrincipals: VentureCapitalAndPrivateEquityPrincipals,
	SecuritiesAndCommodityExchanges:          SecuritiesAndCommodityExchanges,
	FundsAndTrusts:                           FundsAndTrusts,

	// Healthcare & Medical
	Hospitals:                     Hospitals,
	MedicalPractices:              MedicalPractices,
	MedicalEquipmentManufacturing: MedicalEquipmentManufacturing,
	PublicHealth:                  PublicHealth,
	VeterinaryServices:            VeterinaryServices,
	BiotechnologyResearch:         BiotechnologyResearch,

	// Manufacturing
	Manufacturing:                              Manufacturing,
	ComputersAndElectronicsManufacturing:       ComputersAndElectronicsManufacturing,
	SemiconductorManufacturing:                 SemiconductorManufacturing,
	MachineryManufacturing:                     MachineryManufacturing,
	IndustrialMachineryManufacturing:           IndustrialMachineryManufacturing,
	FoodAndBeverageManufacturing:               FoodAndBeverageManufacturing,
	TextileManufacturing:                       TextileManufacturing,
	MotorVehicleManufacturing:                  MotorVehicleManufacturing,
	MotorVehiclePartsManufacturing:             MotorVehiclePartsManufacturing,
	AviationAndAerospaceComponentManufacturing: AviationAndAerospaceComponentManufacturing,
	DefenseAndSpaceManufacturing:               DefenseAndSpaceManufacturing,
	PlasticsManufacturing:                      PlasticsManufacturing,
	RubberProductsManufacturing:                RubberProductsManufacturing,
	PaperAndForestProductManufacturing:         PaperAndForestProductManufacturing,
	WoodProductManufacturing:                   WoodProductManufacturing,
	FurnitureAndHomeFurnishingsManufacturing:   FurnitureAndHomeFurnishingsManufacturing,
	SportingGoodsManufacturing:                 SportingGoodsManufacturing,
	PrintingServices:                           PrintingServices,

	// Retail & Consumer Goods
	Retail:                   Retail,
	RetailGroceries:          RetailGroceries,
	OnlineAndMailOrderRetail: OnlineAndMailOrderRetail,
	RetailApparelAndFashion:  RetailApparelAndFashion,
	RetailAppliancesElectricalAndElectronicEquipment: RetailAppliancesElectricalAndElectronicEquipment,
	RetailBooksAndPrintedNews:                        RetailBooksAndPrintedNews,
	RetailBuildingMaterialsAndGardenEquipment:        RetailBuildingMaterialsAndGardenEquipment,
	RetailFurnitureAndHomeFurnishings:                RetailFurnitureAndHomeFurnishings,
	RetailHealthAndPersonalCareProducts:              RetailHealthAndPersonalCareProducts,
	RetailLuxuryGoodsAndJewelry:                      RetailLuxuryGoodsAndJewelry,
	RetailMotorVehicles:                              RetailMotorVehicles,
	RetailOfficeEquipment:                            RetailOfficeEquipment,
	RetailOfficeSuppliesAndGifts:                     RetailOfficeSuppliesAndGifts,

	// Professional Services
	ProfessionalServices:                     ProfessionalServices,
	Accounting:                               Accounting,
	LegalServices:                            LegalServices,
	LawPractice:                              LawPractice,
	BusinessConsultingAndServices:            BusinessConsultingAndServices,
	StrategicManagementServices:              StrategicManagementServices,
	HumanResourcesServices:                   HumanResourcesServices,
	MarketingServices:                        MarketingServices,
	AdvertisingServices:                      AdvertisingServices,
	PublicRelationsAndCommunicationsServices: PublicRelationsAndCommunicationsServices,
	MarketResearch:                           MarketResearch,
	ArchitectureAndPlanning:                  ArchitectureAndPlanning,
	DesignServices:                           DesignServices,
	GraphicDesign:                            GraphicDesign,
	InteriorDesign:                           InteriorDesign,
	EngineeringServices:                      EngineeringServices,
	EnvironmentalServices:                    EnvironmentalServices,
	ResearchServices:                         ResearchServices,
	ThinkTanks:                               ThinkTanks,
	Photography:                              Photography,
	TranslationAndLocalization:               TranslationAndLocalization,
	WritingAndEditing:                        WritingAndEditing,

	// Education
	Education:                       Education,
	HigherEducation:                 HigherEducation,
	ProfessionalTrainingAndCoaching: ProfessionalTrainingAndCoaching,
	SportsAndRecreationInstruction:  SportsAndRecreationInstruction,

	// Transportation & Logistics
	TransportationLogisticsSupplyChainAndStorage: TransportationLogisticsSupplyChainAndStorage,
	AirlinesAndAviation:                          AirlinesAndAviation,
	FreightAndPackageTransportation:              FreightAndPackageTransportation,
	MaritimeTransportation:                       MaritimeTransportation,
	RailTransportation:                           RailTransportation,
	TruckTransportation:                          TruckTransportation,
	WarehousingAndStorage:                        WarehousingAndStorage,
	PostalServices:                               PostalServices,

	// Energy & Utilities
	Utilities:                      Utilities,
	ElectricPowerGeneration:        ElectricPowerGeneration,
	RenewableEnergyPowerGeneration: RenewableEnergyPowerGeneration,
	OilAndGas:                      OilAndGas,
	Mining:                         Mining,
	OilGasAndMining:                OilGasAndMining,

	// Media & Entertainment
	TechnologyInformationAndMedia:           TechnologyInformationAndMedia,
	BroadcastMediaProductionAndDistribution: BroadcastMediaProductionAndDistribution,
	RadioAndTelevisionBroadcasting:          RadioAndTelevisionBroadcasting,
	MoviesVideosAndSound:                    MoviesVideosAndSound,
	MediaProduction:                         MediaProduction,
	SoundRecording:                          SoundRecording,
	BookAndPeriodicalPublishing:             BookAndPeriodicalPublishing,
	NewspaperPublishing:                     NewspaperPublishing,
	PeriodicalPublishing:                    PeriodicalPublishing,
	EntertainmentProviders:                  EntertainmentProviders,
	ArtistsAndWriters:                       ArtistsAndWriters,
	Musicians:                               Musicians,

	// Construction & Real Estate
	Construction:               Construction,
	CivilEngineering:           CivilEngineering,
	RealEstate:                 RealEstate,
	RealEstateAgentsAndBrokers: RealEstateAgentsAndBrokers,

	// Hospitality & Services
	Hospitality:                Hospitality,
	HotelsAndMotels:            HotelsAndMotels,
	Restaurants:                Restaurants,
	FoodAndBeverageServices:    FoodAndBeverageServices,
	TravelArrangements:         TravelArrangements,
	EventsServices:             EventsServices,
	WellnessAndFitnessServices: WellnessAndFitnessServices,
	ConsumerServices:           ConsumerServices,

	// Government & Non-Profit
	ArmedForces:                 ArmedForces,
	GovernmentRelationsServices: GovernmentRelationsServices,
	NonProfitOrganizations:      NonProfitOrganizations,
	CivicAndSocialOrganizations: CivicAndSocialOrganizations,
	PoliticalOrganizations:      PoliticalOrganizations,
	ProfessionalOrganizations:   ProfessionalOrganizations,
	Fundraising:                 Fundraising,

	// Wholesale & Distribution
	Wholesale:                      Wholesale,
	WholesaleImportAndExport:       WholesaleImportAndExport,
	WholesaleComputerEquipment:     WholesaleComputerEquipment,
	WholesaleFoodAndBeverage:       WholesaleFoodAndBeverage,
	WholesaleBuildingMaterials:     WholesaleBuildingMaterials,
	WholesaleMachinery:             WholesaleMachinery,
	WholesaleMotorVehiclesAndParts: WholesaleMotorVehiclesAndParts,

	// Other Services
	StaffingAndRecruiting:     StaffingAndRecruiting,
	ExecutiveSearchServices:   ExecutiveSearchServices,
	OfficeAdministration:      OfficeAdministration,
	SecurityAndInvestigations: SecurityAndInvestigations,
	EquipmentRentalServices:   EquipmentRentalServices,
	Libraries:                 Libraries,
}
