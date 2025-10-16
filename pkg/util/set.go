package util

import (
	"iter"
	"maps"
	"slices"
)

// Set is a generic collection of unique items.
type Set[T comparable] map[T]struct{}

// NewSet creates and returns a new Set with the given items, deduplicating them.
func NewSet[T comparable](items ...T) *Set[T] {
	ret := make(Set[T], len(items))
	ret.Add(items...)
	return &ret
}

// Contains checks if an item is present in the set.
func (s *Set[T]) Contains(item T) bool {
	_, exists := (*s)[item]
	return exists
}

// Add inserts the given items into the set, deduplicating them.
func (s *Set[T]) Add(items ...T) *Set[T] {
	for _, item := range items {
		(*s)[item] = struct{}{}
	}
	return s
}

// Delete removes the given items from the set if it contains them.
func (s *Set[T]) Delete(items ...T) *Set[T] {
	for _, item := range items {
		delete((*s), item)
	}
	return s
}

// Length returns the number of items in the set.
func (s *Set[T]) Length() int {
	return len(*s)
}

// Items returns a slice containing all the items in the set.
// The order of items in the slice is not guaranteed.
func (s *Set[T]) Items() []T {
	return slices.Collect(s.ItemsSeq())
}

// ItemsSeq returns an iterator that yields all the items in the set.
// The order of items is not guaranteed.
func (s *Set[T]) ItemsSeq() iter.Seq[T] {
	return maps.Keys(*s)
}

// Union returns a new set containing all the items from the original set and all the provided sets, deduplicating them.
func (s *Set[T]) Union(sets ...*Set[T]) *Set[T] {
	sum := s.Length()
	for _, ss := range sets {
		sum = sum + ss.Length()
	}

	ret := make(map[T]struct{}, sum)
	for k := range *s {
		ret[k] = struct{}{}
	}
	for _, ss := range sets {
		for k := range *ss {
			ret[k] = struct{}{}
		}
	}

	rs := Set[T](ret)
	return &rs
}

// Intersection returns a new set containing only the items that are present in both the original set and s2.
func (s *Set[T]) Intersection(s2 *Set[T]) *Set[T] {
	ret := make(map[T]struct{}, Min(s.Length(), s2.Length()))
	for k := range *s {
		if s2.Contains(k) {
			ret[k] = struct{}{}
		}
	}

	rs := Set[T](ret)
	return &rs
}

// Difference returns a new set containing items that are in the original set but not in s2.
func (s *Set[T]) Difference(s2 *Set[T]) *Set[T] {
	ret := make(map[T]struct{}, s.Length())
	for k := range *s {
		if !s2.Contains(k) {
			ret[k] = struct{}{}
		}
	}

	rs := Set[T](ret)
	return &rs
}
