package polyphonic

import (
	"fmt"
	"os"
)

// WriteDict writes gse dictionary
func WriteDict(
	polyphonics map[string]string,
	targetFilePath string,
) error {
	result := ""
	chineseSlice := []string{}
	for chinese := range polyphonics {
		chineseSlice = append(chineseSlice, chinese)
	}
	for i, chinese := range chineseSlice {
		result += fmt.Sprintf(
			"%s 100 n",
			chinese,
		)
		if i != len(chineseSlice)-1 {
			result += "\n"
		}
	}
	resultByte := []byte(result)
	err := os.WriteFile(
		targetFilePath,
		resultByte,
		0644,
	)
	return err
}
