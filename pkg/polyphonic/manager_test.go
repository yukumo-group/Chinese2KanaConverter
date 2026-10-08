package polyphonic

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/yukumo-group/Chinese2KanaConverter/internal/cpyconverter"
	"github.com/yukumo-group/Chinese2KanaConverter/pkg/converter"
)

// TestManager tests functions related to manager
func TestManager(t *testing.T) {
	t.Parallel()
	err := cpyconverter.InitGSEDict(
		"./testdata/dict.txt",
	)
	if err != nil {
		t.Error(err)
	}
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
		"./testdata/dict.txt",
	)
	err = newManager.Save()
	if err != nil {
		t.Error(err)
	}
	loadedManager, err := NewManagerFromFile(
		fileName1,
		"./testdata/dict.txt",
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
	err := cpyconverter.InitGSEDict(
		"./testdata/dict.txt",
	)
	if err != nil {
		t.Error(err)
	}
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

// TestWithDict tests polyphonics manager with gse dict generation
func TestWithDict(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	fileName1 := fmt.Sprintf(
		"%s/%s",
		tmpDir,
		"114.json",
	)
	dictName := fmt.Sprintf(
		"%s/%s",
		tmpDir,
		"dict.txt",
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
		dictName,
	)
	err := newManager.Save()
	if err != nil {
		t.Error(err)
	}
	err = newManager.SaveGSEDict()
	if err != nil {
		t.Error(err)
	}
	// Load new manager
	fmt.Println("start loading new manager")
	reloadedManager, err := NewManagerFromFile(
		fileName1,
		dictName,
	)
	if err != nil {
		t.Error(err)
	}
	data := reloadedManager.GetData()
	_, exists = data["都会区"]
	if !exists {
		t.Errorf(
			"%s does not exists in data",
			"都会区",
		)
	}
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

// TestTotalNewManager tests loading of new manager
func TestTotalNewManager(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	fileName1 := fmt.Sprintf(
		"%s/%s",
		tmpDir,
		"114.json",
	)
	dictName := fmt.Sprintf(
		"%s/%s",
		tmpDir,
		"dict.txt",
	)
	_, err := NewManagerFromFile(
		fileName1,
		dictName,
	)
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat(
		fileName1,
	)
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat(
		dictName,
	)
	if err != nil {
		t.Error(err)
	}
}

// TestChineseConvert tests the converting
func TestChineseConvert(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	targetFilePath := fmt.Sprintf(
		"%s/%s",
		tmpDir,
		"test.json",
	)
	dictFilePath := fmt.Sprintf(
		"%s/%s",
		tmpDir,
		"test.txt",
	)
	newPolyphonicManager, err := NewManagerFromFile(
		targetFilePath,
		dictFilePath,
	)
	if err != nil {
		t.Error(err)
	}
	newPolyphonicManager.AddPolyphonic(
		"都市",
		"du shi",
	)
	newPolyphonicManager.AddPolyphonic(
		"银行",
		"yin hang",
	)
	err = newPolyphonicManager.SaveGSEDict()
	if err != nil {
		t.Error(err)
	}
	result, err := converter.SingleChinesePieceToKana(
		"都市银行",
		true,
	)
	if err != nil {
		t.Error(
			err,
		)
	}
	const expecteResult string = "トゥーシーインハン"
	if result != expecteResult {
		t.Errorf(
			"expected %s, got %s",
			expecteResult,
			result,
		)
	}
}

// TestPolyphonicsCRUD tests the crud function of polyphonics
func TestPolyphonicsCRUD(t *testing.T) {
	t.Parallel()
}
