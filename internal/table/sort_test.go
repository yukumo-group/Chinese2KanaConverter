package table

import (
	"testing"
)

// TestSort tests the sort of string slice
func TestSort(t *testing.T) {
	t.Parallel()
	testStrings := []string{
		"bca",
		"cbadc",
		"e",
		"kksk",
		"cj",
	}
	expectedResult := []string{
		"cbadc",
		"kksk",
		"bca",
		"cj",
		"e",
	}
	result := SortStringSlice(
		testStrings,
	)
	for i, s := range result {
		if s != expectedResult[i] {
			t.Errorf(
				"expected postion %d of the result to be %s, got %s",
				i,
				expectedResult[i],
				s,
			)
		}
	}
}
