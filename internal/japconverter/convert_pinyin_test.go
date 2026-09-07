package japconverter

import (
	"testing"
)

// TestToKanaByInitialsAndFinals tests the function of converting pinyin to kana by matching initials and finals
func TestToKanaByInitialsAndFinals(
	t *testing.T,
) {
	t.Parallel()
	result, err := ToKanaByInitialsAndFinals(
		"ban",
	)
	if err != nil {
		t.Error(err)
	}
	if result != "パン" {
		t.Errorf(
			"expected %s, got %s",
			"パン",
			result,
		)
	}
	result2, err := ToKanaByInitialsAndFinals(
		"zhang",
	)
	if err != nil {
		t.Error(err)
	}
	if result2 != "チャン" {
		t.Errorf(
			"expected %s, got %s",
			"チャン",
			result2,
		)
	}
}

// TestToKana tests the normal to kana function
func TestToKana(t *testing.T) {
	t.Parallel()
	testData := []string{
		"ban",
		"ji",
	}
	result, err := ToKana(testData)
	if err != nil {
		t.Error(err)
	}
	if result != "パンチー" {
		t.Errorf(
			"expected %s, got %s",
			"パンチー",
			result,
		)
	}
}
