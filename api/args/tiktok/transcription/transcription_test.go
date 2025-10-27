package transcription_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/v2/api/args/tiktok/transcription"
)

var _ = Describe("TikTokTranscriptionArguments", func() {
	Describe("Marshalling and unmarshalling", func() {
		It("should unmarshal valid arguments", func() {
			var args transcription.Arguments
			jsonData := []byte(`{"type":"transcription","video_url":"https://tiktok.com/@user/video/123","language":"en-us"}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.VideoURL).To(Equal("https://tiktok.com/@user/video/123"))
			Expect(args.Language).To(Equal("en-us"))
		})

		It("should unmarshal valid arguments without language", func() {
			var args transcription.Arguments
			jsonData := []byte(`{"type":"transcription","video_url":"https://tiktok.com/@user/video/123"}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.VideoURL).To(Equal("https://tiktok.com/@user/video/123"))
			Expect(args.Language).To(Equal("eng-US")) // Default language should be set
		})

		It("should fail unmarshal with invalid JSON", func() {
			var args transcription.Arguments
			jsonData := []byte(`{"type":"transcription","video_url":"https://tiktok.com/@user/video/123"`)
			err := json.Unmarshal(jsonData, &args)
			Expect(err).To(HaveOccurred())
		})

		It("should fail unmarshal when video_url is missing", func() {
			var args transcription.Arguments
			jsonData := []byte(`{"type":"transcription","language":"en-us"}`)
			err := json.Unmarshal(jsonData, &args)
			Expect(errors.Is(err, transcription.ErrVideoURLRequired)).To(BeTrue())
		})
	})

	Describe("Validation", func() {
		It("should succeed with valid arguments", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://tiktok.com/@user/video/123"
			args.Language = "en-us"
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should succeed with valid arguments without language", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://tiktok.com/@user/video/123"
			args.Language = ""
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when video_url is missing", func() {
			args := transcription.NewArguments()
			args.Language = "en-us"
			err := args.Validate()
			Expect(errors.Is(err, transcription.ErrVideoURLRequired)).To(BeTrue())
		})

		It("should fail with an invalid URL format", func() {
			args := transcription.NewArguments()
			args.VideoURL = "not-a-url"
			args.Language = "en-us"
			err := args.Validate()
			Expect(errors.Is(err, transcription.ErrInvalidTikTokURL)).To(BeTrue())
		})

		It("should fail with non-TikTok URL", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://youtube.com/watch?v=123"
			args.Language = "en-us"
			err := args.Validate()
			Expect(errors.Is(err, transcription.ErrInvalidTikTokURL)).To(BeTrue())
		})

		It("should fail with invalid language code format", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://tiktok.com/@user/video/123"
			args.Language = "invalid"
			err := args.Validate()
			Expect(errors.Is(err, transcription.ErrInvalidLanguageCode)).To(BeTrue())
		})
	})

	Describe("TikTok URL validation", func() {
		It("should accept tiktok.com URLs", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://tiktok.com/@user/video/123"
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept www.tiktok.com URLs", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://www.tiktok.com/@user/video/123"
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept m.tiktok.com URLs", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://m.tiktok.com/@user/video/123"
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should reject non-TikTok URLs", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://youtube.com/watch?v=123"
			err := args.Validate()
			Expect(errors.Is(err, transcription.ErrInvalidTikTokURL)).To(BeTrue())
		})
	})

	Describe("Language code validation", func() {
		It("should accept valid 2-letter language codes", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://tiktok.com/@user/video/123"
			args.Language = "en-us"
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept valid 3-letter language codes", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://tiktok.com/@user/video/123"
			args.Language = "eng-us"
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept mixed case language codes", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://tiktok.com/@user/video/123"
			args.Language = "EN-US"
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should reject invalid language format", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://tiktok.com/@user/video/123"
			args.Language = "english"
			err := args.Validate()
			Expect(errors.Is(err, transcription.ErrInvalidLanguageCode)).To(BeTrue())
		})

		It("should reject too many parts", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://tiktok.com/@user/video/123"
			args.Language = "en-us-extra"
			err := args.Validate()
			Expect(errors.Is(err, transcription.ErrInvalidLanguageCode)).To(BeTrue())
		})

		It("should reject too few parts", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://tiktok.com/@user/video/123"
			args.Language = "en"
			err := args.Validate()
			Expect(errors.Is(err, transcription.ErrInvalidLanguageCode)).To(BeTrue())
		})

		It("should reject invalid region length", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://tiktok.com/@user/video/123"
			args.Language = "en-usa"
			err := args.Validate()
			Expect(errors.Is(err, transcription.ErrInvalidLanguageCode)).To(BeTrue())
		})

		It("should reject invalid language length", func() {
			args := transcription.NewArguments()
			args.VideoURL = "https://tiktok.com/@user/video/123"
			args.Language = "english-us"
			err := args.Validate()
			Expect(errors.Is(err, transcription.ErrInvalidLanguageCode)).To(BeTrue())
		})
	})

	Describe("Helper methods", func() {
		It("should return true when language preference is set", func() {
			args := transcription.NewArguments()
			args.Language = "en-us"
			Expect(args.HasLanguagePreference()).To(BeTrue())
		})

		It("should return false when language preference is not set", func() {
			args := transcription.NewArguments()
			args.Language = ""
			Expect(args.HasLanguagePreference()).To(BeFalse())
		})

		It("should return the language code when set", func() {
			args := transcription.NewArguments()
			args.Language = "en-us"
			Expect(args.GetLanguageCode()).To(Equal("en-us"))
		})

		It("should return default language code when not set", func() {
			args := transcription.NewArguments()
			args.Language = ""
			args.SetDefaultValues()
			Expect(args.GetLanguageCode()).To(Equal("eng-US"))
		})

		It("should return the video URL", func() {
			expected := "https://tiktok.com/@user/video/123"
			args := transcription.NewArguments()
			args.VideoURL = expected
			Expect(args.GetVideoURL()).To(Equal(expected))
		})
	})
})
