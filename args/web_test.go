package args_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-types/args"
	"github.com/masa-finance/tee-types/types"
)

var _ = Describe("WebArguments", func() {
	Describe("Marshalling and unmarshalling", func() {
		It("should set default values", func() {
			webArgs := args.WebArguments{
				QueryType: types.WebScraper,
				URL:       "https://example.com",
				MaxDepth:  0,
				MaxPages:  0,
			}
			jsonData, err := json.Marshal(webArgs)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &webArgs)
			Expect(err).ToNot(HaveOccurred())
			Expect(webArgs.MaxPages).To(Equal(1))
		})

		It("should override default values", func() {
			webArgs := args.WebArguments{
				QueryType: types.WebScraper,
				URL:       "https://example.com",
				MaxDepth:  2,
				MaxPages:  5,
			}
			jsonData, err := json.Marshal(webArgs)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &webArgs)
			Expect(err).ToNot(HaveOccurred())
			Expect(webArgs.MaxPages).To(Equal(5))
		})

		It("should fail unmarshal when url is missing", func() {
			var webArgs args.WebArguments
			jsonData := []byte(`{"type":"scraper","max_depth":1,"max_pages":1}`)
			err := json.Unmarshal(jsonData, &webArgs)
			Expect(errors.Is(err, args.ErrWebURLRequired)).To(BeTrue())
		})
	})

	Describe("Validation", func() {
		It("should succeed with valid arguments", func() {
			webArgs := &args.WebArguments{
				QueryType: types.WebScraper,
				URL:       "https://example.com",
				MaxDepth:  2,
				MaxPages:  3,
			}
			err := webArgs.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when url is missing", func() {
			webArgs := &args.WebArguments{
				QueryType: types.WebScraper,
				MaxDepth:  0,
				MaxPages:  1,
			}
			err := webArgs.Validate()
			Expect(errors.Is(err, args.ErrWebURLRequired)).To(BeTrue())
		})

		It("should fail with an invalid URL format", func() {
			webArgs := &args.WebArguments{
				QueryType: types.WebScraper,
				URL:       "http:// invalid.com",
				MaxDepth:  0,
				MaxPages:  1,
			}
			err := webArgs.Validate()
			Expect(errors.Is(err, args.ErrWebURLInvalid)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("invalid URL format"))
		})

		It("should fail when scheme is missing", func() {
			webArgs := &args.WebArguments{
				QueryType: types.WebScraper,
				URL:       "example.com",
				MaxDepth:  0,
				MaxPages:  1,
			}
			err := webArgs.Validate()
			Expect(errors.Is(err, args.ErrWebURLSchemeMissing)).To(BeTrue())
		})

		It("should fail when max depth is negative", func() {
			webArgs := &args.WebArguments{
				QueryType: types.WebScraper,
				URL:       "https://example.com",
				MaxDepth:  -1,
				MaxPages:  1,
			}
			err := webArgs.Validate()
			Expect(errors.Is(err, args.ErrWebMaxDepth)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("got -1"))
		})

		It("should fail when max pages is less than 1", func() {
			webArgs := &args.WebArguments{
				QueryType: types.WebScraper,
				URL:       "https://example.com",
				MaxDepth:  0,
				MaxPages:  0,
			}
			err := webArgs.Validate()
			Expect(errors.Is(err, args.ErrWebMaxPages)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("got 0"))
		})
	})

	Describe("Job capability", func() {
		It("should return the scraper capability", func() {
			webArgs := &args.WebArguments{}
			Expect(webArgs.GetCapability()).To(Equal(types.CapScraper))
		})

		It("should validate capability for WebJob", func() {
			webArgs := &args.WebArguments{
				QueryType: types.WebScraper,
				URL:       "https://example.com",
				MaxDepth:  1,
				MaxPages:  1,
			}
			err := webArgs.ValidateForJobType(types.WebJob)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("ToWebScraperRequest", func() {
		It("should map fields correctly", func() {
			webArgs := args.WebArguments{
				QueryType: types.WebScraper,
				URL:       "https://example.com",
				MaxDepth:  2,
				MaxPages:  3,
			}
			req := webArgs.ToWebScraperRequest()
			Expect(req.StartUrls).To(HaveLen(1))
			Expect(req.StartUrls[0].URL).To(Equal("https://example.com"))
			Expect(req.StartUrls[0].Method).To(Equal("GET"))
			Expect(req.MaxCrawlDepth).To(Equal(2))
			Expect(req.MaxCrawlPages).To(Equal(3))
			Expect(req.RespectRobotsTxtFile).To(BeFalse())
			Expect(req.SaveMarkdown).To(BeTrue())
		})
	})
})
