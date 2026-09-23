package japconverter

import (
	"strings"
	"testing"

	"github.com/yukumo-group/Chinese2KanaConverter/internal/cpyconverter"
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

// TestWithCpyConverter tests the converting of chinese to kana with polyphonics
func TestWithCpyConverter(t *testing.T) {
	t.Parallel()
	data := map[string]string{
		"都会区": "dū huì qū",
	}
	resNoDump := cpyconverter.ToPinyin("西雅图都会区; 长夜漫漫, winter is coming!", true)
	t.Log(resNoDump)
	cpyconverter.DumpHeteronymMap(data)
	ExpectedResult := []string{
		"du",
		"hui",
		"qu",
	}
	res := cpyconverter.ToPinyin("都会区", true)
	for i, py := range res {
		if len(py) < 1 {
			t.Errorf(
				"Pinyin for charaacter %d not generated",
				i,
			)
		}
		if strings.TrimSpace(py) != strings.TrimSpace(ExpectedResult[i]) {
			t.Errorf(
				"Expected %s, got %s",
				ExpectedResult[i],
				py,
			)
		}
	}
	kanaRes, err := ToKana(
		res,
	)
	if err != nil {
		t.Error(err)
	}
	const expectedKanaRes string = "トゥーホイチュイ"
	if kanaRes != expectedKanaRes {
		t.Errorf(
			"expected %s, got %s",
			expectedKanaRes,
			kanaRes,
		)
	}
}
