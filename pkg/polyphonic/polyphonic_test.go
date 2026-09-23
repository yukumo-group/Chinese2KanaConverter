package polyphonic

import (
	"strings"
	"testing"

	"github.com/yukumo-group/Chinese2KanaConverter/internal/cpyconverter"
)

// TestLoad tests the loading of heteronym dict
func TestLoad(t *testing.T) {
	t.Parallel()
	err := cpyconverter.InitGSEDict(
		"./testdata/dict.txt",
	)
	if err != nil {
		t.Error(err)
	}
	data := map[string]string{
		"都会区": "du hui qu",
	}
	LoadPolyphonics(data)
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
}

// TestSafeLoad tests the safe loading
func TestSafeLoad(t *testing.T) {
	err := cpyconverter.InitGSEDict(
		"./testdata/dict.txt",
	)
	if err != nil {
		t.Error(err)
	}
	resultChan := make(chan int, 2)
	data := map[string]string{
		"都会区": "du hui qu",
	}
	func1 := func() {
		SafeLoadPolyphonics(data)
		resultChan <- 0
	}
	func2 := func() {
		SafeLoadPolyphonics(data)
		resultChan <- 0
	}
	go func1()
	go func2()
	for i := 0; i < 2; i++ {
		<-resultChan
	}
}
