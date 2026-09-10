package converter

import (
	"github.com/yukumo-group/Chinese2KanaConverter/internal/cpyconverter"
	"github.com/yukumo-group/Chinese2KanaConverter/internal/japconverter"
	"github.com/yukumo-group/Chinese2KanaConverter/internal/language"
)

// SingleChinesePieceToKana converts single piece of Chinese to kana
func SingleChinesePieceToKana(
	chineseText string,
	useHeteronym bool,
) (string, error) {
	result := ""
	chunks := language.SeparateToChunks(
		chineseText,
	)
	for _, chunk := range chunks {
		if chunk.IsChinese {
			pinyins := cpyconverter.ToPinyin(
				chunk.Text,
				useHeteronym,
			)
			kanas, err := japconverter.ToKana(
				pinyins,
			)
			if err != nil {
				return result, err
			}
			result += kanas
		} else {
			result += japconverter.OthersToKana(
				chunk.Text,
			)
		}
	}
	return result, nil
}
