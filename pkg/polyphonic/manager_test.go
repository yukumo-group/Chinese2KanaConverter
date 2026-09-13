package polyphonic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/yukumo-group/Chinese2KanaConverter/internal/cpyconverter"
)

// TestManager tests functions related to manager
func TestManager(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	fileName1 := fmt.Sprintf(
		"%s/%s",
		tmpDir,
		"114.json",
	)
	newManager := NewManager()
	newManager.AddPolyphonic(
		"都会区",
		"dū huì qū",
	)
	loadedData := newManager.GetData()
	_, exists := loadedData["都会区"]
	if !exists {
		t.Errorf(
			"%s does not exists",
			"都会区",
		)
	}
	newManager.SetTargetFile(
		fileName1,
	)
	err := newManager.Save()
	if err != nil {
		t.Error(err)
	}
	loadedManager, err := NewManagerFromFile(
		fileName1,
	)
	if err != nil {
		t.Error(err)
	}
	loadedData = loadedManager.GetData()
	_, exists = loadedData["都会区"]
	if !exists {
		t.Errorf(
			"%s does not exists",
			"都会区",
		)
	}
}

// TestLoadMultiples tests the loading for polyphonics
func TestLoadMultiples(t *testing.T) {
	t.Parallel()
	newManager := NewManager()
	newManager.AddPolyphonic(
		"都会区",
		"dū huì qū",
	)
	newManager.AddPolyphonic(
		"都会区",
		"dū huì qū",
	)
	newManager.Initialize()
	newManager.Initialize()
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
