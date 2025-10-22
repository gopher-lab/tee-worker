package search_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/api/args/twitter/search"
	"github.com/masa-finance/tee-worker/api/types"
)

var _ = Describe("TwitterSearchArguments", func() {
	Describe("Marshalling and unmarshalling", func() {
		It("should unmarshal valid arguments with all fields", func() {
			var args search.Arguments
			jsonData := []byte(`{
				"type": "searchbyquery",
				"query": "test query",
				"count": 50,
				"start_time": "2023-01-01T00:00:00Z",
				"end_time": "2023-12-31T23:59:59Z",
				"max_results": 100,
				"next_cursor": "cursor123"
			}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.Query).To(Equal("test query"))
			Expect(args.Count).To(Equal(50))
			Expect(args.StartTime).To(Equal("2023-01-01T00:00:00Z"))
			Expect(args.EndTime).To(Equal("2023-12-31T23:59:59Z"))
			Expect(args.MaxResults).To(Equal(100))
			Expect(args.NextCursor).To(Equal("cursor123"))
		})

		It("should unmarshal valid arguments with minimal fields", func() {
			var args search.Arguments
			jsonData := []byte(`{
				"type": "searchbyquery",
				"query": "minimal test"
			}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.Query).To(Equal("minimal test"))
			Expect(args.Count).To(Equal(0))
			Expect(args.MaxResults).To(Equal(10)) // SetDefaultValues() sets this to MaxResults
		})

		It("should fail unmarshal with invalid JSON", func() {
			var args search.Arguments
			jsonData := []byte(`{"type":"searchbyquery","query":"test"`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).To(HaveOccurred())
			// The error is a JSON syntax error, not wrapped with ErrUnmarshalling
			// since the JSON is malformed before reaching the custom UnmarshalJSON method
		})

		It("should set default values after unmarshalling", func() {
			var args search.Arguments
			jsonData := []byte(`{"type":"searchbyquery","query":"test"}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).ToNot(HaveOccurred())
			// Default values should be set by SetDefaultValues()
			Expect(args.GetCapability()).To(Equal(types.CapSearchByQuery))
		})
	})

	Describe("Validation", func() {
		It("should succeed with valid arguments", func() {
			args := search.NewArguments()
			args.Query = "test query"
			args.Count = 50
			args.MaxResults = 100
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when count is negative", func() {
			args := search.NewArguments()
			args.Query = "test query"
			args.Count = -1
			err := args.Validate()
			Expect(errors.Is(err, search.ErrCountNegative)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("got: -1"))
		})

		It("should fail when count exceeds maximum", func() {
			args := search.NewArguments()
			args.Query = "test query"
			args.Count = 1001
			err := args.Validate()
			Expect(errors.Is(err, search.ErrCountTooLarge)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("got: 1001"))
		})

		It("should fail when max_results is negative", func() {
			args := search.NewArguments()
			args.Query = "test query"
			args.MaxResults = -1
			err := args.Validate()
			Expect(errors.Is(err, search.ErrMaxResultsNegative)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("got: -1"))
		})

		It("should fail when max_results exceeds maximum", func() {
			args := search.NewArguments()
			args.Query = "test query"
			args.MaxResults = 1001
			err := args.Validate()
			Expect(errors.Is(err, search.ErrMaxResultsTooLarge)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("got: 1001"))
		})

		It("should succeed with count at maximum boundary", func() {
			args := search.NewArguments()
			args.Query = "test query"
			args.Count = 1000
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should succeed with max_results at maximum boundary", func() {
			args := search.NewArguments()
			args.Query = "test query"
			args.MaxResults = 1000
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Operation Type Detection", func() {
		Context("Single Tweet Operations", func() {
			It("should identify getbyid as single tweet operation", func() {
				args := search.NewArguments()
				args.Type = types.CapGetById
				Expect(args.IsSingleTweetOperation()).To(BeTrue())
			})

			It("should not identify searchbyquery as single tweet operation", func() {
				args := search.NewArguments()
				// Type is already CapSearchByQuery from NewArguments()
				Expect(args.IsSingleTweetOperation()).To(BeFalse())
			})
		})

		Context("Multiple Tweet Operations", func() {
			It("should identify searchbyquery as multiple tweet operation", func() {
				args := search.NewArguments()
				// Type is already CapSearchByQuery from NewArguments()
				Expect(args.IsMultipleTweetOperation()).To(BeTrue())
			})

			It("should identify searchbyfullarchive as multiple tweet operation", func() {
				args := search.NewArguments()
				args.Type = types.CapSearchByFullArchive
				Expect(args.IsMultipleTweetOperation()).To(BeTrue())
			})

			It("should identify gettweets as multiple tweet operation", func() {
				args := search.NewArguments()
				args.Type = types.CapGetTweets
				Expect(args.IsMultipleTweetOperation()).To(BeTrue())
			})

			It("should identify getreplies as multiple tweet operation", func() {
				args := search.NewArguments()
				args.Type = types.CapGetReplies
				Expect(args.IsMultipleTweetOperation()).To(BeTrue())
			})

			It("should identify getmedia as multiple tweet operation", func() {
				args := search.NewArguments()
				args.Type = types.CapGetMedia
				Expect(args.IsMultipleTweetOperation()).To(BeTrue())
			})

			It("should not identify getbyid as multiple tweet operation", func() {
				args := search.NewArguments()
				args.Type = types.CapGetById
				Expect(args.IsMultipleTweetOperation()).To(BeFalse())
			})
		})

		Context("Single Profile Operations", func() {
			It("should identify getprofilebyid as single profile operation", func() {
				args := search.NewArguments()
				args.Type = types.CapGetProfileById
				Expect(args.IsSingleProfileOperation()).To(BeTrue())
			})

			It("should identify searchbyprofile as single profile operation", func() {
				args := search.NewArguments()
				args.Type = types.CapSearchByProfile
				Expect(args.IsSingleProfileOperation()).To(BeTrue())
			})

			It("should not identify getfollowers as single profile operation", func() {
				args := search.NewArguments()
				args.Type = types.CapGetFollowers
				Expect(args.IsSingleProfileOperation()).To(BeFalse())
			})
		})

		Context("Multiple Profile Operations", func() {
			It("should identify getretweeters as multiple profile operation", func() {
				args := search.NewArguments()
				args.Type = types.CapGetRetweeters
				Expect(args.IsMultipleProfileOperation()).To(BeTrue())
			})

			It("should not identify getprofilebyid as multiple profile operation", func() {
				args := search.NewArguments()
				args.Type = types.CapGetProfileById
				Expect(args.IsMultipleProfileOperation()).To(BeFalse())
			})
		})

		Context("Single Space Operations", func() {
			It("should identify getspace as single space operation", func() {
				args := search.NewArguments()
				args.Type = types.CapGetSpace
				Expect(args.IsSingleSpaceOperation()).To(BeTrue())
			})

			It("should not identify searchbyquery as single space operation", func() {
				args := search.NewArguments()
				// Type is already CapSearchByQuery from NewArguments()
				Expect(args.IsSingleSpaceOperation()).To(BeFalse())
			})
		})

		Context("Trends Operations", func() {
			It("should identify gettrends as trends operation", func() {
				args := search.NewArguments()
				args.Type = types.CapGetTrends
				Expect(args.IsTrendsOperation()).To(BeTrue())
			})

			It("should not identify searchbyquery as trends operation", func() {
				args := search.NewArguments()
				// Type is already CapSearchByQuery from NewArguments()
				Expect(args.IsTrendsOperation()).To(BeFalse())
			})
		})
	})

	Describe("Constants and Error Values", func() {
		It("should have correct MaxResults constant", func() {
			Expect(search.MaxResults).To(Equal(1000))
		})

		It("should have correct error messages", func() {
			Expect(search.ErrCountNegative.Error()).To(Equal("count must be non-negative"))
			Expect(search.ErrCountTooLarge.Error()).To(Equal("count must be less than or equal to 1000"))
			Expect(search.ErrMaxResultsTooLarge.Error()).To(Equal("max_results must be less than or equal to 1000"))
			Expect(search.ErrMaxResultsNegative.Error()).To(Equal("max_results must be non-negative"))
			Expect(search.ErrUnmarshalling.Error()).To(Equal("failed to unmarshal twitter search arguments"))
		})
	})

	Describe("JSON Marshalling", func() {
		It("should marshal arguments correctly", func() {
			args := search.NewArguments()
			args.Query = "test query"
			args.Count = 50
			args.StartTime = "2023-01-01T00:00:00Z"
			args.EndTime = "2023-12-31T23:59:59Z"
			args.MaxResults = 100
			args.NextCursor = "cursor123"
			jsonData, err := json.Marshal(args)
			Expect(err).ToNot(HaveOccurred())

			var unmarshalled search.Arguments
			err = json.Unmarshal(jsonData, &unmarshalled)
			Expect(err).ToNot(HaveOccurred())
			Expect(unmarshalled.Query).To(Equal(args.Query))
			Expect(unmarshalled.Count).To(Equal(args.Count))
			Expect(unmarshalled.StartTime).To(Equal(args.StartTime))
			Expect(unmarshalled.EndTime).To(Equal(args.EndTime))
			Expect(unmarshalled.MaxResults).To(Equal(args.MaxResults))
			Expect(unmarshalled.NextCursor).To(Equal(args.NextCursor))
		})
	})
})
