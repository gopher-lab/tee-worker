package types

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/masa-finance/tee-worker/pkg/util"
)

type JobStatus string

func (j JobStatus) String() string {
	return string(j)
}

const (
	JobStatusNotSaved   JobStatus = "done(not saved)"
	JobStatusSaved      JobStatus = "done(saved)"
	JobStatusDone       JobStatus = "done"
	JobStatusActive     JobStatus = "in progress"
	JobStatusReceived   JobStatus = "received"
	JobStatusError      JobStatus = "error"
	JobStatusRetryError JobStatus = "error(retrying)"
)

func (j JobStatus) IsDone() bool {
	return j == JobStatusSaved || j == JobStatusDone || j == JobStatusNotSaved
}

// note, this could be combined with job type in a future PR / refactor...
type Source string

func (j Source) String() string {
	return string(j)
}

// To add a new RouterType and/or JobType, you will need to add the JobType to tee-types and the Source below, and add the new mapping to the SourceFor function. This is necessary basically because of Twitter, which has 3 JobTypes but a single Router.
const (
	TwitterSource   Source = "twitter"
	WebSource       Source = "web"
	TiktokSource    Source = "tiktok"
	RedditSource    Source = "reddit"
	LinkedInSource  Source = "linkedin"
	TelemetrySource Source = "telemetry"
	UnknownSource   Source = ""
)

const UnknownJob = JobType("")

var sourceMap = map[JobType]Source{
	TwitterJob:   TwitterSource,
	WebJob:       WebSource,
	TiktokJob:    TiktokSource,
	RedditJob:    RedditSource,
	LinkedInJob:  LinkedInSource,
	TelemetryJob: TelemetrySource,
	UnknownJob:   UnknownSource,
}

var Sources = slices.Compact(slices.Sorted(maps.Values(sourceMap)))

func SourceFor(j JobType) Source {
	source, ok := sourceMap[j]
	if ok {
		return source
	}
	return UnknownSource
}

type JobType string
type Capability string

type WorkerCapabilities map[JobType][]Capability

type JobArguments map[string]any

func (ja JobArguments) Unmarshal(i any) error {
	dat, err := json.Marshal(ja)
	if err != nil {
		return err
	}
	return json.Unmarshal(dat, i)
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

func (j Job) String() string {
	return fmt.Sprintf("UUID: %s Type: %s Arguments: %s", j.UUID, j.Type, j.Arguments)
}

// String returns the string representation of the JobType
func (j JobType) String() string {
	return string(j)
}

// ValidateCapability validates that a capability is supported for this job type
// If the capability is CapEmpty, it will be set to the default capability for the job type
func (j JobType) ValidateCapability(capability *Capability) error {
	// Set default capability if empty
	if *capability == CapEmpty {
		defaultCap, exists := JobDefaultCapabilityMap[j]
		if !exists {
			return fmt.Errorf("no default capability configured for job type: %s", j)
		}
		*capability = defaultCap
	}

	// Validate the capability
	validCaps, exists := JobCapabilityMap[j]
	if !exists {
		return fmt.Errorf("unknown job type: %s", j)
	}

	if !slices.Contains(validCaps, *capability) {
		return fmt.Errorf("capability '%s' is not valid for job type '%s'. valid capabilities: %v",
			*capability, j, validCaps)
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
	WebJob       JobType = "web"
	TelemetryJob JobType = "telemetry"
	TiktokJob    JobType = "tiktok"
	TwitterJob   JobType = "twitter"
	LinkedInJob  JobType = "linkedin"
	RedditJob    JobType = "reddit"
)

// Capability constants - typed to prevent typos and enable discoverability
const (

	// Twitter (credential-based) capabilities
	CapSearchByQuery   Capability = "searchbyquery"
	CapSearchByProfile Capability = "searchbyprofile"
	CapGetById         Capability = "getbyid"
	CapGetReplies      Capability = "getreplies"
	CapGetRetweeters   Capability = "getretweeters"
	CapGetMedia        Capability = "getmedia"
	CapGetProfileById  Capability = "getprofilebyid"
	CapGetTrends       Capability = "gettrends"
	CapGetSpace        Capability = "getspace"
	CapGetProfile      Capability = "getprofile"
	CapGetTweets       Capability = "gettweets"

	// Twitter (apify-based) capabilities
	CapGetFollowing Capability = "getfollowing"
	CapGetFollowers Capability = "getfollowers"

	// Twitter (api-based) capabilities
	CapSearchByFullArchive Capability = "searchbyfullarchive"

	CapScraper          Capability = "scraper"
	CapSearchByTrending Capability = "searchbytrending"
	CapTelemetry        Capability = "telemetry"
	CapTranscription    Capability = "transcription"

	// Reddit capabilities
	CapScrapeUrls        Capability = "scrapeurls"
	CapSearchPosts       Capability = "searchposts"
	CapSearchUsers       Capability = "searchusers"
	CapSearchCommunities Capability = "searchcommunities"

	CapEmpty Capability = ""
)

// Capability group constants for easy reuse
var (
	AlwaysAvailableTelemetryCaps = []Capability{CapTelemetry}
	AlwaysAvailableTiktokCaps    = []Capability{CapTranscription}

	// AlwaysAvailableCapabilities defines the job capabilities that are always available regardless of configuration
	AlwaysAvailableCapabilities = WorkerCapabilities{
		TelemetryJob: AlwaysAvailableTelemetryCaps,
		TiktokJob:    AlwaysAvailableTiktokCaps,
	}

	// Twitter capabilities
	TwitterCaps = []Capability{
		CapSearchByQuery, CapSearchByProfile, CapSearchByFullArchive,
		CapGetById, CapGetReplies, CapGetRetweeters, CapGetTweets, CapGetMedia, CapGetProfileById,
		CapGetTrends, CapGetFollowing, CapGetFollowers, CapGetSpace, CapGetProfile,
	}

	// TiktokSearchCaps are Tiktok capabilities available with Apify
	TiktokSearchCaps = []Capability{CapSearchByQuery, CapSearchByTrending}

	// RedditCaps are all the Reddit capabilities (only available with Apify)
	RedditCaps = []Capability{CapScrapeUrls, CapSearchPosts, CapSearchUsers, CapSearchCommunities}

	// WebCaps are all the Web capabilities (only available with Apify)
	WebCaps = []Capability{CapScraper}

	// LinkedInCaps are all the LinkedIn capabilities (only available with Apify)
	LinkedInCaps = []Capability{CapSearchByProfile}
)

// JobCapabilityMap defines which capabilities are valid for each job type
var JobCapabilityMap = map[JobType][]Capability{
	// Twitter job capabilities
	TwitterJob: TwitterCaps,

	// Web job capabilities
	WebJob: WebCaps,

	// LinkedIn job capabilities
	LinkedInJob: LinkedInCaps,

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
	TwitterJob:   CapSearchByQuery,
	WebJob:       CapScraper,
	TiktokJob:    CapTranscription,
	RedditJob:    CapScrapeUrls,
	TelemetryJob: CapTelemetry,
	LinkedInJob:  CapSearchByProfile,
}

// JobResponse represents a response to a job submission
type JobResponse struct {
	UID string `json:"uid"`
}

// JobResult represents the result of a job execution
type JobResult struct {
	Error      string `json:"error"`
	Data       []byte `json:"data"`
	Job        Job    `json:"job"`
	NextCursor string `json:"next_cursor"`
}

// Success returns true if the job was successful.
func (jr JobResult) Success() bool {
	return jr.Error == ""
}

// Unmarshal unmarshals the job result data.
func (jr JobResult) Unmarshal(i interface{}) error {
	return json.Unmarshal(jr.Data, i)
}

// JobRequest represents a request to execute a job
type JobRequest struct {
	EncryptedJob string `json:"encrypted_job"`
}

// JobError represents an error in job execution
type JobError struct {
	Error string `json:"error"`
}

// Key represents a key request
type Key struct {
	Key       string `json:"key"`
	Signature string `json:"signature"`
}

// KeyResponse represents a response to a key operation
type KeyResponse struct {
	Status string `json:"status"`
}

type ResultResponse struct {
	UUID  string `json:"uuid"`
	Error string `json:"error"`
}

// Document represents a document stored in the vector store. We need to put it in this package because of circular dependencies.
type Document struct {
	Id        string         `json:"id"`
	Source    Source         `json:"source"`
	Content   string         `json:"content"`
	Metadata  map[string]any `json:"metadata"`
	Embedding []float32      `json:"embedding,omitempty"`
	Score     float32        `json:"score,omitempty"` // For similarity search results
	UpdatedAt time.Time      `json:"updated_at"`
	// SearchText is used only for embedding/indexing and SHOULD NOT be serialized or stored.
	SearchText string `json:"-"`
}

func (d Document) String() string {
	return fmt.Sprintf("%s/%s\n%s\n%s", d.Source, d.Id, d.Metadata, d.Content)
}

// CollectionStats represents collection statistics from Milvus
type CollectionStats struct {
	CollectionName string `json:"collection_name,omitempty"`
	RowCount       uint   `json:"row_count"`
}

// JobResult is the struct that is stored in the NATS KV store
type IndexerJobResult struct {
	Status JobStatus  `json:"status"`
	Docs   []Document `json:"docs,omitempty"`
	Error  string     `json:"error"`
}
