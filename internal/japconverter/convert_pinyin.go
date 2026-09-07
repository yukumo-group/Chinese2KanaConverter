package japconverter

import (
	"fmt"
	"strings"

	"github.com/yukumo-group/Chinese2KanaConverter/internal/table"
)

// ToKanaByInitialsAndFinals converts chinese to kana through matching of initials and finals
func ToKanaByInitialsAndFinals(
	text string,
) (string, error) {
	withReplacecFinalsText := text
	for _, final := range table.AllFinals {
		kana, exists := table.FinalToKana[final]
		if !exists {
			return "", fmt.Errorf(
				"%s does not exists in the finals to kana map",
				final,
			)
		}
		withReplacecFinalsText = strings.ReplaceAll(
			withReplacecFinalsText,
			final,
			strings.TrimSpace(kana),
		)
	}
	allProcessedString := withReplacecFinalsText
	for _, initial := range table.AllInitials {
		kana, exists := table.InitialToKana[initial]
		if !exists {
			return "", fmt.Errorf(
				"%s does not exists in the initial to kana map",
				initial,
			)
		}
		allProcessedString = strings.ReplaceAll(
			allProcessedString,
			initial,
			strings.TrimSpace(kana),
		)
	}
	return allProcessedString, nil
}

// ToKana converts pinyin list to kana
func ToKana(
	pinyinList []string,
) (string, error) {
	result := ""
	for _, pinyin := range pinyinList {
		processedPinyin := strings.TrimSpace(
			pinyin,
		)
		japKana, exists := table.PinyinToJapMap[processedPinyin]
		if exists {
			result += strings.TrimSpace(japKana)
		} else {
			resultByMatching, err := ToKanaByInitialsAndFinals(
				pinyin,
			)
			if err != nil {
				return "", err
			}
			result += strings.TrimSpace(resultByMatching)
		}
	}
	return result, nil
}
