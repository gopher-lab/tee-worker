package args_test

import (
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/api/args"
	"github.com/masa-finance/tee-worker/api/types"
)

var _ = Describe("RedditArguments", func() {
	Describe("Marshalling and unmarshalling", func() {
		It("should set default values", func() {
			redditArgs := args.RedditArguments{
				QueryType: types.RedditSearchPosts,
				Queries:   []string{"Zaphod", "Ford"},
			}
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
			redditArgs := args.RedditArguments{
				QueryType:      types.RedditSearchPosts,
				Queries:        []string{"Zaphod", "Ford"},
				MaxItems:       20,
				MaxPosts:       21,
				MaxComments:    22,
				MaxCommunities: 23,
				MaxUsers:       24,
				Sort:           types.RedditSortTop,
			}
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
			redditArgs := &args.RedditArguments{
				QueryType: types.RedditSearchPosts,
				Queries:   []string{"test"},
				Sort:      types.RedditSortNew,
			}
			err := redditArgs.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should succeed with valid scrapeurls arguments", func() {
			redditArgs := &args.RedditArguments{
				QueryType: types.RedditScrapeUrls,
				URLs:      []string{"https://www.reddit.com/r/golang/comments/foo/bar"},
				Sort:      types.RedditSortNew,
			}
			err := redditArgs.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with an invalid type", func() {
			redditArgs := &args.RedditArguments{
				QueryType: "invalidtype",
				Queries:   []string{"test"},
				Sort:      types.RedditSortNew,
			}
			err := redditArgs.Validate()
			Expect(err).To(MatchError(args.ErrRedditInvalidType))
		})

		It("should fail with an invalid sort", func() {
			redditArgs := &args.RedditArguments{
				QueryType: types.RedditSearchPosts,
				Queries:   []string{"test"},
				Sort:      "invalidsort",
			}
			err := redditArgs.Validate()
			Expect(err).To(MatchError(args.ErrRedditInvalidSort))
		})

		It("should fail if the after time is in the future", func() {
			redditArgs := &args.RedditArguments{
				QueryType: types.RedditSearchPosts,
				Queries:   []string{"test"},
				Sort:      types.RedditSortNew,
				After:     time.Now().Add(24 * time.Hour),
			}
			err := redditArgs.Validate()
			Expect(err).To(MatchError(args.ErrRedditTimeInTheFuture))
		})

		It("should fail if queries are not provided for searchposts", func() {
			redditArgs := &args.RedditArguments{
				QueryType: types.RedditSearchPosts,
				Sort:      types.RedditSortNew,
			}
			err := redditArgs.Validate()
			Expect(err).To(MatchError(args.ErrRedditNoQueries))
		})

		It("should fail if urls are not provided for scrapeurls", func() {
			redditArgs := &args.RedditArguments{
				QueryType: types.RedditScrapeUrls,
				Sort:      types.RedditSortNew,
			}
			err := redditArgs.Validate()
			Expect(err).To(MatchError(args.ErrRedditNoUrls))
		})

		It("should fail if queries are provided for scrapeurls", func() {
			redditArgs := &args.RedditArguments{
				QueryType: types.RedditScrapeUrls,
				Queries:   []string{"test"},
				URLs:      []string{"https://www.reddit.com/r/golang/comments/foo/bar/"},
				Sort:      types.RedditSortNew,
			}
			err := redditArgs.Validate()
			Expect(err).To(MatchError(args.ErrRedditQueriesNotAllowed))
		})

		It("should fail if urls are provided for searchposts", func() {
			redditArgs := &args.RedditArguments{
				QueryType: types.RedditSearchPosts,
				Queries:   []string{"test"},
				URLs:      []string{"https://www.reddit.com/r/golang/comments/foo/bar"},
				Sort:      types.RedditSortNew,
			}
			err := redditArgs.Validate()
			Expect(err).To(MatchError(args.ErrRedditUrlsNotAllowed))
		})

		It("should fail with an invalid URL", func() {
			redditArgs := &args.RedditArguments{
				QueryType: types.RedditScrapeUrls,
				URLs:      []string{"ht tp://invalid-url.com"},
				Sort:      types.RedditSortNew,
			}
			err := redditArgs.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("is not a valid URL"))
		})

		It("should fail with an invalid domain", func() {
			redditArgs := &args.RedditArguments{
				QueryType: types.RedditScrapeUrls,
				URLs:      []string{"https://www.google.com"},
				Sort:      types.RedditSortNew,
			}
			err := redditArgs.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid Reddit URL"))
		})

		It("should fail if the URL is not a post or comment", func() {
			redditArgs := &args.RedditArguments{
				QueryType: types.RedditScrapeUrls,
				URLs:      []string{"https://www.reddit.com/r/golang/"},
				Sort:      types.RedditSortNew,
			}
			err := redditArgs.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not a Reddit post or comment URL"))
		})
	})
})
