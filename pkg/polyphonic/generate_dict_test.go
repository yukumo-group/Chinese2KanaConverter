package polyphonic

import (
	"fmt"
	"os"
	"testing"
)

// TestDictWriting tests the dict writing
func TestDictWriting(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	filePath := fmt.Sprintf(
		"%s/%s",
		tmpDir,
		"test.txt",
	)
	err := WriteDict(
		map[string]string{
			"都会区": "du hui qu",
		},
		filePath,
	)
	if err != nil {
		t.Error(err)
	}
	data, err := os.ReadFile(
		filePath,
	)
	if err != nil {
		t.Error(err)
	}
	stringData := string(data)
	if stringData != "都会区 100 n" {
		t.Errorf(
			"expected %s, got %s",
			"都会区 100 n",
			stringData,
		)
	}
}
