package table

import (
	"slices"
)

// SortStringSlice sorts string slice
func SortStringSlice(
	stringSlice []string,
) []string {
	sortedSlice := slices.Clone(stringSlice)
	slices.SortFunc(
		sortedSlice,
		func(a string, b string) int {
			return len([]rune(b)) - len([]rune(a))
		},
	)
	return sortedSlice
}
