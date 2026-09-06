package japconverter

import (
	"github.com/yukumo-group/Chinese2KanaConverter/internal/table"
)

// ToKanaByInitialsAndFinals converts chinese to kana through matching of initials and finals
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
