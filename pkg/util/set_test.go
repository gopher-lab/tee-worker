package util_test

import (
	"slices"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-types/pkg/util"
)

var _ = Describe("Set", func() {
	It("should return a slice of all its elements", func() {
		s := util.NewSet(0, 1, 2, 2, 2, 3, 4, 5, 5, 5, 5, 6)
		Expect(s.Items()).To(ConsistOf(0, 1, 2, 3, 4, 5, 6))
		Expect(s.Length()).To(Equal(7))
	})

	It("should check whether an item is included in the Set or not", func() {
		s := util.NewSet(1, 2, 3, 4, 5, 6)
		Expect(s.Contains(2)).To(BeTrue())
		Expect(s.Contains(42)).To(BeFalse())
	})

	It("should add items to the set without duplicating", func() {
		s := util.NewSet(1, 2, 3, 4, 5, 6)
		s.Add(7, 8, 9, 2, 4)
		Expect(s.Items()).To(ConsistOf(1, 2, 3, 4, 5, 6, 7, 8, 9))
	})

	It("should delete items from the set if they exist", func() {
		s := util.NewSet(1, 2, 3, 4, 5, 6)
		s.Delete(7, 8, 9, 2, 4, 42)
		Expect(s.Items()).To(ConsistOf(1, 3, 5, 6))
	})

	It("should return the number of items in the set", func() {
		s := util.NewSet(0, 1, 2, 3, 4, 5, 6)
		Expect(s.Length()).To(Equal(7))
		s.Add(7, 8, 9, 2, 4)
		Expect(s.Length()).To(Equal(10))
		s.Delete(0, 1)
		Expect(s.Length()).To(Equal(8))
	})

	It("should return a sequence of all its elements", func() {
		s := util.NewSet(0, 1, 2, 3, 4, 5, 6)
		items := slices.Collect(s.ItemsSeq())
		Expect(items).To(ConsistOf(0, 1, 2, 3, 4, 5, 6))
	})

	It("should return the union of multiple sets", func() {
		s1 := util.NewSet(0, 1, 2, 3, 4)
		s2 := util.NewSet(3, 4, 5, 6, 7)
		s3 := util.NewSet(8, 9, 0)
		s4 := s1.Union(s2, s3)
		Expect(s4.Items()).To(ConsistOf(0, 1, 2, 3, 4, 5, 6, 7, 8, 9))
	})

	It("should return the intersection of two sets", func() {
		s1 := util.NewSet(0, 1, 2, 3, 4)
		s2 := util.NewSet(3, 4, 5, 6, 7)
		s3 := s1.Intersection(s2)
		Expect(s3.Items()).To(ConsistOf(3, 4))
	})

	It("should return the difference of two sets", func() {
		s1 := util.NewSet(0, 1, 2, 3, 4)
		s2 := util.NewSet(3, 4, 5, 6, 7)
		s3 := s1.Difference(s2)
		Expect(s3.Items()).To(ConsistOf(0, 1, 2))
	})
})
