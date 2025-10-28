package types

import (
	"time"
)

// WebStartURL represents a single start URL configuration for web scraping
type WebStartURL struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

type WebQueryType string

const (
	WebScraper WebQueryType = "scraper"
)

// WebScraperRequest represents the customizable configuration for web scraping operations
type WebScraperRequest struct {
	StartUrls            []WebStartURL `json:"startUrls"`
	MaxCrawlDepth        int           `json:"maxCrawlDepth"`
	MaxCrawlPages        int           `json:"maxCrawlPages"`
	RespectRobotsTxtFile bool          `json:"respectRobotsTxtFile"`
	SaveMarkdown         bool          `json:"saveMarkdown"`
	Proxy                ApifyProxy    `json:"proxyConfiguration"`
}

// WebCrawlInfo contains information about the crawling process
type WebCrawlInfo struct {
	LoadedURL      string    `json:"loadedUrl"`
	LoadedTime     time.Time `json:"loadedTime"`
	ReferrerURL    string    `json:"referrerUrl"`
	Depth          int       `json:"depth"`
	HTTPStatusCode int       `json:"httpStatusCode"`
}

// WebMetadata contains metadata extracted from the scraped page
type WebMetadata struct {
	CanonicalURL string  `json:"canonicalUrl"`
	Title        string  `json:"title"`
	Description  *string `json:"description"`
	Author       *string `json:"author"`
	Keywords     *string `json:"keywords"`
	LanguageCode *string `json:"languageCode"`
}

// WebScraperResult represents the complete result from web scraping a single page
type WebScraperResult struct {
	URL         string       `json:"url"`
	Crawl       WebCrawlInfo `json:"crawl"`
	Metadata    WebMetadata  `json:"metadata"`
	Text        string       `json:"text"`
	Markdown    string       `json:"markdown"`
	LLMResponse string       `json:"llmresponse,omitempty"` // populated by LLM processor
}
