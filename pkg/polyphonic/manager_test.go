package polyphonic

import (
	"fmt"
	"testing"
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
