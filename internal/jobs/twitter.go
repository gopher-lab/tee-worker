package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/masa-finance/tee-worker/v2/internal/jobs/twitterx"
	"github.com/masa-finance/tee-worker/v2/pkg/client"

	"github.com/masa-finance/tee-worker/v2/api/args"
	twitterargs "github.com/masa-finance/tee-worker/v2/api/args/twitter"
	"github.com/masa-finance/tee-worker/v2/api/types"
	"github.com/masa-finance/tee-worker/v2/internal/config"
	"github.com/masa-finance/tee-worker/v2/internal/jobs/stats"
	"github.com/masa-finance/tee-worker/v2/internal/jobs/twitter"
	"github.com/masa-finance/tee-worker/v2/internal/jobs/twitterapify"

	twitterscraper "github.com/imperatrona/twitter-scraper"
	"github.com/sirupsen/logrus"
)

func (ts *TwitterScraper) convertTwitterScraperTweetToTweetResult(tweet twitterscraper.Tweet) *types.TweetResult {
	id, err := strconv.ParseInt(tweet.ID, 10, 64)
	if err != nil {
		logrus.Warnf("failed to convert tweet ID to int64: %s", tweet.ID)
		id = 0 // set to 0 if conversion fails
	}

	createdAt := time.Unix(tweet.Timestamp, 0).UTC()

	logrus.Debug("Converting Tweet ID: ", id) // Changed to Debug
	return &types.TweetResult{
		ID:             id,
		TweetID:        tweet.ID,
		ConversationID: tweet.ConversationID,
		UserID:         tweet.UserID,
		Text:           tweet.Text,
		CreatedAt:      createdAt,
		Timestamp:      tweet.Timestamp,
		IsQuoted:       tweet.IsQuoted,
		IsPin:          tweet.IsPin,
		IsReply:        tweet.IsReply, // Corrected from tweet.IsPin
		IsRetweet:      tweet.IsRetweet,
		IsSelfThread:   tweet.IsSelfThread,
		Likes:          tweet.Likes,
		Hashtags:       tweet.Hashtags,
		HTML:           tweet.HTML,
		Replies:        tweet.Replies,
		Retweets:       tweet.Retweets,
		URLs:           tweet.URLs,
		Username:       tweet.Username,
		Photos: func() []types.Photo {
			var photos []types.Photo
			for _, photo := range tweet.Photos {
				photos = append(photos, types.Photo{
					ID:  photo.ID,
					URL: photo.URL,
				})
			}
			return photos
		}(),
		Videos: func() []types.Video {
			var videos []types.Video
			for _, video := range tweet.Videos {
				videos = append(videos, types.Video{
					ID:      video.ID,
					Preview: video.Preview,
					URL:     video.URL,
					HLSURL:  video.HLSURL,
				})
			}
			return videos
		}(),
		RetweetedStatusID: tweet.RetweetedStatusID,
		Views:             tweet.Views,
		SensitiveContent:  tweet.SensitiveContent,
	}
}

func parseAccounts(accountPairs []string) []*twitter.TwitterAccount {
	return filterMap(accountPairs, func(pair string) (*twitter.TwitterAccount, bool) {
		credentials := strings.Split(pair, ":")
		if len(credentials) != 2 {
			logrus.Warnf("invalid account credentials: %s", pair)
			return nil, false
		}
		return &twitter.TwitterAccount{
			Username: strings.TrimSpace(credentials[0]),
			Password: strings.TrimSpace(credentials[1]),
		}, true
	})
}

func parseApiKeys(apiKeys []string) []*twitter.TwitterApiKey {
	return filterMap(apiKeys, func(key string) (*twitter.TwitterApiKey, bool) {
		return &twitter.TwitterApiKey{
			Key: strings.TrimSpace(key),
		}, true
	})
}

// getCredentialScraper returns a credential-based scraper and account
func (ts *TwitterScraper) getCredentialScraper(j types.Job, baseDir string) (*twitter.Scraper, *twitter.TwitterAccount, error) {
	if baseDir == "" {
		baseDir = ts.configuration.DataDir
	}

	account := ts.accountManager.GetNextAccount()
	if account == nil {
		ts.statsCollector.Add(j.WorkerID, stats.TwitterAuthErrors, 1)
		return nil, nil, fmt.Errorf("no Twitter credentials available")
	}

	authConfig := twitter.AuthConfig{
		Account: account,
		BaseDir: baseDir,
	}
	scraper := twitter.NewScraper(authConfig)
	if scraper == nil {
		ts.statsCollector.Add(j.WorkerID, stats.TwitterAuthErrors, 1)
		logrus.Errorf("Authentication failed for %s", account.Username)
		return nil, account, fmt.Errorf("twitter authentication failed for %s", account.Username)
	}

	return scraper, account, nil
}

// getApiScraper returns a TwitterX API scraper and API key
func (ts *TwitterScraper) getApiScraper(j types.Job) (*twitterx.TwitterXScraper, *twitter.TwitterApiKey, error) {
	apiKey := ts.accountManager.GetNextApiKey()
	if apiKey == nil {
		ts.statsCollector.Add(j.WorkerID, stats.TwitterAuthErrors, 1)
		return nil, nil, fmt.Errorf("no Twitter API keys available")
	}

	apiClient := client.NewTwitterXClient(apiKey.Key)
	twitterXScraper := twitterx.NewTwitterXScraper(apiClient)

	return twitterXScraper, apiKey, nil
}

// getApifyScraper returns an Apify client
func (ts *TwitterScraper) getApifyScraper(j types.Job) (*twitterapify.TwitterApifyClient, error) {
	if ts.configuration.ApifyApiKey == "" {
		ts.statsCollector.Add(j.WorkerID, stats.TwitterAuthErrors, 1)
		return nil, fmt.Errorf("no Apify API key available")
	}

	apifyScraper, err := twitterapify.NewTwitterApifyClient(ts.configuration.ApifyApiKey)
	if err != nil {
		ts.statsCollector.Add(j.WorkerID, stats.TwitterAuthErrors, 1)
		return nil, fmt.Errorf("failed to create apify scraper: %w", err)
	}
	return apifyScraper, nil
}

func (ts *TwitterScraper) handleError(j types.Job, err error, account *twitter.TwitterAccount) bool {
	if strings.Contains(err.Error(), "Rate limit exceeded") || strings.Contains(err.Error(), "status code 429") {
		ts.statsCollector.Add(j.WorkerID, stats.TwitterRateErrors, 1)
		if account != nil {
			ts.accountManager.MarkAccountRateLimited(account)
			logrus.Warnf("rate limited: %s", account.Username)
		} else {
			logrus.Warn("Rate limited (API Key or no specific account)")
		}
		return true
	}
	ts.statsCollector.Add(j.WorkerID, stats.TwitterErrors, 1)
	return false
}

func filterMap[T any, R any](slice []T, f func(T) (R, bool)) []R {
	result := make([]R, 0, len(slice))
	for _, v := range slice {
		if r, ok := f(v); ok {
			result = append(result, r)
		}
	}
	return result
}

func (ts *TwitterScraper) SearchByProfile(j types.Job, baseDir string, username string) (twitterscraper.Profile, error) {
	scraper, account, err := ts.getCredentialScraper(j, baseDir)
	if err != nil {
		logrus.Errorf("failed to get credential scraper: %v", err)
		return twitterscraper.Profile{}, err
	}
	ts.statsCollector.Add(j.WorkerID, stats.TwitterScrapes, 1)
	profile, err := scraper.GetProfile(username)
	if err != nil {
		logrus.Errorf("scraper.GetProfile failed for username %s: %v", username, err)
		_ = ts.handleError(j, err, account)
		return twitterscraper.Profile{}, err
	}
	ts.statsCollector.Add(j.WorkerID, stats.TwitterProfiles, 1)
	return profile, nil
}

func (ts *TwitterScraper) SearchByQuery(j types.Job, baseDir string, query string, count int) ([]*types.TweetResult, error) {
	scraper, account, err := ts.getCredentialScraper(j, baseDir)
	if err != nil {
		return nil, err
	}
	return ts.scrapeTweetsWithCredentials(j, query, count, scraper, account)
}

func (ts *TwitterScraper) SearchByFullArchive(j types.Job, baseQueryEndpoint string, query string, count int) ([]*types.TweetResult, error) {
	twitterXScraper, apiKey, err := ts.getApiScraper(j)
	if err != nil {
		return nil, err
	}
	return ts.scrapeTweetsWithAPI(j, baseQueryEndpoint, query, count, twitterXScraper, apiKey)
}

func (ts *TwitterScraper) scrapeTweetsWithCredentials(j types.Job, query string, count int, scraper *twitter.Scraper, account *twitter.TwitterAccount) ([]*types.TweetResult, error) {
	ts.statsCollector.Add(j.WorkerID, stats.TwitterScrapes, 1)
	tweets := make([]*types.TweetResult, 0, count)

	ctx, cancel := context.WithTimeout(context.Background(), j.Timeout)
	defer cancel()

	scraper.SetSearchMode(twitterscraper.SearchLatest)

	for tweetScraped := range scraper.SearchTweets(ctx, query, count) {
		if tweetScraped.Error != nil {
			_ = ts.handleError(j, tweetScraped.Error, account)
			return nil, tweetScraped.Error
		}
		newTweetResult := ts.convertTwitterScraperTweetToTweetResult(tweetScraped.Tweet)
		tweets = append(tweets, newTweetResult)
	}

	ts.statsCollector.Add(j.WorkerID, stats.TwitterTweets, uint(len(tweets)))
	return tweets, nil
}

func (ts *TwitterScraper) scrapeTweetsWithAPI(j types.Job, baseQueryEndpoint string, query string, count int, twitterXScraper *twitterx.TwitterXScraper, apiKey *twitter.TwitterApiKey) ([]*types.TweetResult, error) {
	ts.statsCollector.Add(j.WorkerID, stats.TwitterScrapes, 1)

	if baseQueryEndpoint == twitterx.TweetsAll && apiKey.Type == twitter.TwitterApiKeyTypeBase {
		return nil, fmt.Errorf("this API key is a base/Basic key and does not have access to full archive search. Please use an elevated/Pro API key")
	}

	tweets := make([]*types.TweetResult, 0, count)

	cursor := ""
	deadline := time.Now().Add(j.Timeout)

	for len(tweets) < count && time.Now().Before(deadline) {
		numToFetch := count - len(tweets)
		if numToFetch <= 0 {
			break
		}

		result, err := twitterXScraper.ScrapeTweetsByQuery(baseQueryEndpoint, query, numToFetch, cursor)
		if err != nil {
			if ts.handleError(j, err, nil) {
				if len(tweets) > 0 {
					logrus.Warnf("Rate limit hit, returning partial results (%d tweets) for query: %s", len(tweets), query)
					break
				}
			}
			return nil, err
		}

		if result == nil || len(result.Data) == 0 {
			if len(tweets) == 0 {
				logrus.Infof("No tweets found for query: %s with API key.", query)
			}
			break
		}

		for _, tX := range result.Data {
			tweetIDInt, convErr := strconv.ParseInt(tX.ID, 10, 64)
			if convErr != nil {
				logrus.Errorf("Failed to convert tweet ID from twitterx '%s' to int64: %v", tX.ID, convErr)
				return nil, fmt.Errorf("failed to parse tweet ID '%s' from twitterx: %w", tX.ID, convErr)
			}

			newTweet := &types.TweetResult{
				ID:             tweetIDInt,
				TweetID:        tX.ID,
				AuthorID:       tX.AuthorID,
				Text:           tX.Text,
				ConversationID: tX.ConversationID,
				UserID:         tX.AuthorID,
				CreatedAt:      tX.CreatedAt,
				Username:       tX.Username,
				Lang:           tX.Lang,
			}
			//if result.Meta != nil {
			newTweet.NewestID = result.Meta.NewestID
			newTweet.OldestID = result.Meta.OldestID
			newTweet.ResultCount = result.Meta.ResultCount
			//}

			//if tX.PublicMetrics != nil {
			newTweet.PublicMetrics = types.PublicMetrics{
				RetweetCount:  tX.PublicMetrics.RetweetCount,
				ReplyCount:    tX.PublicMetrics.ReplyCount,
				LikeCount:     tX.PublicMetrics.LikeCount,
				QuoteCount:    tX.PublicMetrics.QuoteCount,
				BookmarkCount: tX.PublicMetrics.BookmarkCount,
			}
			//}
			// if tX.PossiblySensitive is available in twitterx.TweetData and types.TweetResult has PossiblySensitive:
			// newTweet.PossiblySensitive = tX.PossiblySensitive
			// Also, fields like IsQuoted, Photos, Videos etc. would need to be populated if tX provides them.
			// Currently, this mapping is simpler than convertTwitterScraperTweetToTweetResult.

			tweets = append(tweets, newTweet)
			if len(tweets) >= count {
				goto EndLoop
			}
		}

		if result.Meta.NextCursor != "" {
			cursor = result.Meta.NextCursor
		} else {
			cursor = ""
		}

		if cursor == "" {
			break
		}
	}
EndLoop:

	logrus.Infof("Scraped %d tweets (target: %d) using API key for query: %s", len(tweets), count, query)
	ts.statsCollector.Add(j.WorkerID, stats.TwitterTweets, uint(len(tweets)))
	return tweets, nil
}

func (ts *TwitterScraper) GetTweet(j types.Job, baseDir, tweetID string) (*types.TweetResult, error) {
	scraper, account, err := ts.getCredentialScraper(j, baseDir)
	if err != nil {
		return nil, err
	}

	ts.statsCollector.Add(j.WorkerID, stats.TwitterScrapes, 1)
	scrapedTweet, err := scraper.GetTweet(tweetID)
	if err != nil {
		_ = ts.handleError(j, err, account)
		return nil, err
	}
	if scrapedTweet == nil {
		return nil, fmt.Errorf("scrapedTweet not found or error occurred, but error was nil")
	}
	tweetResult := ts.convertTwitterScraperTweetToTweetResult(*scrapedTweet)
	ts.statsCollector.Add(j.WorkerID, stats.TwitterTweets, 1)
	return tweetResult, nil
}

func (ts *TwitterScraper) GetTweetReplies(j types.Job, baseDir, tweetID string, cursor string) ([]*types.TweetResult, error) {
	scraper, account, err := ts.getCredentialScraper(j, baseDir)
	if err != nil {
		return nil, err
	}

	ts.statsCollector.Add(j.WorkerID, stats.TwitterScrapes, 1)
	var replies []*types.TweetResult

	scrapedTweets, threadEntries, err := scraper.GetTweetReplies(tweetID, cursor)
	if err != nil {
		_ = ts.handleError(j, err, account)
		return nil, err
	}

	for i, scrapedTweet := range scrapedTweets {
		newTweetResult := ts.convertTwitterScraperTweetToTweetResult(*scrapedTweet)
		if i < len(threadEntries) {
			// Assuming types.TweetResult has a ThreadCursor field (struct, not pointer)
			newTweetResult.ThreadCursor.Cursor = threadEntries[i].Cursor
			newTweetResult.ThreadCursor.CursorType = threadEntries[i].CursorType
			newTweetResult.ThreadCursor.FocalTweetID = threadEntries[i].FocalTweetID
			newTweetResult.ThreadCursor.ThreadID = threadEntries[i].ThreadID
		}
		// Removed newTweetResult.Error = err as err is for the GetTweetReplies operation itself.
		replies = append(replies, newTweetResult)
	}

	ts.statsCollector.Add(j.WorkerID, stats.TwitterTweets, uint(len(replies)))
	return replies, nil
}

func (ts *TwitterScraper) GetTweetRetweeters(j types.Job, baseDir, tweetID string, count int, cursor string) ([]*twitterscraper.Profile, error) {
	scraper, account, err := ts.getCredentialScraper(j, baseDir)
	if err != nil {
		return nil, err
	}

	ts.statsCollector.Add(j.WorkerID, stats.TwitterScrapes, 1)
	retweeters, _, err := scraper.GetTweetRetweeters(tweetID, count, cursor)
	if err != nil {
		_ = ts.handleError(j, err, account)
		return nil, err
	}

	ts.statsCollector.Add(j.WorkerID, stats.TwitterProfiles, uint(len(retweeters)))
	return retweeters, nil
}

func (ts *TwitterScraper) GetUserTweets(j types.Job, baseDir, username string, count int, cursor string) ([]*types.TweetResult, string, error) {
	scraper, account, err := ts.getCredentialScraper(j, baseDir)
	if err != nil {
		return nil, "", err
	}
	ts.statsCollector.Add(j.WorkerID, stats.TwitterScrapes, 1)

	var tweets []*types.TweetResult
	var nextCursor string

	if cursor != "" {
		fetchedTweets, fetchCursor, fetchErr := scraper.FetchTweets(username, count, cursor)
		if fetchErr != nil {
			_ = ts.handleError(j, fetchErr, account)
			return nil, "", fetchErr
		}
		for _, tweet := range fetchedTweets {
			newTweetResult := ts.convertTwitterScraperTweetToTweetResult(*tweet)
			tweets = append(tweets, newTweetResult)
		}
		nextCursor = fetchCursor
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), j.Timeout)
		defer cancel()
		for tweetScraped := range scraper.GetTweets(ctx, username, count) {
			if tweetScraped.Error != nil {
				_ = ts.handleError(j, tweetScraped.Error, account)
				return nil, "", tweetScraped.Error
			}
			newTweetResult := ts.convertTwitterScraperTweetToTweetResult(tweetScraped.Tweet)
			tweets = append(tweets, newTweetResult)
		}
		if len(tweets) > 0 {
			nextCursor = strconv.FormatInt(tweets[len(tweets)-1].ID, 10)
		}
	}
	ts.statsCollector.Add(j.WorkerID, stats.TwitterTweets, uint(len(tweets)))
	return tweets, nextCursor, nil
}

func (ts *TwitterScraper) GetUserMedia(j types.Job, baseDir, username string, count int, cursor string) ([]*types.TweetResult, string, error) {
	scraper, account, err := ts.getCredentialScraper(j, baseDir)
	if err != nil {
		return nil, "", err
	}
	ts.statsCollector.Add(j.WorkerID, stats.TwitterScrapes, 1)

	var media []*types.TweetResult
	var nextCursor string
	ctx, cancel := context.WithTimeout(context.Background(), j.Timeout)
	defer cancel()

	if cursor != "" {
		fetchedTweets, fetchCursor, fetchErr := scraper.FetchTweetsAndReplies(username, count, cursor)
		if fetchErr != nil {
			_ = ts.handleError(j, fetchErr, account)
			return nil, "", fetchErr
		}
		for _, tweet := range fetchedTweets {
			if len(tweet.Photos) > 0 || len(tweet.Videos) > 0 {
				newTweetResult := ts.convertTwitterScraperTweetToTweetResult(*tweet)
				media = append(media, newTweetResult)
			}
			if len(media) >= count {
				break
			}
		}
		nextCursor = fetchCursor
	} else {
		// Fetch more tweets initially as GetTweetsAndReplies doesn't guarantee 'count' media items.
		// Adjust multiplier as needed; it's a heuristic.
		initialFetchCount := count * 5
		if initialFetchCount == 0 && count > 0 { // handle count=0 case for initialFetchCount if count is very small
			initialFetchCount = 100 // a reasonable default if count is tiny but non-zero
		} else if count == 0 {
			initialFetchCount = 0 // if specifically asking for 0 media items
		}

		for tweetScraped := range scraper.GetTweetsAndReplies(ctx, username, initialFetchCount) {
			if tweetScraped.Error != nil {
				if ts.handleError(j, tweetScraped.Error, account) {
					return nil, "", tweetScraped.Error
				}
				continue
			}
			if len(tweetScraped.Tweet.Photos) > 0 || len(tweetScraped.Tweet.Videos) > 0 {
				newTweetResult := ts.convertTwitterScraperTweetToTweetResult(tweetScraped.Tweet)
				media = append(media, newTweetResult)
				if len(media) >= count && count > 0 { // ensure count > 0 for breaking
					break
				}
			}
		}
		if len(media) > 0 {
			nextCursor = strconv.FormatInt(media[len(media)-1].ID, 10)
		}
	}
	ts.statsCollector.Add(j.WorkerID, stats.TwitterOther, uint(len(media)))
	return media, nextCursor, nil
}

func (ts *TwitterScraper) GetProfileByID(j types.Job, baseDir, userID string) (*twitterscraper.Profile, error) {
	scraper, account, err := ts.getCredentialScraper(j, baseDir)
	if err != nil {
		return nil, err
	}

	ts.statsCollector.Add(j.WorkerID, stats.TwitterScrapes, 1)
	profile, err := scraper.GetProfileByID(userID)
	if err != nil {
		_ = ts.handleError(j, err, account)
		return nil, err
	}
	ts.statsCollector.Add(j.WorkerID, stats.TwitterProfiles, 1)
	return &profile, nil
}

func (ts *TwitterScraper) GetTrends(j types.Job, baseDir string) ([]string, error) {
	scraper, account, err := ts.getCredentialScraper(j, baseDir)
	if err != nil {
		return nil, err
	}

	ts.statsCollector.Add(j.WorkerID, stats.TwitterScrapes, 1)
	trends, err := scraper.GetTrends()
	if err != nil {
		_ = ts.handleError(j, err, account)
		return nil, err
	}
	ts.statsCollector.Add(j.WorkerID, stats.TwitterOther, uint(len(trends)))
	return trends, nil
}

func (ts *TwitterScraper) getFollowersApify(j types.Job, username string, maxResults uint, cursor client.Cursor) ([]*types.ProfileResultApify, client.Cursor, error) {
	apifyScraper, err := ts.getApifyScraper(j)
	if err != nil {
		return nil, "", err
	}

	ts.statsCollector.Add(j.WorkerID, stats.TwitterScrapes, 1)

	followers, nextCursor, err := apifyScraper.GetFollowers(username, maxResults, cursor)
	if err != nil {
		return nil, "", err
	}

	ts.statsCollector.Add(j.WorkerID, stats.TwitterFollowers, uint(len(followers)))
	return followers, nextCursor, nil
}

func (ts *TwitterScraper) getFollowingApify(j types.Job, username string, maxResults uint, cursor client.Cursor) ([]*types.ProfileResultApify, client.Cursor, error) {
	apifyScraper, err := ts.getApifyScraper(j)
	if err != nil {
		return nil, "", err
	}

	ts.statsCollector.Add(j.WorkerID, stats.TwitterScrapes, 1)

	following, nextCursor, err := apifyScraper.GetFollowing(username, cursor, maxResults)
	if err != nil {
		return nil, "", err
	}

	ts.statsCollector.Add(j.WorkerID, stats.TwitterFollowers, uint(len(following)))
	return following, nextCursor, nil
}

func (ts *TwitterScraper) GetSpace(j types.Job, baseDir, spaceID string) (*twitterscraper.Space, error) {
	scraper, account, err := ts.getCredentialScraper(j, baseDir)
	if err != nil {
		return nil, err
	}

	ts.statsCollector.Add(j.WorkerID, stats.TwitterScrapes, 1)
	space, err := scraper.GetSpace(spaceID)
	if err != nil {
		_ = ts.handleError(j, err, account)
		return nil, err
	}
	ts.statsCollector.Add(j.WorkerID, stats.TwitterOther, 1)
	return space, nil
}

type TwitterScraper struct {
	configuration  config.TwitterScraperConfig
	accountManager *twitter.TwitterAccountManager
	statsCollector *stats.StatsCollector
	capabilities   map[types.Capability]bool
}

func NewTwitterScraper(jc config.JobConfiguration, c *stats.StatsCollector) *TwitterScraper {
	// Use direct config access instead of JSON marshaling/unmarshaling
	config := jc.GetTwitterConfig()

	accounts := parseAccounts(config.Accounts)
	apiKeys := parseApiKeys(config.ApiKeys)
	accountManager := twitter.NewTwitterAccountManager(accounts, apiKeys)
	accountManager.DetectAllApiKeyTypes()

	config.SkipLoginVerification = jc.GetBool("twitter_skip_login_verification", false)

	return &TwitterScraper{
		configuration:  config,
		accountManager: accountManager,
		statsCollector: c,
		capabilities: map[types.Capability]bool{
			// Credential-based capabilities
			types.CapSearchByQuery:   true,
			types.CapSearchByProfile: true,
			types.CapGetById:         true,
			types.CapGetReplies:      true,
			types.CapGetTweets:       true,
			types.CapGetMedia:        true,
			types.CapGetProfileById:  true,
			types.CapGetTrends:       true,
			types.CapGetSpace:        true,
			types.CapGetProfile:      true,

			// API-based capabilities
			types.CapSearchByFullArchive: true,

			// Apify-based capabilities
			types.CapGetFollowing: true,
			types.CapGetFollowers: true,
		},
	}
}

// executeCapability routes the job to the appropriate method based on capability
func (ts *TwitterScraper) executeCapability(j types.Job, jobArgs *twitterargs.SearchArguments) (types.JobResult, error) {
	capability := jobArgs.GetCapability()

	switch capability {
	// Apify-based capabilities
	case types.CapGetFollowers:
		followers, nextCursor, err := ts.getFollowersApify(j, jobArgs.Query, uint(jobArgs.MaxResults), client.Cursor(jobArgs.NextCursor))
		return processResponse(followers, nextCursor.String(), err)
	case types.CapGetFollowing:
		following, nextCursor, err := ts.getFollowingApify(j, jobArgs.Query, uint(jobArgs.MaxResults), client.Cursor(jobArgs.NextCursor))
		return processResponse(following, nextCursor.String(), err)

	// API-based capabilities
	case types.CapSearchByFullArchive:
		tweets, err := ts.SearchByFullArchive(j, twitterx.TweetsAll, jobArgs.Query, jobArgs.MaxResults)
		return processResponse(tweets, "", err)

	// Credential-based capabilities
	case types.CapSearchByQuery:
		tweets, err := ts.SearchByQuery(j, ts.configuration.DataDir, jobArgs.Query, jobArgs.MaxResults)
		return processResponse(tweets, "", err)
	case types.CapSearchByProfile:
		profile, err := ts.SearchByProfile(j, ts.configuration.DataDir, jobArgs.Query)
		return processResponse(profile, "", err)
	case types.CapGetById:
		tweet, err := ts.GetTweet(j, ts.configuration.DataDir, jobArgs.Query)
		return processResponse(tweet, "", err)
	case types.CapGetReplies:
		replies, err := ts.GetTweetReplies(j, ts.configuration.DataDir, jobArgs.Query, jobArgs.NextCursor)
		return processResponse(replies, jobArgs.NextCursor, err)
	case types.CapGetRetweeters:
		retweeters, err := ts.GetTweetRetweeters(j, ts.configuration.DataDir, jobArgs.Query, jobArgs.MaxResults, jobArgs.NextCursor)
		return processResponse(retweeters, jobArgs.NextCursor, err)
	case types.CapGetMedia:
		media, nextCursor, err := ts.GetUserMedia(j, ts.configuration.DataDir, jobArgs.Query, jobArgs.MaxResults, jobArgs.NextCursor)
		return processResponse(media, nextCursor, err)
	case types.CapGetProfileById:
		profile, err := ts.GetProfileByID(j, ts.configuration.DataDir, jobArgs.Query)
		return processResponse(profile, "", err)
	case types.CapGetTrends:
		trends, err := ts.GetTrends(j, ts.configuration.DataDir)
		return processResponse(trends, "", err)
	case types.CapGetSpace:
		space, err := ts.GetSpace(j, ts.configuration.DataDir, jobArgs.Query)
		return processResponse(space, "", err)
	case types.CapGetProfile:
		profile, err := ts.SearchByProfile(j, ts.configuration.DataDir, jobArgs.Query)
		return processResponse(profile, "", err)
	case types.CapGetTweets:
		tweets, nextCursor, err := ts.GetUserTweets(j, ts.configuration.DataDir, jobArgs.Query, jobArgs.MaxResults, jobArgs.NextCursor)
		return processResponse(tweets, nextCursor, err)

	default:
		return types.JobResult{Error: fmt.Sprintf("unsupported capability: %s", capability)}, fmt.Errorf("unsupported capability: %s", capability)
	}
}

func processResponse(response any, nextCursor string, err error) (types.JobResult, error) {
	if err != nil {
		logrus.Debugf("Processing response with error: %v, NextCursor: %s", err, nextCursor)
		return types.JobResult{Error: err.Error(), NextCursor: nextCursor}, err
	}
	dat, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		logrus.Errorf("Error marshalling response: %v", marshalErr)
		return types.JobResult{Error: marshalErr.Error()}, marshalErr
	}
	return types.JobResult{Data: dat, NextCursor: nextCursor}, nil
}

// ExecuteJob runs a Twitter job using capability-based routing.
// It first unmarshals the job arguments using the centralized type-safe unmarshaller.
// Then it routes to the appropriate method based on the capability.
func (ts *TwitterScraper) ExecuteJob(j types.Job) (types.JobResult, error) {
	// Use the centralized unmarshaller from tee-types
	jobArgs, err := args.UnmarshalJobArguments(types.JobType(j.Type), map[string]any(j.Arguments))
	if err != nil {
		logrus.Errorf("Error while unmarshalling job arguments for job ID %s, type %s: %v", j.UUID, j.Type, err)
		return types.JobResult{Error: "error unmarshalling job arguments"}, err
	}

	// Type assert to Twitter arguments
	args, ok := jobArgs.(*twitterargs.SearchArguments)
	if !ok {
		logrus.Errorf("Expected Twitter arguments for job ID %s, type %s", j.UUID, j.Type)
		return types.JobResult{Error: "invalid argument type for Twitter job"}, fmt.Errorf("invalid argument type")
	}

	// Log the capability for debugging
	logrus.Debugf("Executing Twitter job ID %s with capability: %s", j.UUID, args.GetCapability())

	// Route based on capability
	jobResult, err := ts.executeCapability(j, args)
	if err != nil {
		logrus.Errorf("Error executing job ID %s, type %s: %v", j.UUID, j.Type, err)
		return types.JobResult{Error: "error executing job"}, err
	}

	// Check if raw data is empty
	if len(jobResult.Data) == 0 {
		logrus.Errorf("Job result data is empty for job ID %s, type %s", j.UUID, j.Type)
		return types.JobResult{Error: "job result data is empty"}, fmt.Errorf("job result data is empty")
	}

	// Validate the result based on operation type
	switch {
	case args.Type == types.CapGetFollowers || args.Type == types.CapGetFollowing:
		var results []*types.ProfileResultApify
		if err := jobResult.Unmarshal(&results); err != nil {
			logrus.Errorf("Error while unmarshalling followers/following result for job ID %s, type %s: %v", j.UUID, j.Type, err)
			return types.JobResult{Error: "error unmarshalling followers/following result for final validation"}, err
		}
	case args.IsSingleTweetOperation():
		var result *types.TweetResult
		if err := jobResult.Unmarshal(&result); err != nil {
			logrus.Errorf("Error while unmarshalling single tweet result for job ID %s, type %s: %v", j.UUID, j.Type, err)
			return types.JobResult{Error: "error unmarshalling single tweet result for final validation"}, err
		}
	case args.IsMultipleTweetOperation():
		var results []*types.TweetResult
		if err := jobResult.Unmarshal(&results); err != nil {
			logrus.Errorf("Error while unmarshalling multiple tweet result for job ID %s, type %s: %v", j.UUID, j.Type, err)
			return types.JobResult{Error: "error unmarshalling multiple tweet result for final validation"}, err
		}
	case args.IsSingleProfileOperation():
		var result *twitterscraper.Profile
		if err := jobResult.Unmarshal(&result); err != nil {
			logrus.Errorf("Error while unmarshalling single profile result for job ID %s, type %s: %v", j.UUID, j.Type, err)
			return types.JobResult{Error: "error unmarshalling single profile result for final validation"}, err
		}
	case args.IsMultipleProfileOperation():
		var results []*twitterscraper.Profile
		if err := jobResult.Unmarshal(&results); err != nil {
			logrus.Errorf("Error while unmarshalling multiple profile result for job ID %s, type %s: %v", j.UUID, j.Type, err)
			return types.JobResult{Error: "error unmarshalling multiple profile result for final validation"}, err
		}
	case args.IsSingleSpaceOperation():
		var result *twitterscraper.Space
		if err := jobResult.Unmarshal(&result); err != nil {
			logrus.Errorf("Error while unmarshalling single space result for job ID %s, type %s: %v", j.UUID, j.Type, err)
			return types.JobResult{Error: "error unmarshalling single space result for final validation"}, err
		}
	case args.IsTrendsOperation():
		var results []string
		if err := jobResult.Unmarshal(&results); err != nil {
			logrus.Errorf("Error while unmarshalling trends result for job ID %s, type %s: %v", j.UUID, j.Type, err)
			return types.JobResult{Error: "error unmarshalling trends result for final validation"}, err
		}
	default:
		logrus.Errorf("Invalid operation type for job ID %s, type %s", j.UUID, j.Type)
		return types.JobResult{Error: "invalid operation type"}, fmt.Errorf("invalid operation type")
	}

	return jobResult, nil
}
