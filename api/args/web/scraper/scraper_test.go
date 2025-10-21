package scraper_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/api/args/web"
	"github.com/masa-finance/tee-worker/api/args/web/scraper"
)

var _ = Describe("WebArguments", func() {
	Describe("Marshalling and unmarshalling", func() {
		It("should set default values", func() {
			webArgs := web.NewScraperArguments()
			webArgs.URL = "https://example.com"
			webArgs.MaxDepth = 0
			webArgs.MaxPages = 0
			jsonData, err := json.Marshal(webArgs)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &webArgs)
			Expect(err).ToNot(HaveOccurred())
			Expect(webArgs.MaxPages).To(Equal(1))
		})

		It("should override default values", func() {
			webArgs := web.NewScraperArguments()
			webArgs.URL = "https://example.com"
			webArgs.MaxDepth = 2
			webArgs.MaxPages = 5
			jsonData, err := json.Marshal(webArgs)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &webArgs)
			Expect(err).ToNot(HaveOccurred())
			Expect(webArgs.MaxPages).To(Equal(5))
		})

		It("should fail unmarshal when url is missing", func() {
			var webArgs web.ScraperArguments
			jsonData := []byte(`{"type":"scraper","max_depth":1,"max_pages":1}`)
			err := json.Unmarshal(jsonData, &webArgs)
			Expect(errors.Is(err, scraper.ErrURLRequired)).To(BeTrue())
		})
	})

	Describe("Validation", func() {
		It("should succeed with valid arguments", func() {
			webArgs := web.NewScraperArguments()
			webArgs.URL = "https://example.com"
			webArgs.MaxDepth = 2
			webArgs.MaxPages = 3
			err := webArgs.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when url is missing", func() {
			webArgs := web.NewScraperArguments()
			webArgs.MaxDepth = 0
			webArgs.MaxPages = 1
			err := webArgs.Validate()
			Expect(errors.Is(err, scraper.ErrURLRequired)).To(BeTrue())
		})

		It("should fail with an invalid URL format", func() {
			webArgs := web.NewScraperArguments()
			webArgs.URL = "http:// invalid.com"
			webArgs.MaxDepth = 0
			webArgs.MaxPages = 1
			err := webArgs.Validate()
			Expect(errors.Is(err, scraper.ErrURLInvalid)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("invalid URL format"))
		})

		It("should fail when scheme is missing", func() {
			webArgs := web.NewScraperArguments()
			webArgs.URL = "example.com"
			webArgs.MaxDepth = 0
			webArgs.MaxPages = 1
			err := webArgs.Validate()
			Expect(errors.Is(err, scraper.ErrURLSchemeMissing)).To(BeTrue())
		})

		It("should fail when max depth is negative", func() {
			webArgs := web.NewScraperArguments()
			webArgs.URL = "https://example.com"
			webArgs.MaxDepth = -1
			webArgs.MaxPages = 1
			err := webArgs.Validate()
			Expect(errors.Is(err, scraper.ErrMaxDepth)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("got -1"))
		})

		It("should fail when max pages is less than 1", func() {
			webArgs := web.NewScraperArguments()
			webArgs.URL = "https://example.com"
			webArgs.MaxDepth = 0
			webArgs.MaxPages = 0
			err := webArgs.Validate()
			Expect(errors.Is(err, scraper.ErrMaxPages)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("got 0"))
		})
	})

	Describe("ToWebScraperRequest", func() {
		It("should map fields correctly", func() {
			webArgs := web.NewScraperArguments()
			webArgs.URL = "https://example.com"
			webArgs.MaxDepth = 2
			webArgs.MaxPages = 3
			req := webArgs.ToScraperRequest()
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
