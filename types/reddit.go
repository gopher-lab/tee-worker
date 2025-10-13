package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/masa-finance/tee-types/pkg/util"
)

type RedditQueryType string

const (
	RedditScrapeUrls        RedditQueryType = "scrapeurls"
	RedditSearchPosts       RedditQueryType = "searchposts"
	RedditSearchUsers       RedditQueryType = "searchusers"
	RedditSearchCommunities RedditQueryType = "searchcommunities"
)

var AllRedditQueryTypes = util.NewSet(RedditScrapeUrls, RedditSearchPosts, RedditSearchUsers, RedditSearchCommunities)

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

type RedditTypeSwitch struct {
	Type RedditItemType `json:"type"`
}

type RedditItem struct {
	TypeSwitch *RedditTypeSwitch
	User       *RedditUser
	Post       *RedditPost
	Comment    *RedditComment
	Community  *RedditCommunity
}

func (t *RedditItem) UnmarshalJSON(data []byte) error {
	t.TypeSwitch = &RedditTypeSwitch{}
	if err := json.Unmarshal(data, &t.TypeSwitch); err != nil {
		return fmt.Errorf("failed to unmarshal reddit response type: %w", err)
	}

	switch t.TypeSwitch.Type {
	case RedditUserItem:
		t.User = &RedditUser{}
		if err := json.Unmarshal(data, t.User); err != nil {
			return fmt.Errorf("failed to unmarshal reddit user: %w", err)
		}
	case RedditPostItem:
		t.Post = &RedditPost{}
		if err := json.Unmarshal(data, t.Post); err != nil {
			return fmt.Errorf("failed to unmarshal reddit post: %w", err)
		}
	case RedditCommentItem:
		t.Comment = &RedditComment{}
		if err := json.Unmarshal(data, t.Comment); err != nil {
			return fmt.Errorf("failed to unmarshal reddit comment: %w", err)
		}
	case RedditCommunityItem:
		t.Community = &RedditCommunity{}
		if err := json.Unmarshal(data, t.Community); err != nil {
			return fmt.Errorf("failed to unmarshal reddit community: %w", err)
		}
	default:
		return fmt.Errorf("unknown Reddit response type: %s", t.TypeSwitch.Type)
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for RedditResponse.
// It unwraps the inner struct (User, Post, Comment, or Community) and marshals it directly.
func (t *RedditItem) MarshalJSON() ([]byte, error) {
	if t.TypeSwitch == nil {
		return []byte("null"), nil
	}

	switch t.TypeSwitch.Type {
	case RedditUserItem:
		return json.Marshal(t.User)
	case RedditPostItem:
		return json.Marshal(t.Post)
	case RedditCommentItem:
		return json.Marshal(t.Comment)
	case RedditCommunityItem:
		return json.Marshal(t.Community)
	default:
		return nil, fmt.Errorf("unknown Reddit response type: %s", t.TypeSwitch.Type)
	}
}
