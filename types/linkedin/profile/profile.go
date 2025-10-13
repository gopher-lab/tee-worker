package profile

import (
	"time"

	"github.com/masa-finance/tee-types/pkg/util"
)

type ScraperMode string

const (
	ScraperModeShort     ScraperMode = "Short"
	ScraperModeFull      ScraperMode = "Full"
	ScraperModeFullEmail ScraperMode = "Full + email search"
)

var AllScraperModes = util.NewSet(ScraperModeShort, ScraperModeFull, ScraperModeFullEmail)

// Profile represents a complete profile response
type Profile struct {
	ID                string    `json:"id"`
	PublicIdentifier  string    `json:"publicIdentifier,omitempty"`
	URL               string    `json:"linkedinUrl"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Headline          *string   `json:"headline,omitempty"`
	About             *string   `json:"about,omitempty"`
	Summary           *string   `json:"summary,omitempty"`
	OpenToWork        bool      `json:"openToWork,omitempty"`
	OpenProfile       bool      `json:"openProfile,omitempty"`
	Hiring            bool      `json:"hiring,omitempty"`
	Photo             *string   `json:"photo,omitempty"`
	PictureUrl        *string   `json:"pictureUrl,omitempty"`
	Premium           bool      `json:"premium,omitempty"`
	Influencer        bool      `json:"influencer,omitempty"`
	Location          Location  `json:"location,omitempty"`
	Verified          bool      `json:"verified,omitempty"`
	RegisteredAt      time.Time `json:"registeredAt,omitempty"`
	TopSkills         *string   `json:"topSkills,omitempty"`
	ConnectionsCount  int       `json:"connectionsCount,omitempty"`
	FollowerCount     int       `json:"followerCount,omitempty"`
	ComposeOptionType *string   `json:"composeOptionType,omitempty"`

	// Full mode
	CurrentPosition []CurrentPosition `json:"currentPosition,omitempty"`

	// Short mode
	CurrentPositions []ShortCurrentPosition `json:"currentPositions,omitempty"`

	Experience              []Experience     `json:"experience,omitempty"`
	Education               []Education      `json:"education,omitempty"`
	Certifications          []Certification  `json:"certifications,omitempty"`
	Projects                []Project        `json:"projects,omitempty"`
	Volunteering            []Volunteering   `json:"volunteering,omitempty"`
	ReceivedRecommendations []Recommendation `json:"receivedRecommendations,omitempty"`
	Skills                  []Skill          `json:"skills,omitempty"`
	Courses                 []Course         `json:"courses,omitempty"`
	Publications            []Publication    `json:"publications,omitempty"`
	Patents                 []Patent         `json:"patents,omitempty"`
	HonorsAndAwards         []HonorAndAward  `json:"honorsAndAwards,omitempty"`
	Languages               []Language       `json:"languages,omitempty"`
	Featured                any              `json:"featured,omitempty"`
	MoreProfiles            []MoreProfile    `json:"moreProfiles,omitempty"`

	// Email mode
	Emails          []string         `json:"emails,omitempty"`
	CompanyWebsites []CompanyWebsite `json:"companyWebsites,omitempty"`
}

// Location represents the location information
type Location struct {
	Text        string         `json:"linkedinText"`
	CountryCode string         `json:"countryCode,omitempty"`
	Parsed      ParsedLocation `json:"parsed,omitempty"`
}

// ParsedLocation represents the parsed location details
type ParsedLocation struct {
	Text        string  `json:"text,omitempty"`
	CountryCode string  `json:"countryCode,omitempty"`
	RegionCode  *string `json:"regionCode,omitempty"`
	Country     string  `json:"country,omitempty"`
	CountryFull string  `json:"countryFull,omitempty"`
	State       string  `json:"state,omitempty"`
	City        string  `json:"city,omitempty"`
}

// CurrentPosition represents current position information
type CurrentPosition struct {
	CompanyID          *string     `json:"companyId,omitempty"`
	CompanyLinkedinUrl *string     `json:"companyLinkedinUrl,omitempty"`
	CompanyName        string      `json:"companyName"`
	DateRange          *DatePeriod `json:"dateRange,omitempty"`
}

// Experience represents work experience
type Experience struct {
	Position             string    `json:"position"`
	Location             *string   `json:"location,omitempty"`
	EmploymentType       *string   `json:"employmentType,omitempty"`
	WorkplaceType        *string   `json:"workplaceType,omitempty"`
	CompanyName          string    `json:"companyName"`
	CompanyURL           *string   `json:"companyUrl,omitempty"`
	CompanyID            *string   `json:"companyId,omitempty"`
	CompanyUniversalName *string   `json:"companyUniversalName,omitempty"`
	Duration             string    `json:"duration"`
	Description          *string   `json:"description,omitempty"`
	Skills               []string  `json:"skills,omitempty"`
	StartDate            DateRange `json:"startDate"`
	EndDate              DateRange `json:"endDate"`
}

// DateRange represents a date range with month, year, and text
type DateRange struct {
	Month *string `json:"month,omitempty"`
	Year  *int    `json:"year,omitempty"`
	Text  string  `json:"text"`
}

// Education represents educational background
type Education struct {
	SchoolName   string    `json:"schoolName,omitempty"`
	SchoolURL    string    `json:"schoolUrl,omitempty"`
	Degree       string    `json:"degree,omitempty"`
	FieldOfStudy *string   `json:"fieldOfStudy,omitempty"`
	Skills       []string  `json:"skills,omitempty"`
	StartDate    DateRange `json:"startDate,omitempty"`
	EndDate      DateRange `json:"endDate,omitempty"`
	Period       string    `json:"period,omitempty"`
}

// Certification represents a certification
type Certification struct {
	Title        string `json:"title,omitempty"`
	IssuedAt     string `json:"issuedAt,omitempty"`
	IssuedBy     string `json:"issuedBy,omitempty"`
	IssuedByLink string `json:"issuedByLink,omitempty"`
}

// Project represents a project
type Project struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Duration    string    `json:"duration,omitempty"`
	StartDate   DateRange `json:"startDate,omitempty"`
	EndDate     DateRange `json:"endDate,omitempty"`
}

// Volunteering represents volunteer experience
type Volunteering struct {
	Role             string     `json:"role,omitempty"`
	Duration         string     `json:"duration,omitempty"`
	StartDate        *DateRange `json:"startDate,omitempty"`
	EndDate          *DateRange `json:"endDate,omitempty"`
	OrganizationName string     `json:"organizationName,omitempty"`
	OrganizationURL  *string    `json:"organizationUrl,omitempty"`
	Cause            string     `json:"cause,omitempty"`
}

// Skill represents a skill with optional positions and endorsements
type Skill struct {
	Name         string   `json:"name,omitempty"`
	Positions    []string `json:"positions,omitempty"`
	Endorsements string   `json:"endorsements,omitempty"`
}

// Course represents a course
type Course struct {
	Title              string `json:"title,omitempty"`
	AssociatedWith     string `json:"associatedWith,omitempty"`
	AssociatedWithLink string `json:"associatedWithLink,omitempty"`
}

// Publication represents a publication
type Publication struct {
	Title       string `json:"title,omitempty"`
	PublishedAt string `json:"publishedAt,omitempty"`
	Link        string `json:"link,omitempty"`
}

// HonorAndAward represents an honor or award
type HonorAndAward struct {
	Title              string `json:"title,omitempty"`
	IssuedBy           string `json:"issuedBy,omitempty"`
	IssuedAt           string `json:"issuedAt,omitempty"`
	Description        string `json:"description,omitempty"`
	AssociatedWith     string `json:"associatedWith,omitempty"`
	AssociatedWithLink string `json:"associatedWithLink,omitempty"`
}

// Language represents a language with proficiency level
type Language struct {
	Name        string `json:"name,omitempty"`
	Proficiency string `json:"proficiency,omitempty"`
}

// MoreProfile represents a related profile
type MoreProfile struct {
	ID               string  `json:"id,omitempty"`
	FirstName        string  `json:"firstName,omitempty"`
	LastName         string  `json:"lastName,omitempty"`
	Position         *string `json:"position,omitempty"`
	PublicIdentifier string  `json:"publicIdentifier,omitempty"`
	URL              string  `json:"linkedinUrl,omitempty"`
}

// ShortCurrentPosition represents the short profile current positions array
type ShortCurrentPosition struct {
	TenureAtPosition   *Tenure    `json:"tenureAtPosition,omitempty"`
	CompanyName        string     `json:"companyName,omitempty"`
	Title              *string    `json:"title,omitempty"`
	Current            *bool      `json:"current,omitempty"`
	TenureAtCompany    *Tenure    `json:"tenureAtCompany,omitempty"`
	StartedOn          *StartedOn `json:"startedOn,omitempty"`
	CompanyID          *string    `json:"companyId,omitempty"`
	CompanyLinkedinUrl *string    `json:"companyLinkedinUrl,omitempty"`
}

type Tenure struct {
	NumYears  *int `json:"numYears,omitempty"`
	NumMonths *int `json:"numMonths,omitempty"`
}

type StartedOn struct {
	Month int `json:"month,omitempty"`
	Year  int `json:"year,omitempty"`
}

// DatePeriod represents a date period with optional start and end parts
type DatePeriod struct {
	Start *DateParts `json:"start,omitempty"`
	End   *DateParts `json:"end,omitempty"`
}

type DateParts struct {
	Month *int `json:"month,omitempty"`
	Year  *int `json:"year,omitempty"`
	Day   *int `json:"day,omitempty"`
}

// CompanyWebsite represents company website with validation hint
type CompanyWebsite struct {
	URL              string `json:"url,omitempty"`
	Domain           string `json:"domain,omitempty"`
	ValidEmailServer *bool  `json:"validEmailServer,omitempty"`
}

// Recommendation captures received recommendations
type Recommendation struct {
	GivenBy     *string `json:"givenBy,omitempty"`
	GivenByLink *string `json:"givenByLink,omitempty"`
	GivenAt     *string `json:"givenAt,omitempty"`
	Description string  `json:"description,omitempty"`
}

// Patent represents a patent entry
type Patent struct {
	Title    string  `json:"title,omitempty"`
	Number   *string `json:"number,omitempty"`
	IssuedAt string  `json:"issuedAt,omitempty"`
}
