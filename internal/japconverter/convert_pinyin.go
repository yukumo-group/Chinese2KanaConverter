package japconverter

import (
	"github.com/yukumo-group/Chinese2KanaConverter/internal/table"
)

// ToKanaByInitialsAndFinals
func ToKanaByInitialsAndFinals(
	text string,
) string {
	return ""
}

// ToKana converts pinyin list to kana
func ToKana(
	pinyinList []string,
) string {
	result := ""
	for _, pinyin := range pinyinList {
		japKana, exists := table.PinyinToJapMap[pinyin]
		if exists {
			result += japKana
		}
	}
	return result
}
