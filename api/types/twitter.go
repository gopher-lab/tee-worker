// Package types provides shared types between tee-worker and tee-indexer
package types

import "time"

type TweetResult struct {
	ID             int64     `json:"id"`
	TweetID        string    `json:"tweet_id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	Text           string    `json:"text"`
	CreatedAt      time.Time `json:"created_at"`
	Timestamp      int64     `json:"timestamp"`

	ThreadCursor struct {
		FocalTweetID string `json:"focal_tweet_id"`
		ThreadID     string `json:"thread_id"`
		Cursor       string `json:"cursor"`
		CursorType   string `json:"cursor_type"`
	}
	IsQuoted     bool     `json:"is_quoted"`
	IsPin        bool     `json:"is_pin"`
	IsReply      bool     `json:"is_reply"`
	IsRetweet    bool     `json:"is_retweet"`
	IsSelfThread bool     `json:"is_self_thread"`
	Likes        int      `json:"likes"`
	Hashtags     []string `json:"hashtags"`
	HTML         string   `json:"html"`
	Replies      int      `json:"replies"`
	Retweets     int      `json:"retweets"`
	URLs         []string `json:"urls"`
	Username     string   `json:"username"`

	Photos []Photo `json:"photos"`

	// Video type.
	Videos []Video `json:"videos"`

	RetweetedStatusID string `json:"retweeted_status_id"`
	Views             int    `json:"views"`
	SensitiveContent  bool   `json:"sensitive_content"`

	// from twitterx
	AuthorID          string        `json:"author_id"`
	PublicMetrics     PublicMetrics `json:"public_metrics"`
	PossiblySensitive bool          `json:"possibly_sensitive"`
	Lang              string        `json:"lang"`
	NewestID          string        `json:"newest_id"`
	OldestID          string        `json:"oldest_id"`
	ResultCount       int           `json:"result_count"`

	Error error `json:"error"`
}

type PublicMetrics struct {
	RetweetCount    int `json:"retweet_count"`
	ReplyCount      int `json:"reply_count"`
	LikeCount       int `json:"like_count"`
	QuoteCount      int `json:"quote_count"`
	BookmarkCount   int `json:"bookmark_count"`
	ImpressionCount int `json:"impression_count"`
}
type Photo struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Video struct {
	ID      string `json:"id"`
	Preview string `json:"preview"`
	URL     string `json:"url"`
	HLSURL  string `json:"hls_url"`
}

type ProfileResultApify struct {
	ID                             int64           `json:"id"`
	IDStr                          string          `json:"id_str"`
	Name                           string          `json:"name"`
	ScreenName                     string          `json:"screen_name"`
	Location                       string          `json:"location"`
	Description                    string          `json:"description"`
	URL                            *string         `json:"url"`
	Entities                       ProfileEntities `json:"entities"`
	Protected                      bool            `json:"protected"`
	FollowersCount                 int             `json:"followers_count"`
	FastFollowersCount             int             `json:"fast_followers_count"`
	NormalFollowersCount           int             `json:"normal_followers_count"`
	FriendsCount                   int             `json:"friends_count"`
	ListedCount                    int             `json:"listed_count"`
	CreatedAt                      string          `json:"created_at"`
	FavouritesCount                int             `json:"favourites_count"`
	UTCOffset                      *int            `json:"utc_offset"`
	TimeZone                       *string         `json:"time_zone"`
	GeoEnabled                     bool            `json:"geo_enabled"`
	Verified                       bool            `json:"verified"`
	StatusesCount                  int             `json:"statuses_count"`
	MediaCount                     int             `json:"media_count"`
	Lang                           *string         `json:"lang"`
	ContributorsEnabled            bool            `json:"contributors_enabled"`
	IsTranslator                   bool            `json:"is_translator"`
	IsTranslationEnabled           bool            `json:"is_translation_enabled"`
	ProfileBackgroundColor         string          `json:"profile_background_color"`
	ProfileBackgroundImageURL      *string         `json:"profile_background_image_url"`
	ProfileBackgroundImageURLHTTPS *string         `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool            `json:"profile_background_tile"`
	ProfileImageURL                string          `json:"profile_image_url"`
	ProfileImageURLHTTPS           string          `json:"profile_image_url_https"`
	ProfileLinkColor               string          `json:"profile_link_color"`
	ProfileSidebarBorderColor      string          `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string          `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string          `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool            `json:"profile_use_background_image"`
	HasExtendedProfile             bool            `json:"has_extended_profile"`
	DefaultProfile                 bool            `json:"default_profile"`
	DefaultProfileImage            bool            `json:"default_profile_image"`
	PinnedTweetIDs                 []int64         `json:"pinned_tweet_ids"`
	PinnedTweetIDsStr              []string        `json:"pinned_tweet_ids_str"`
	HasCustomTimelines             bool            `json:"has_custom_timelines"`
	CanMediaTag                    bool            `json:"can_media_tag"`
	FollowedBy                     bool            `json:"followed_by"`
	Following                      bool            `json:"following"`
	LiveFollowing                  bool            `json:"live_following"`
	FollowRequestSent              bool            `json:"follow_request_sent"`
	Notifications                  bool            `json:"notifications"`
	Muting                         bool            `json:"muting"`
	Blocking                       bool            `json:"blocking"`
	BlockedBy                      bool            `json:"blocked_by"`
	AdvertiserAccountType          string          `json:"advertiser_account_type"`
	AdvertiserAccountServiceLevels []string        `json:"advertiser_account_service_levels"`
	AnalyticsType                  string          `json:"analytics_type"`
	BusinessProfileState           string          `json:"business_profile_state"`
	TranslatorType                 string          `json:"translator_type"`
	WithheldInCountries            []string        `json:"withheld_in_countries"`
	RequireSomeConsent             bool            `json:"require_some_consent"`
	Type                           string          `json:"type"`
	TargetUsername                 string          `json:"target_username"`
	Email                          *string         `json:"email"`
}

type ProfileEntities struct {
	URL         *URLEntities `json:"url,omitempty"`
	Description *URLEntities `json:"description,omitempty"`
}

type URLEntities struct {
	URLs []URLEntity `json:"urls,omitempty"`
}

type URLEntity struct {
	URL         string `json:"url"`
	ExpandedURL string `json:"expanded_url"`
	DisplayURL  string `json:"display_url"`
	Indices     []int  `json:"indices"`
}

type ProfileResultScraper struct {
	Avatar               string     `json:"avatar"`
	Banner               string     `json:"banner"`
	Biography            string     `json:"biography"`
	Birthday             string     `json:"birthday"`
	FollowersCount       int        `json:"followers_count"`
	FollowingCount       int        `json:"following_count"`
	FriendsCount         int        `json:"friends_count"`
	IsPrivate            bool       `json:"is_private"`
	IsVerified           bool       `json:"is_verified"`
	IsBlueVerified       bool       `json:"is_blue_verified"`
	Joined               *time.Time `json:"joined"`
	LikesCount           int        `json:"likes_count"`
	ListedCount          int        `json:"listed_count"`
	Location             string     `json:"location"`
	Name                 string     `json:"name"`
	PinnedTweetIDs       []string   `json:"pinned_tweet_ids"`
	TweetsCount          int        `json:"tweets_count"`
	URL                  string     `json:"url"`
	UserID               string     `json:"user_id"`
	Username             string     `json:"username"`
	Website              string     `json:"website"`
	Sensitive            bool       `json:"sensitive"`
	Following            bool       `json:"following"`
	FollowedBy           bool       `json:"followed_by"`
	MediaCount           int        `json:"media_count"`
	FastFollowersCount   int        `json:"fast_followers_count"`
	NormalFollowersCount int        `json:"normal_followers_count"`
	ProfileImageShape    string     `json:"profile_image_shape"`
	HasGraduatedAccess   bool       `json:"has_graduated_access"`
	CanHighlightTweets   bool       `json:"can_highlight_tweets"`
}
