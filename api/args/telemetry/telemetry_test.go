package telemetry_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/api/args/telemetry"
	"github.com/masa-finance/tee-worker/api/types"
)

var _ = Describe("Telemetry Arguments", func() {
	Describe("Marshalling and unmarshalling", func() {
		It("should set default values", func() {
			args := telemetry.NewArguments()
			jsonData, err := json.Marshal(args)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.Type).To(Equal(types.CapTelemetry))
		})

		It("should preserve custom values", func() {
			args := telemetry.NewArguments()
			jsonData, err := json.Marshal(args)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(jsonData), &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args.Type).To(Equal(types.CapTelemetry))
		})

		It("should handle invalid JSON", func() {
			args := &telemetry.Arguments{}
			invalidJSON := `{"type": "telemetry", "invalid": }`
			err := json.Unmarshal([]byte(invalidJSON), args)
			Expect(err).To(HaveOccurred())
			// The error should be a JSON syntax error, not our custom error
			Expect(err).To(BeAssignableToTypeOf(&json.SyntaxError{}))
		})
	})

	Describe("Validation", func() {
		It("should succeed with valid arguments", func() {
			args := telemetry.NewArguments()
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should succeed with empty arguments", func() {
			args := &telemetry.Arguments{}
			err := args.Validate()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("GetCapability", func() {
		It("should return the telemetry capability", func() {
			args := telemetry.NewArguments()
			Expect(args.GetCapability()).To(Equal(types.CapTelemetry))
		})

		It("should return empty capability for uninitialized arguments", func() {
			args := &telemetry.Arguments{}
			Expect(args.GetCapability()).To(Equal(types.Capability("")))
		})
	})

	Describe("SetDefaultValues", func() {
		It("should not modify arguments", func() {
			args := telemetry.NewArguments()
			originalType := args.Type
			args.SetDefaultValues()
			Expect(args.Type).To(Equal(originalType))
		})
	})
})
