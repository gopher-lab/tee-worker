// Package types provides shared types between tee-worker and tee-indexer
package types

// TikTokTranscriptionResult defines the structure of the result data for a TikTok transcription
type TikTokTranscriptionResult struct {
	TranscriptionText string `json:"transcription_text"`
	DetectedLanguage  string `json:"detected_language,omitempty"`
	VideoTitle        string `json:"video_title,omitempty"`
	OriginalURL       string `json:"original_url"`
	ThumbnailURL      string `json:"thumbnail_url,omitempty"`
}

type TikTokSearchByQueryResult struct {
	URL                   string            `json:"url"`
	ID                    string            `json:"id"`
	Desc                  string            `json:"desc"`
	CreateTime            string            `json:"createTime"`
	ScheduleTime          int64             `json:"scheduleTime"`
	Video                 TikTokVideo       `json:"video"`
	Author                string            `json:"author"`
	Music                 TikTokMusic       `json:"music"`
	Challenges            []any             `json:"challenges"` // we don't have examples of this data yet...
	Stats                 TikTokStats       `json:"stats"`
	IsActivityItem        bool              `json:"isActivityItem"`
	DuetInfo              TikTokDuetInfo    `json:"duetInfo"`
	WarnInfo              []any             `json:"warnInfo"` // we don't have examples of this data yet...
	OriginalItem          bool              `json:"originalItem"`
	OfficialItem          bool              `json:"officalItem"`
	TextExtra             []TikTokTextExtra `json:"textExtra"`
	Secret                bool              `json:"secret"`
	ForFriend             bool              `json:"forFriend"`
	Digged                bool              `json:"digged"`
	ItemCommentStatus     int               `json:"itemCommentStatus"`
	ShowNotPass           bool              `json:"showNotPass"`
	VL1                   bool              `json:"vl1"`
	TakeDown              int               `json:"takeDown"`
	ItemMute              bool              `json:"itemMute"`
	EffectStickers        []any             `json:"effectStickers"` // we don't have examples of this data yet...
	AuthorStats           TikTokAuthorStats `json:"authorStats"`
	PrivateItem           bool              `json:"privateItem"`
	DuetEnabled           bool              `json:"duetEnabled"`
	StitchEnabled         bool              `json:"stitchEnabled"`
	StickersOnItem        []any             `json:"stickersOnItem"` // we don't have examples of this data yet...
	IsAd                  bool              `json:"isAd"`
	ShareEnabled          bool              `json:"shareEnabled"`
	Comments              []any             `json:"comments"` // we don't have examples of this data yet...
	DuetDisplay           int               `json:"duetDisplay"`
	StitchDisplay         int               `json:"stitchDisplay"`
	IndexEnabled          bool              `json:"indexEnabled"`
	DiversificationLabels []string          `json:"diversificationLabels"`
	AdAuthorization       bool              `json:"adAuthorization"`
	AdLabelVersion        int               `json:"adLabelVersion"`
	LocationCreated       string            `json:"locationCreated"`
	Nickname              string            `json:"nickname"`
	AuthorID              string            `json:"authorId"`
	AuthorSecID           string            `json:"authorSecId"`
	AvatarThumb           string            `json:"avatarThumb"`
	DownloadSetting       int               `json:"downloadSetting"`
	AuthorPrivate         bool              `json:"authorPrivate"`
}

type TikTokSearchByTrending struct {
	CountryCode string `json:"country_code"`
	Cover       string `json:"cover"`
	Duration    int    `json:"duration"`
	ID          string `json:"id"`
	ItemID      string `json:"item_id"`
	ItemURL     string `json:"item_url"`
	Region      string `json:"region"`
	Title       string `json:"title"`
}

type TikTokVideo struct {
	ID            string              `json:"id"`
	Height        int                 `json:"height"`
	Width         int                 `json:"width"`
	Duration      int                 `json:"duration"`
	Ratio         string              `json:"ratio"`
	Cover         string              `json:"cover"`
	OriginCover   string              `json:"origin_cover"`
	DynamicCover  string              `json:"dynamic_cover"`
	PlayAddr      string              `json:"play_addr"`
	DownloadAddr  string              `json:"download_addr"`
	ShareCover    []string            `json:"share_cover"`
	ReflowCover   string              `json:"reflowCover"`
	Bitrate       int                 `json:"bitrate"`
	EncodedType   string              `json:"encodedType"`
	Format        string              `json:"format"`
	VideoQuality  string              `json:"videoQuality"`
	EncodeUserTag string              `json:"encodeUserTag"`
	CodecType     string              `json:"codecType"`
	Definition    string              `json:"definition"`
	SubtitleInfos []any               `json:"subtitleInfos"` // we don't have examples of this data yet...
	ZoomCover     TikTokZoomCover     `json:"zoomCover"`
	VolumeInfo    TikTokVolumeInfo    `json:"volumeInfo"`
	BitrateInfo   []TikTokBitrateInfo `json:"bitrateInfo"`
}

type TikTokZoomCover struct {
	Cover240 string `json:"240"`
	Cover480 string `json:"480"`
	Cover720 string `json:"720"`
	Cover960 string `json:"960"`
}

type TikTokVolumeInfo struct {
	Loudness float64 `json:"Loudness"`
	Peak     float64 `json:"Peak"`
}

type TikTokBitrateInfo struct {
	GearName    string         `json:"GearName"`
	Bitrate     int            `json:"bitrate"`
	QualityType int            `json:"QualityType"`
	PlayAddr    TikTokPlayAddr `json:"PlayAddr"`
	CodecType   string         `json:"CodecType"`
}

type TikTokPlayAddr struct {
	Uri      string   `json:"Uri"`
	UrlList  []string `json:"UrlList"`
	DataSize string   `json:"DataSize"`
	UrlKey   string   `json:"UrlKey"`
	FileHash string   `json:"FileHash"`
	FileCs   string   `json:"FileCs"`
}

type TikTokMusic struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	PlayURL            string `json:"playUrl"`
	CoverLarge         string `json:"coverLarge"`
	CoverMedium        string `json:"coverMedium"`
	CoverThumb         string `json:"coverThumb"`
	AuthorName         string `json:"authorName"`
	Original           bool   `json:"original"`
	Duration           int    `json:"duration"`
	Album              string `json:"album"`
	ScheduleSearchTime int64  `json:"scheduleSearchTime"`
}

type TikTokStats struct {
	DiggCount    int64 `json:"diggCount"`
	ShareCount   int64 `json:"shareCount"`
	CommentCount int64 `json:"commentCount"`
	PlayCount    int64 `json:"playCount"`
}

type TikTokDuetInfo struct {
	DuetFromID string `json:"duetFromId"`
}

type TikTokTextExtra struct {
	AwemeID      string `json:"awemeId"`
	Start        int    `json:"start"`
	End          int    `json:"end"`
	HashtagID    string `json:"hashtagId"`
	HashtagName  string `json:"hashtagName"`
	Type         int    `json:"type"`
	SubType      int    `json:"subType"`
	UserID       string `json:"userId"`
	IsCommerce   bool   `json:"isCommerce"`
	UserUniqueID string `json:"userUniqueId"`
	SecUID       string `json:"secUid"`
}

type TikTokAuthorStats struct {
	FollowerCount  int64 `json:"followerCount"`
	FollowingCount int64 `json:"followingCount"`
	Heart          int64 `json:"heart"`
	HeartCount     int64 `json:"heartCount"`
	VideoCount     int64 `json:"videoCount"`
	DiggCount      int64 `json:"diggCount"`
}
