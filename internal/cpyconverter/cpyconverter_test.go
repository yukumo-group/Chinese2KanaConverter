package cpyconverter

import (
	"strings"
	"testing"
)

// TestDumpHeteronymDict tests the dumping of heteronym dict
func TestDumpHeteronymDict(t *testing.T) {
	t.Parallel()
	err := InitGSEDict(
		"./testdata/dict.txt",
	)
	if err != nil {
		t.Error(err)
	}
	data := map[string]string{
		"都会区": "dū huì qū",
	}
	resNoDump := ToPinyin("西雅图都会区; 长夜漫漫, winter is coming!", true)
	t.Log(resNoDump)
	DumpHeteronymMap(data)
	ExpectedResult := []string{
		"du",
		"hui",
		"qu",
	}
	res := ToPinyin("都会区", true)
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
}
