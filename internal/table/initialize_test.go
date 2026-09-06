package table

import (
	"testing"
)

// TestInitialsAndFinals tests for initials and finals
func TestInitialsAndFinals(
	t *testing.T,
) {
	t.Parallel()
	if len(AllInitials) == 0 {
		t.Error(
			"the length of AllInitials cannot be 0",
		)
	}
	if len(AllFinals) == 0 {
		t.Error(
			"the length of AllFinals cannot be 0",
		)
	}
	for _, initial := range AllInitials {
		_, exists := InitialToKana[initial]
		if !exists {
			t.Errorf(
				"%s does not exists in the initial to kana map",
				initial,
			)
		}
	}
	for _, final := range AllFinals {
		_, exists := FinalToKana[final]
		if !exists {
			t.Errorf(
				"%s does not exists in the initial to kana map",
				final,
			)
		}
	}
	for i := 1; i < len(AllInitials); i++ {
		if len([]rune(AllInitials[i])) > len([]rune(AllInitials[i-1])) {
			t.Log(AllInitials)
			t.Errorf(
				"%d: %s has a length larger than %d:%s",
				i,
				AllInitials[i],
				i-1,
				AllInitials[i-1],
			)
		}
	}
	for i := 1; i < len(AllFinals); i++ {
		if len([]rune(AllFinals[i])) > len([]rune(AllFinals[i-1])) {
			t.Log(AllFinals)
			t.Errorf(
				"%d: %s has a length larger than %d:%s",
				i,
				AllFinals[i],
				i-1,
				AllFinals[i-1],
			)
		}
	}
}
