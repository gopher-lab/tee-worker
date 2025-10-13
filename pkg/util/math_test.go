package util_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-types/pkg/util"
)

var _ = Describe("Math functions", func() {
	Describe("Min", func() {
		It("should calculate the minimum of a series of orderable values regardless of parameter order", func() {
			Expect(util.Min(1, 2, 3, 4, 5, 6)).To(Equal(1))
			Expect(util.Min(2, 3, 8, -1, 4, 42)).To(Equal(-1))
		})
	})

	Describe("Max", func() {
		It("should calculate the maximum of a series of orderable values regardless of parameter order", func() {
			Expect(util.Max(1, 2, 3, 4, 5, 6)).To(Equal(6))
			Expect(util.Max(2, 3, 8, -12, 4, 42)).To(Equal(42))
		})
	})
})
