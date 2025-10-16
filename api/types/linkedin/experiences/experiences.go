package experiences

import "github.com/masa-finance/tee-worker/pkg/util"

// id represents a LinkedIn experience level identifier
type Id string

// Experience level constants
const (
	LessThanAYear    Id = "1"
	OneToTwoYears    Id = "2"
	ThreeToFiveYears Id = "3"
	SixToTenYears    Id = "4"
	MoreThanTenYears Id = "5"
)

var All = util.NewSet(
	LessThanAYear,
	OneToTwoYears,
	ThreeToFiveYears,
	SixToTenYears,
	MoreThanTenYears,
)

type ExperiencesConfig struct {
	All              util.Set[Id]
	LessThanAYear    Id
	OneToTwoYears    Id
	ThreeToFiveYears Id
	SixToTenYears    Id
	MoreThanTenYears Id
}

var Experiences = ExperiencesConfig{
	All:              *All,
	LessThanAYear:    LessThanAYear,
	OneToTwoYears:    OneToTwoYears,
	ThreeToFiveYears: ThreeToFiveYears,
	SixToTenYears:    SixToTenYears,
	MoreThanTenYears: MoreThanTenYears,
}
