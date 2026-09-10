package language

import (
	"testing"

	"github.com/yukumo-group/Chinese2KanaConverter/internal/process"
)

const testParagraph string = "abbccdd，,一一四五一四114514abcdabcd四四八八七七"

// TestChunk tests the chunking
func TestChunk(t *testing.T) {
	t.Parallel()
	testChunks := SeparateToChunks(testParagraph)
	if len(testChunks) != 4 {
		t.Errorf(
			"expected the length of the separated text to be %d, got %d",
			3,
			len(testChunks),
		)
	}
	expectedResult := []bool{
		false,
		true,
		false,
		true,
		false,
	}
	for i, chunk := range testChunks {
		if chunk.IsChinese != expectedResult[i] {
			t.Errorf(
				"expected %s for IsChinese of %d, got %s",
				process.ConvertBoolToString(chunk.IsChinese),
				i,
				process.ConvertBoolToString(expectedResult[i]),
			)
		}
	}
}
