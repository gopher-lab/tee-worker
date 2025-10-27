package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/masa-finance/tee-worker/v2/pkg/util"
)

var AllRedditQueryTypes = util.NewSet(CapScrapeUrls, CapSearchPosts, CapSearchUsers, CapSearchCommunities)

type RedditSortType string

const (
	RedditSortRelevance RedditSortType = "relevance"
	RedditSortHot       RedditSortType = "hot"
	RedditSortTop       RedditSortType = "top"
	RedditSortNew       RedditSortType = "new"
	RedditSortRising    RedditSortType = "rising"
	RedditSortComments  RedditSortType = "comments"
)

var AllRedditSortTypes = util.NewSet(
	RedditSortRelevance,
	RedditSortHot,
	RedditSortTop,
	RedditSortNew,
	RedditSortRising,
	RedditSortComments,
)

// RedditStartURL represents a single start URL for the Apify Reddit scraper.
type RedditStartURL struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

type RedditItemType string

const (
	RedditUserItem      RedditItemType = "user"
	RedditPostItem      RedditItemType = "post"
	RedditCommentItem   RedditItemType = "comment"
	RedditCommunityItem RedditItemType = "community"
)

// RedditUser represents the data structure for a Reddit user from the Apify scraper.
type RedditUser struct {
	ID           string    `json:"id"`
	URL          string    `json:"url"`
	Username     string    `json:"username"`
	UserIcon     string    `json:"userIcon"`
	PostKarma    int       `json:"postKarma"`
	CommentKarma int       `json:"commentKarma"`
	Description  string    `json:"description"`
	Over18       bool      `json:"over18"`
	CreatedAt    time.Time `json:"createdAt"`
	ScrapedAt    time.Time `json:"scrapedAt"`
	DataType     string    `json:"dataType"`
}

// RedditPost represents the data structure for a Reddit post from the Apify scraper.
type RedditPost struct {
	ID                  string    `json:"id"`
	ParsedID            string    `json:"parsedId"`
	URL                 string    `json:"url"`
	Username            string    `json:"username"`
	Title               string    `json:"title"`
	CommunityName       string    `json:"communityName"`
	ParsedCommunityName string    `json:"parsedCommunityName"`
	Body                string    `json:"body"`
	HTML                *string   `json:"html"`
	NumberOfComments    int       `json:"numberOfComments"`
	UpVotes             int       `json:"upVotes"`
	IsVideo             bool      `json:"isVideo"`
	IsAd                bool      `json:"isAd"`
	Over18              bool      `json:"over18"`
	CreatedAt           time.Time `json:"createdAt"`
	ScrapedAt           time.Time `json:"scrapedAt"`
	DataType            string    `json:"dataType"`
}

// RedditComment represents the data structure for a Reddit comment from the Apify scraper.
type RedditComment struct {
	ID              string    `json:"id"`
	ParsedID        string    `json:"parsedId"`
	URL             string    `json:"url"`
	ParentID        string    `json:"parentId"`
	Username        string    `json:"username"`
	Category        string    `json:"category"`
	CommunityName   string    `json:"communityName"`
	Body            string    `json:"body"`
	CreatedAt       time.Time `json:"createdAt"`
	ScrapedAt       time.Time `json:"scrapedAt"`
	UpVotes         int       `json:"upVotes"`
	NumberOfReplies int       `json:"numberOfreplies"`
	HTML            string    `json:"html"`
	DataType        string    `json:"dataType"`
}

// RedditCommunity represents the data structure for a Reddit community from the Apify scraper.
type RedditCommunity struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Title           string    `json:"title"`
	HeaderImage     string    `json:"headerImage"`
	Description     string    `json:"description"`
	Over18          bool      `json:"over18"`
	CreatedAt       time.Time `json:"createdAt"`
	ScrapedAt       time.Time `json:"scrapedAt"`
	NumberOfMembers int       `json:"numberOfMembers"`
	URL             string    `json:"url"`
	DataType        string    `json:"dataType"`
}

// RedditResponse represents a Reddit API response that can be any of the Reddit item types
type RedditResponse struct {
	Type      RedditItemType   `json:"type"`
	User      *RedditUser      `json:"user,omitempty"`
	Post      *RedditPost      `json:"post,omitempty"`
	Comment   *RedditComment   `json:"comment,omitempty"`
	Community *RedditCommunity `json:"community,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for RedditResponse
func (r *RedditResponse) UnmarshalJSON(data []byte) error {
	// First, unmarshal into a map to get the type
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Get the type field (check both 'type' and 'dataType' for compatibility)
	var itemType RedditItemType
	if typeData, exists := raw["type"]; exists {
		if err := json.Unmarshal(typeData, &itemType); err != nil {
			return fmt.Errorf("failed to unmarshal reddit response type: %w", err)
		}
	} else if typeData, exists := raw["dataType"]; exists {
		if err := json.Unmarshal(typeData, &itemType); err != nil {
			return fmt.Errorf("failed to unmarshal reddit response dataType: %w", err)
		}
	} else {
		return fmt.Errorf("missing 'type' or 'dataType' field in reddit response")
	}

	r.Type = itemType

	// Unmarshal the appropriate struct based on type
	switch itemType {
	case RedditUserItem:
		r.User = &RedditUser{}
		if err := json.Unmarshal(data, r.User); err != nil {
			return fmt.Errorf("failed to unmarshal reddit user: %w", err)
		}
	case RedditPostItem:
		r.Post = &RedditPost{}
		if err := json.Unmarshal(data, r.Post); err != nil {
			return fmt.Errorf("failed to unmarshal reddit post: %w", err)
		}
	case RedditCommentItem:
		r.Comment = &RedditComment{}
		if err := json.Unmarshal(data, r.Comment); err != nil {
			return fmt.Errorf("failed to unmarshal reddit comment: %w", err)
		}
	case RedditCommunityItem:
		r.Community = &RedditCommunity{}
		if err := json.Unmarshal(data, r.Community); err != nil {
			return fmt.Errorf("failed to unmarshal reddit community: %w", err)
		}
	default:
		return fmt.Errorf("unknown Reddit response type: %s", itemType)
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface for RedditResponse.
// It unwraps the inner struct (User, Post, Comment, or Community) and marshals it directly.
func (r *RedditResponse) MarshalJSON() ([]byte, error) {
	switch r.Type {
	case RedditUserItem:
		return json.Marshal(r.User)
	case RedditPostItem:
		return json.Marshal(r.Post)
	case RedditCommentItem:
		return json.Marshal(r.Comment)
	case RedditCommunityItem:
		return json.Marshal(r.Community)
	default:
		return nil, fmt.Errorf("unknown Reddit response type: %s", r.Type)
	}
}

// RedditItem is an alias for RedditResponse for backward compatibility
type RedditItem = RedditResponse
