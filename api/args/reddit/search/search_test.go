package search_test

import (
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/v2/api/args/reddit"
	"github.com/masa-finance/tee-worker/v2/api/args/reddit/search"
	"github.com/masa-finance/tee-worker/v2/api/types"
)

var _ = Describe("RedditArguments", func() {
	Describe("Marshalling and unmarshalling", func() {
		It("should set default values", func() {
			redditArgs := reddit.NewSearchPostsArguments()
			redditArgs.Queries = []string{"Zaphod", "Ford"}
			jsonData, err := json.Marshal(redditArgs)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &redditArgs)
			Expect(err).ToNot(HaveOccurred())
			Expect(redditArgs.MaxItems).To(Equal(uint(10)))
			Expect(redditArgs.MaxPosts).To(Equal(uint(10)))
			Expect(redditArgs.MaxComments).To(Equal(uint(10)))
			Expect(redditArgs.MaxCommunities).To(Equal(uint(2)))
			Expect(redditArgs.MaxUsers).To(Equal(uint(2)))
			Expect(redditArgs.Sort).To(Equal(types.RedditSortNew))
			Expect(redditArgs.MaxResults).To(Equal(redditArgs.MaxItems))
		})

		It("should override default values", func() {
			redditArgs := reddit.NewSearchPostsArguments()
			redditArgs.Queries = []string{"Zaphod", "Ford"}
			redditArgs.MaxItems = 20
			redditArgs.MaxPosts = 21
			redditArgs.MaxComments = 22
			redditArgs.MaxCommunities = 23
			redditArgs.MaxUsers = 24
			redditArgs.Sort = types.RedditSortTop
			jsonData, err := json.Marshal(redditArgs)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &redditArgs)
			Expect(err).ToNot(HaveOccurred())
			Expect(redditArgs.MaxItems).To(Equal(uint(20)))
			Expect(redditArgs.MaxPosts).To(Equal(uint(21)))
			Expect(redditArgs.MaxComments).To(Equal(uint(22)))
			Expect(redditArgs.MaxCommunities).To(Equal(uint(23)))
			Expect(redditArgs.MaxUsers).To(Equal(uint(24)))
			Expect(redditArgs.MaxResults).To(Equal(uint(20)))
			Expect(redditArgs.Sort).To(Equal(types.RedditSortTop))
		})

	})

	Describe("Validation", func() {
		It("should succeed with valid arguments", func() {
			redditArgs := reddit.NewSearchPostsArguments()
			redditArgs.Queries = []string{"test"}
			redditArgs.Sort = types.RedditSortNew
			err := redditArgs.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should succeed with valid scrapeurls arguments", func() {
			redditArgs := reddit.NewScrapeUrlsArguments()
			redditArgs.URLs = []string{"https://www.reddit.com/r/golang/comments/foo/bar"}
			redditArgs.Sort = types.RedditSortNew
			err := redditArgs.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with an invalid type", func() {
			redditArgs := reddit.NewSearchPostsArguments()
			redditArgs.Type = "invalidtype" // Override the default
			redditArgs.Queries = []string{"test"}
			redditArgs.Sort = types.RedditSortNew
			err := redditArgs.Validate()
			Expect(err).To(MatchError(search.ErrInvalidType))
		})

		It("should fail with an invalid sort", func() {
			redditArgs := search.NewSearchPostsArguments()
			redditArgs.Queries = []string{"test"}
			redditArgs.Sort = "invalidsort"
			err := redditArgs.Validate()
			Expect(err).To(MatchError(search.ErrInvalidSort))
		})

		It("should fail if the after time is in the future", func() {
			redditArgs := reddit.NewSearchPostsArguments()
			redditArgs.Queries = []string{"test"}
			redditArgs.Sort = types.RedditSortNew
			redditArgs.After = time.Now().Add(24 * time.Hour)
			err := redditArgs.Validate()
			Expect(err).To(MatchError(search.ErrTimeInTheFuture))
		})

		It("should fail if queries are not provided for searchposts", func() {
			redditArgs := search.NewSearchPostsArguments()
			redditArgs.Sort = types.RedditSortNew
			err := redditArgs.Validate()
			Expect(err).To(MatchError(search.ErrNoQueries))
		})

		It("should fail if urls are not provided for scrapeurls", func() {
			redditArgs := search.NewScrapeUrlsArguments()
			redditArgs.Sort = types.RedditSortNew
			err := redditArgs.Validate()
			Expect(err).To(MatchError(search.ErrNoUrls))
		})

		It("should fail if queries are provided for scrapeurls", func() {
			redditArgs := search.NewScrapeUrlsArguments()
			redditArgs.Queries = []string{"test"}
			redditArgs.URLs = []string{"https://www.reddit.com/r/golang/comments/foo/bar/"}
			redditArgs.Sort = types.RedditSortNew
			err := redditArgs.Validate()
			Expect(err).To(MatchError(search.ErrQueriesNotAllowed))
		})

		It("should fail if urls are provided for searchposts", func() {
			redditArgs := search.NewSearchPostsArguments()
			redditArgs.Queries = []string{"test"}
			redditArgs.URLs = []string{"https://www.reddit.com/r/golang/comments/foo/bar"}
			redditArgs.Sort = types.RedditSortNew
			err := redditArgs.Validate()
			Expect(err).To(MatchError(search.ErrUrlsNotAllowed))
		})

		It("should fail with an invalid URL", func() {
			redditArgs := reddit.NewScrapeUrlsArguments()
			redditArgs.URLs = []string{"ht tp://invalid-url.com"}
			redditArgs.Sort = types.RedditSortNew
			err := redditArgs.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("is not a valid URL"))
		})

		It("should fail with an invalid domain", func() {
			redditArgs := search.NewScrapeUrlsArguments()
			redditArgs.URLs = []string{"https://www.google.com"}
			redditArgs.Sort = types.RedditSortNew
			err := redditArgs.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid Reddit URL"))
		})

		It("should fail if the URL is not a post or comment", func() {
			redditArgs := search.NewScrapeUrlsArguments()
			redditArgs.URLs = []string{"https://www.reddit.com/r/golang/"}
			redditArgs.Sort = types.RedditSortNew
			err := redditArgs.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not a Reddit post or comment URL"))
		})
	})
})
