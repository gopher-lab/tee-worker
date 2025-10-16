package trending_test

import (
	"encoding/json"
	"errors"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/api/args/tiktok/trending"
	"github.com/masa-finance/tee-worker/api/types"
)

var _ = Describe("TikTokTrendingArguments", func() {
	Describe("Marshalling and unmarshalling", func() {
		It("should unmarshal valid arguments with all fields", func() {
			var args trending.Arguments
			jsonData := []byte(`{"type":"searchbytrending","country_code":"US","sort_by":"vv","max_items":50,"period":"7"}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.Type).To(Equal(types.CapSearchByTrending))
			Expect(args.CountryCode).To(Equal("US"))
			Expect(args.SortBy).To(Equal("vv"))
			Expect(args.MaxItems).To(Equal(50))
			Expect(args.Period).To(Equal("7"))
		})

		It("should unmarshal valid arguments with minimal fields", func() {
			var args trending.Arguments
			jsonData := []byte(`{"type":"searchbytrending"}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.Type).To(Equal(types.CapSearchByTrending))
			Expect(args.CountryCode).To(Equal("US")) // Default
			Expect(args.SortBy).To(Equal("vv"))      // Default
			Expect(args.Period).To(Equal("7"))       // Default
			Expect(args.MaxItems).To(Equal(0))       // No default for MaxItems
		})

		It("should fail unmarshal with invalid JSON", func() {
			var args trending.Arguments
			jsonData := []byte(`{"type":"searchbytrending","country_code":"US"`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).To(HaveOccurred())
		})

		It("should fail unmarshal with invalid country code", func() {
			var args trending.Arguments
			jsonData := []byte(`{"type":"searchbytrending","country_code":"INVALID"}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "country_code is required")).To(BeTrue())
		})

		It("should fail unmarshal with invalid sort_by", func() {
			var args trending.Arguments
			jsonData := []byte(`{"type":"searchbytrending","sort_by":"invalid"}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "sort_by is required")).To(BeTrue())
		})

		It("should fail unmarshal with invalid period", func() {
			var args trending.Arguments
			jsonData := []byte(`{"type":"searchbytrending","period":"invalid"}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "period is required")).To(BeTrue())
		})

		It("should fail unmarshal with negative max_items", func() {
			var args trending.Arguments
			jsonData := []byte(`{"type":"searchbytrending","max_items":-1}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "max_items must be non-negative")).To(BeTrue())
		})
	})

	Describe("Validation", func() {
		It("should succeed with valid arguments", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				SortBy:      "vv",
				MaxItems:    50,
				Period:      "7",
			}
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with invalid country code", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "INVALID",
				SortBy:      "vv",
				MaxItems:    50,
				Period:      "7",
			}
			err := args.Validate()
			Expect(errors.Is(err, trending.ErrTrendingCountryCodeRequired)).To(BeTrue())
		})

		It("should fail with invalid sort_by", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				SortBy:      "invalid",
				MaxItems:    50,
				Period:      "7",
			}
			err := args.Validate()
			Expect(errors.Is(err, trending.ErrTrendingSortByRequired)).To(BeTrue())
		})

		It("should fail with invalid period", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				SortBy:      "vv",
				MaxItems:    50,
				Period:      "invalid",
			}
			err := args.Validate()
			Expect(errors.Is(err, trending.ErrTrendingPeriodRequired)).To(BeTrue())
		})

		It("should fail with negative max_items", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				SortBy:      "vv",
				MaxItems:    -1,
				Period:      "7",
			}
			err := args.Validate()
			Expect(errors.Is(err, trending.ErrTrendingMaxItemsNegative)).To(BeTrue())
		})
	})

	Describe("Default values", func() {
		It("should set default country_code when not provided", func() {
			args := &trending.Arguments{
				Type:   types.CapSearchByTrending,
				SortBy: "vv",
				Period: "7",
			}
			args.SetDefaultValues()
			Expect(args.CountryCode).To(Equal("US"))
		})

		It("should set default sort_by when not provided", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				Period:      "7",
			}
			args.SetDefaultValues()
			Expect(args.SortBy).To(Equal("vv"))
		})

		It("should set default period when not provided", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				SortBy:      "vv",
			}
			args.SetDefaultValues()
			Expect(args.Period).To(Equal("7"))
		})

		It("should not override existing values", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "CA",
				SortBy:      "like",
				Period:      "30",
			}
			args.SetDefaultValues()
			Expect(args.CountryCode).To(Equal("CA"))
			Expect(args.SortBy).To(Equal("like"))
			Expect(args.Period).To(Equal("30"))
		})
	})

	Describe("Country code validation", func() {
		It("should accept valid country codes", func() {
			validCountries := []string{"US", "CA", "GB", "AU", "DE", "FR", "JP", "KR", "BR"}
			for _, country := range validCountries {
				args := &trending.Arguments{
					Type:        types.CapSearchByTrending,
					CountryCode: country,
					SortBy:      "vv",
					Period:      "7",
				}
				err := args.Validate()
				Expect(err).ToNot(HaveOccurred(), "Country %s should be valid", country)
			}
		})

		It("should accept lowercase country codes", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "us",
				SortBy:      "vv",
				Period:      "7",
			}
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should reject invalid country codes", func() {
			invalidCountries := []string{"INVALID", "XX", "123", ""}
			for _, country := range invalidCountries {
				args := &trending.Arguments{
					Type:        types.CapSearchByTrending,
					CountryCode: country,
					SortBy:      "vv",
					Period:      "7",
				}
				err := args.Validate()
				Expect(err).To(HaveOccurred(), "Country %s should be invalid", country)
			}
		})
	})

	Describe("Sort by validation", func() {
		It("should accept valid sort options", func() {
			validSorts := []string{"vv", "like", "comment", "repost"}
			for _, sort := range validSorts {
				args := &trending.Arguments{
					Type:        types.CapSearchByTrending,
					CountryCode: "US",
					SortBy:      sort,
					Period:      "7",
				}
				err := args.Validate()
				Expect(err).ToNot(HaveOccurred(), "Sort %s should be valid", sort)
			}
		})

		It("should accept uppercase sort options", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				SortBy:      "LIKE",
				Period:      "7",
			}
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should reject invalid sort options", func() {
			invalidSorts := []string{"invalid", "views", "likes", ""}
			for _, sort := range invalidSorts {
				args := &trending.Arguments{
					Type:        types.CapSearchByTrending,
					CountryCode: "US",
					SortBy:      sort,
					Period:      "7",
				}
				err := args.Validate()
				Expect(err).To(HaveOccurred(), "Sort %s should be invalid", sort)
			}
		})
	})

	Describe("Period validation", func() {
		It("should accept valid periods", func() {
			validPeriods := []string{"7", "30"}
			for _, period := range validPeriods {
				args := &trending.Arguments{
					Type:        types.CapSearchByTrending,
					CountryCode: "US",
					SortBy:      "vv",
					Period:      period,
				}
				err := args.Validate()
				Expect(err).ToNot(HaveOccurred(), "Period %s should be valid", period)
			}
		})

		It("should reject invalid periods", func() {
			invalidPeriods := []string{"1", "14", "60", "invalid", ""}
			for _, period := range invalidPeriods {
				args := &trending.Arguments{
					Type:        types.CapSearchByTrending,
					CountryCode: "US",
					SortBy:      "vv",
					Period:      period,
				}
				err := args.Validate()
				Expect(err).To(HaveOccurred(), "Period %s should be invalid", period)
			}
		})
	})

	Describe("MaxItems validation", func() {
		It("should accept zero max_items", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				SortBy:      "vv",
				MaxItems:    0,
				Period:      "7",
			}
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept positive max_items", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				SortBy:      "vv",
				MaxItems:    100,
				Period:      "7",
			}
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should reject negative max_items", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				SortBy:      "vv",
				MaxItems:    -1,
				Period:      "7",
			}
			err := args.Validate()
			Expect(errors.Is(err, trending.ErrTrendingMaxItemsNegative)).To(BeTrue())
		})
	})

	Describe("Job capability", func() {
		It("should return the searchbytrending capability", func() {
			args := trending.NewArguments()
			Expect(args.GetCapability()).To(Equal(types.CapSearchByTrending))
		})

		It("should validate capability for TiktokJob", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				SortBy:      "vv",
				MaxItems:    50,
				Period:      "7",
			}
			err := args.ValidateCapability(types.TiktokJob)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail validation for incompatible job type", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				SortBy:      "vv",
				MaxItems:    50,
				Period:      "7",
			}
			err := args.ValidateCapability(types.TwitterJob)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Edge cases", func() {
		It("should handle mixed case country codes", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "us",
				SortBy:      "vv",
				Period:      "7",
			}
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle mixed case sort options", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				SortBy:      "LIKE",
				Period:      "7",
			}
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle large max_items values", func() {
			args := &trending.Arguments{
				Type:        types.CapSearchByTrending,
				CountryCode: "US",
				SortBy:      "vv",
				MaxItems:    10000,
				Period:      "7",
			}
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle all supported countries", func() {
			supportedCountries := []string{
				"AU", "BR", "CA", "EG", "FR", "DE", "ID", "IL", "IT", "JP",
				"MY", "PH", "RU", "SA", "SG", "KR", "ES", "TW", "TH", "TR",
				"AE", "GB", "US", "VN",
			}
			for _, country := range supportedCountries {
				args := &trending.Arguments{
					Type:        types.CapSearchByTrending,
					CountryCode: country,
					SortBy:      "vv",
					Period:      "7",
				}
				err := args.Validate()
				Expect(err).ToNot(HaveOccurred(), "Country %s should be supported", country)
			}
		})
	})
})
