package query_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/api/args/tiktok/query"
)

var _ = Describe("TikTokQueryArguments", func() {
	Describe("Marshalling and unmarshalling", func() {
		It("should unmarshal valid arguments with search", func() {
			var args query.Arguments
			jsonData := []byte(`{"type":"searchbyquery","search":["test query","another query"],"max_items":20}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.Search).To(Equal([]string{"test query", "another query"}))
			Expect(args.MaxItems).To(Equal(uint(20)))
		})

		It("should unmarshal valid arguments with start_urls", func() {
			var args query.Arguments
			jsonData := []byte(`{"type":"searchbyquery","start_urls":["https://tiktok.com/@user1","https://tiktok.com/@user2"],"max_items":15}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.StartUrls).To(Equal([]string{"https://tiktok.com/@user1", "https://tiktok.com/@user2"}))
			Expect(args.MaxItems).To(Equal(uint(15)))
		})

		It("should unmarshal valid arguments with both search and start_urls", func() {
			var args query.Arguments
			jsonData := []byte(`{"type":"searchbyquery","search":["test"],"start_urls":["https://tiktok.com/@user"],"max_items":5}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.Search).To(Equal([]string{"test"}))
			Expect(args.StartUrls).To(Equal([]string{"https://tiktok.com/@user"}))
			Expect(args.MaxItems).To(Equal(uint(5)))
		})

		It("should unmarshal valid arguments without max_items (should use default)", func() {
			var args query.Arguments
			jsonData := []byte(`{"type":"searchbyquery","search":["test query"]}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.Search).To(Equal([]string{"test query"}))
			Expect(args.MaxItems).To(Equal(uint(10))) // Default value
		})

		It("should fail unmarshal with invalid JSON", func() {
			var args query.Arguments
			jsonData := []byte(`{"type":"searchbyquery","search":["test query"`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).To(HaveOccurred())
		})

		It("should fail unmarshal when neither search nor start_urls are provided", func() {
			var args query.Arguments
			jsonData := []byte(`{"type":"searchbyquery","max_items":10}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(errors.Is(err, query.ErrSearchOrUrlsRequired)).To(BeTrue())
		})
	})

	Describe("Validation", func() {
		It("should succeed with valid search arguments", func() {
			args := query.NewArguments()
			args.Search = []string{"test query", "another query"}
			args.MaxItems = 20
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should succeed with valid start_urls arguments", func() {
			args := query.NewArguments()
			args.StartUrls = []string{"https://tiktok.com/@user1", "https://tiktok.com/@user2"}
			args.MaxItems = 15
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should succeed with both search and start_urls", func() {
			args := query.NewArguments()
			args.Search = []string{"test"}
			args.StartUrls = []string{"https://tiktok.com/@user"}
			args.MaxItems = 5
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when both search and start_urls are empty", func() {
			args := query.NewArguments()
			args.MaxItems = 10
			err := args.Validate()
			Expect(errors.Is(err, query.ErrSearchOrUrlsRequired)).To(BeTrue())
		})

		It("should fail when search is empty slice", func() {
			args := query.NewArguments()
			args.Search = []string{}
			args.MaxItems = 10
			err := args.Validate()
			Expect(errors.Is(err, query.ErrSearchOrUrlsRequired)).To(BeTrue())
		})

		It("should fail when start_urls is empty slice", func() {
			args := query.NewArguments()
			args.StartUrls = []string{}
			args.MaxItems = 10
			err := args.Validate()
			Expect(errors.Is(err, query.ErrSearchOrUrlsRequired)).To(BeTrue())
		})
	})

	Describe("Default values", func() {
		It("should set default max_items when not provided", func() {
			args := query.NewArguments()
			args.Search = []string{"test"}
			args.SetDefaultValues()
			Expect(args.MaxItems).To(Equal(uint(10)))
		})

		It("should not override existing max_items", func() {
			args := query.NewArguments()
			args.Search = []string{"test"}
			args.MaxItems = 25
			args.SetDefaultValues()
			Expect(args.MaxItems).To(Equal(uint(25)))
		})

		It("should not override zero max_items if explicitly set", func() {
			args := query.NewArguments()
			args.Search = []string{"test"}
			args.MaxItems = 0
			args.SetDefaultValues()
			Expect(args.MaxItems).To(Equal(uint(10))) // Should set default
		})
	})

	Describe("Edge cases", func() {
		It("should handle empty search strings", func() {
			args := query.NewArguments()
			args.Search = []string{"", "valid query", ""}
			args.MaxItems = 10
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle empty start_urls strings", func() {
			args := query.NewArguments()
			args.StartUrls = []string{"", "https://tiktok.com/@user", ""}
			args.MaxItems = 10
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle large max_items values", func() {
			args := query.NewArguments()
			args.Search = []string{"test"}
			args.MaxItems = 1000
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle end_page field", func() {
			args := query.NewArguments()
			args.Search = []string{"test"}
			args.MaxItems = 10
			args.EndPage = 5
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
			Expect(args.EndPage).To(Equal(uint(5)))
		})
	})
})
