package apify

import (
	"github.com/masa-finance/tee-worker/api/types"
)

type ActorId string

type actorIds struct {
	RedditScraper         ActorId
	TikTokSearchScraper   ActorId
	TikTokTrendingScraper ActorId
	LLMDatasetProcessor   ActorId
	TwitterFollowers      ActorId
	WebScraper            ActorId
	LinkedInSearchProfile ActorId
}

var ActorIds = actorIds{
	RedditScraper:         "trudax~reddit-scraper",
	TikTokSearchScraper:   "epctex~tiktok-search-scraper",
	TikTokTrendingScraper: "lexis-solutions~tiktok-trending-videos-scraper",
	LLMDatasetProcessor:   "dusan.vystrcil~llm-dataset-processor",
	TwitterFollowers:      "kaitoeasyapi~premium-x-follower-scraper-following-data",
	WebScraper:            "apify~website-content-crawler",
	LinkedInSearchProfile: "harvestapi~linkedin-profile-search",
}

type defaultActorInput map[string]any

type ActorConfig struct {
	ActorId      ActorId
	DefaultInput defaultActorInput
	Capabilities []types.Capability
	JobType      types.JobType
}

// Actors is a list of actor configurations for Apify.  Omitting LLM for now as it's not a standalone actor / has no dedicated capabilities
var Actors = []ActorConfig{
	{
		ActorId:      ActorIds.RedditScraper,
		DefaultInput: defaultActorInput{},
		Capabilities: types.RedditCaps,
		JobType:      types.RedditJob,
	},
	{
		ActorId:      ActorIds.TikTokSearchScraper,
		DefaultInput: defaultActorInput{"proxy": map[string]any{"useApifyProxy": true}},
		Capabilities: []types.Capability{types.CapSearchByQuery},
		JobType:      types.TiktokJob,
	},
	{
		ActorId:      ActorIds.TikTokTrendingScraper,
		DefaultInput: defaultActorInput{},
		Capabilities: []types.Capability{types.CapSearchByTrending},
		JobType:      types.TiktokJob,
	},
	{
		ActorId:      ActorIds.TwitterFollowers,
		DefaultInput: defaultActorInput{"maxFollowers": 200, "maxFollowings": 200},
		Capabilities: []types.Capability{types.CapGetFollowing, types.CapGetFollowers},
		JobType:      types.TwitterJob,
	},
	{
		ActorId:      ActorIds.WebScraper,
		DefaultInput: defaultActorInput{"startUrls": []map[string]any{{"url": "https://docs.learnbittensor.org"}}},
		Capabilities: types.WebCaps,
		JobType:      types.WebJob,
	},
	{
		ActorId:      ActorIds.LinkedInSearchProfile,
		DefaultInput: defaultActorInput{},
		Capabilities: types.LinkedInCaps,
		JobType:      types.LinkedInJob,
	},
}
