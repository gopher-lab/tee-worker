package page_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/api/args/web/page"
	"github.com/masa-finance/tee-worker/api/types"
)

var _ = Describe("WebArguments", func() {
	Describe("Marshalling and unmarshalling", func() {
		It("should set default values", func() {
			webArgs := page.NewArguments()
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
			webArgs := page.NewArguments()
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
			var webArgs page.Arguments
			jsonData := []byte(`{"type":"scraper","max_depth":1,"max_pages":1}`)
			err := json.Unmarshal(jsonData, &webArgs)
			Expect(errors.Is(err, page.ErrURLRequired)).To(BeTrue())
		})
	})

	Describe("Validation", func() {
		It("should succeed with valid arguments", func() {
			webArgs := page.NewArguments()
			webArgs.URL = "https://example.com"
			webArgs.MaxDepth = 2
			webArgs.MaxPages = 3
			err := webArgs.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when url is missing", func() {
			webArgs := page.NewArguments()
			webArgs.MaxDepth = 0
			webArgs.MaxPages = 1
			err := webArgs.Validate()
			Expect(errors.Is(err, page.ErrURLRequired)).To(BeTrue())
		})

		It("should fail with an invalid URL format", func() {
			webArgs := page.NewArguments()
			webArgs.URL = "http:// invalid.com"
			webArgs.MaxDepth = 0
			webArgs.MaxPages = 1
			err := webArgs.Validate()
			Expect(errors.Is(err, page.ErrURLInvalid)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("invalid URL format"))
		})

		It("should fail when scheme is missing", func() {
			webArgs := page.NewArguments()
			webArgs.URL = "example.com"
			webArgs.MaxDepth = 0
			webArgs.MaxPages = 1
			err := webArgs.Validate()
			Expect(errors.Is(err, page.ErrURLSchemeMissing)).To(BeTrue())
		})

		It("should fail when max depth is negative", func() {
			webArgs := page.NewArguments()
			webArgs.URL = "https://example.com"
			webArgs.MaxDepth = -1
			webArgs.MaxPages = 1
			err := webArgs.Validate()
			Expect(errors.Is(err, page.ErrMaxDepth)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("got -1"))
		})

		It("should fail when max pages is less than 1", func() {
			webArgs := page.NewArguments()
			webArgs.URL = "https://example.com"
			webArgs.MaxDepth = 0
			webArgs.MaxPages = 0
			err := webArgs.Validate()
			Expect(errors.Is(err, page.ErrMaxPages)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("got 0"))
		})
	})

	Describe("Job capability", func() {
		It("should return the scraper capability", func() {
			webArgs := page.NewArguments()
			Expect(webArgs.GetCapability()).To(Equal(types.CapScraper))
		})

		It("should validate capability for WebJob", func() {
			webArgs := page.NewArguments()
			webArgs.URL = "https://example.com"
			webArgs.MaxDepth = 1
			webArgs.MaxPages = 1
			err := webArgs.ValidateCapability(types.WebJob)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("ToWebScraperRequest", func() {
		It("should map fields correctly", func() {
			webArgs := page.NewArguments()
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
