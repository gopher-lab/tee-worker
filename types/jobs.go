package types

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/masa-finance/tee-types/pkg/util"
)

type JobType string

type JobArguments map[string]interface{}

func (j JobArguments) Unmarshal(i interface{}) error {
	d, err := json.Marshal(j)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, i)
}

type Job struct {
	Type         JobType       `json:"type"`
	Arguments    JobArguments  `json:"arguments"`
	UUID         string        `json:"-"`
	Nonce        string        `json:"quote"`
	WorkerID     string        `json:"worker_id"`
	TargetWorker string        `json:"target_worker"`
	Timeout      time.Duration `json:"timeout"`
}

type Capability string
type WorkerCapabilities map[JobType][]Capability

// String returns the string representation of the JobType
func (j JobType) String() string {
	return string(j)
}

// ValidateCapability validates that a capability is supported for this job type
func (j JobType) ValidateCapability(capability Capability) error {
	validCaps, exists := JobCapabilityMap[j]
	if !exists {
		return fmt.Errorf("unknown job type: %s", j)
	}

	if !slices.Contains(validCaps, capability) {
		return fmt.Errorf("capability '%s' is not valid for job type '%s'. valid capabilities: %v",
			capability, j, validCaps)
	}

	return nil
}

// combineCapabilities combines multiple capability slices and ensures uniqueness
func combineCapabilities(capSlices ...[]Capability) []Capability {
	caps := util.NewSet[Capability]()
	for _, capSlice := range capSlices {
		caps.Add(capSlice...)
	}
	return caps.Items()
}

// Job type constants - centralized from tee-indexer and tee-worker
const (
	WebJob               JobType = "web"
	TelemetryJob         JobType = "telemetry"
	TiktokJob            JobType = "tiktok"
	TwitterJob           JobType = "twitter"            // General Twitter scraping (uses best available auth for capability)
	TwitterCredentialJob JobType = "twitter-credential" // Twitter scraping with credentials
	TwitterApiJob        JobType = "twitter-api"        // Twitter scraping with API keys
	TwitterApifyJob      JobType = "twitter-apify"      // Twitter scraping with Apify
	LinkedInJob          JobType = "linkedin"           // LinkedIn scraping, keeping for unmarshalling logic
	RedditJob            JobType = "reddit"             // Reddit scraping with Apify
)

// Capability constants - typed to prevent typos and enable discoverability
const (
	CapScraper             Capability = "scraper"
	CapTelemetry           Capability = "telemetry"
	CapTranscription       Capability = "transcription"
	CapSearchByQuery       Capability = "searchbyquery"
	CapSearchByTrending    Capability = "searchbytrending"
	CapSearchByFullArchive Capability = "searchbyfullarchive"
	CapSearchByProfile     Capability = "searchbyprofile"
	CapGetById             Capability = "getbyid"
	CapGetReplies          Capability = "getreplies"
	CapGetRetweeters       Capability = "getretweeters"
	CapGetTweets           Capability = "gettweets"
	CapGetMedia            Capability = "getmedia"
	CapGetHomeTweets       Capability = "gethometweets"
	CapGetForYouTweets     Capability = "getforyoutweets"
	CapGetProfileById      Capability = "getprofilebyid"
	CapGetTrends           Capability = "gettrends"
	CapGetFollowing        Capability = "getfollowing"
	CapGetFollowers        Capability = "getfollowers"
	CapGetSpace            Capability = "getspace"
	CapGetProfile          Capability = "getprofile"
	// Reddit capabilities
	CapScrapeUrls        Capability = "scrapeurls"
	CapSearchPosts       Capability = "searchposts"
	CapSearchUsers       Capability = "searchusers"
	CapSearchCommunities Capability = "searchcommunities"

	CapEmpty Capability = ""
)

// Capability group constants for easy reuse
var (
	AlwaysAvailableTelemetryCaps = []Capability{CapTelemetry, CapEmpty}
	AlwaysAvailableTiktokCaps    = []Capability{CapTranscription, CapEmpty}
	AlwaysAvailableLinkedInCaps  = []Capability{CapSearchByQuery, CapGetProfile, CapEmpty}

	// AlwaysAvailableCapabilities defines the job capabilities that are always available regardless of configuration
	AlwaysAvailableCapabilities = WorkerCapabilities{
		TelemetryJob: AlwaysAvailableTelemetryCaps,
		TiktokJob:    AlwaysAvailableTiktokCaps,
	}

	// TwitterCredentialCaps are all Twitter capabilities available with credential-based auth
	TwitterCredentialCaps = []Capability{
		CapSearchByQuery, CapSearchByProfile,
		CapGetById, CapGetReplies, CapGetRetweeters, CapGetTweets, CapGetMedia,
		CapGetHomeTweets, CapGetForYouTweets, CapGetProfileById,
		CapGetTrends, CapGetFollowing, CapGetFollowers, CapGetSpace,
		CapEmpty,
	}

	// TwitterAPICaps are basic Twitter capabilities available with API keys
	TwitterAPICaps = []Capability{CapSearchByQuery, CapGetById, CapGetProfileById, CapEmpty}

	// TwitterApifyCaps are Twitter capabilities available with Apify
	TwitterApifyCaps = []Capability{CapGetFollowers, CapGetFollowing, CapEmpty}

	// TiktokSearchCaps are Tiktok capabilities available with Apify
	TiktokSearchCaps = []Capability{CapSearchByQuery, CapSearchByTrending}

	// RedditCaps are all the Reddit capabilities (only available with Apify)
	RedditCaps = []Capability{CapScrapeUrls, CapSearchPosts, CapSearchUsers, CapSearchCommunities}

	// WebCaps are all the Web capabilities (only available with Apify)
	WebCaps = []Capability{CapScraper, CapEmpty}
)

// JobCapabilityMap defines which capabilities are valid for each job type
var JobCapabilityMap = map[JobType][]Capability{
	// Twitter job types and their valid capabilities
	TwitterJob: combineCapabilities(
		TwitterCredentialCaps,
		TwitterAPICaps,
		TwitterApifyCaps,
		[]Capability{CapSearchByFullArchive}, // Elevated API capability
	),
	TwitterCredentialJob: TwitterCredentialCaps,
	TwitterApiJob: combineCapabilities(
		TwitterAPICaps,
		[]Capability{CapSearchByFullArchive}, // Elevated API capability
	),
	TwitterApifyJob: TwitterApifyCaps,

	// Web job capabilities
	WebJob: WebCaps,

	// TikTok job capabilities
	TiktokJob: combineCapabilities(
		AlwaysAvailableTiktokCaps,
		TiktokSearchCaps,
	),

	// Reddit job capabilities
	RedditJob: RedditCaps,

	// Telemetry job capabilities
	TelemetryJob: AlwaysAvailableTelemetryCaps,
}

// if no capability is specified, use the default capability for the job type
var JobDefaultCapabilityMap = map[JobType]Capability{
	TwitterJob:           CapSearchByQuery,
	TwitterCredentialJob: CapSearchByQuery,
	TwitterApiJob:        CapSearchByQuery,
	TwitterApifyJob:      CapGetFollowers,
	WebJob:               CapScraper,
	TiktokJob:            CapTranscription,
	RedditJob:            CapScrapeUrls,
	TelemetryJob:         CapTelemetry,
}
