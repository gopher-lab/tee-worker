package args_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/api/args"
	"github.com/masa-finance/tee-worker/api/args/reddit"
	"github.com/masa-finance/tee-worker/api/args/telemetry"
	"github.com/masa-finance/tee-worker/api/args/tiktok"
	"github.com/masa-finance/tee-worker/api/args/twitter"
	"github.com/masa-finance/tee-worker/api/args/web"
	"github.com/masa-finance/tee-worker/api/types"
)

var _ = Describe("Unmarshaller", func() {
	Describe("UnmarshalJobArguments", func() {
		Context("with a WebJob", func() {
			It("should unmarshal the arguments correctly", func() {
				argsMap := map[string]any{
					"url":       "https://example.com",
					"max_depth": 2,
				}
				jobArgs, err := args.UnmarshalJobArguments(types.WebJob, argsMap)
				Expect(err).ToNot(HaveOccurred())
				webArgs, ok := jobArgs.(*web.ScraperArguments)
				Expect(ok).To(BeTrue())
				Expect(webArgs.URL).To(Equal("https://example.com"))
				Expect(webArgs.MaxDepth).To(Equal(2))
			})
		})

		Context("with a TiktokJob", func() {
			It("should unmarshal the arguments correctly", func() {
				argsMap := map[string]any{
					"type":      "transcription",
					"video_url": "https://www.tiktok.com/@user/video/123",
					"language":  "en-us",
				}
				jobArgs, err := args.UnmarshalJobArguments(types.TiktokJob, argsMap)
				Expect(err).ToNot(HaveOccurred())
				tiktokArgs, ok := jobArgs.(*tiktok.TranscriptionArguments)
				Expect(ok).To(BeTrue())
				Expect(tiktokArgs.VideoURL).To(Equal("https://www.tiktok.com/@user/video/123"))
				Expect(tiktokArgs.Language).To(Equal("en-us"))
			})
		})

		Context("with a TwitterJob", func() {
			It("should unmarshal the arguments correctly", func() {
				argsMap := map[string]any{
					"type":  "searchbyquery",
					"query": "golang",
					"count": 10,
				}
				jobArgs, err := args.UnmarshalJobArguments(types.TwitterJob, argsMap)
				Expect(err).ToNot(HaveOccurred())
				twitterArgs, ok := jobArgs.(*twitter.SearchArguments)
				Expect(ok).To(BeTrue())
				Expect(twitterArgs.Type).To(Equal(types.CapSearchByQuery))
				Expect(twitterArgs.Query).To(Equal("golang"))
				Expect(twitterArgs.Count).To(Equal(10))
			})
		})

		Context("with a RedditJob", func() {
			It("should unmarshal the arguments correctly", func() {
				argsMap := map[string]any{
					"type":    "searchposts",
					"queries": []string{"golang"},
					"sort":    "new",
				}
				jobArgs, err := args.UnmarshalJobArguments(types.RedditJob, argsMap)
				Expect(err).ToNot(HaveOccurred())
				redditArgs, ok := jobArgs.(*reddit.SearchArguments)
				Expect(ok).To(BeTrue())
				Expect(redditArgs.Type).To(Equal(types.CapSearchPosts))
			})
		})

		Context("with a TelemetryJob", func() {
			It("should return a TelemetryArguments struct", func() {
				argsMap := map[string]any{}
				jobArgs, err := args.UnmarshalJobArguments(types.TelemetryJob, argsMap)
				Expect(err).ToNot(HaveOccurred())
				_, ok := jobArgs.(*telemetry.Arguments)
				Expect(ok).To(BeTrue())
			})
		})

		Context("with an unknown job type", func() {
			It("should return an error", func() {
				argsMap := map[string]any{}
				_, err := args.UnmarshalJobArguments("unknown", argsMap)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown job type"))
			})
		})
	})
})
