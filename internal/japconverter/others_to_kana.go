package japconverter

import (
	kanatrans "github.com/Luigi-Pizzolito/English2KanaTransliteration"
)

// OthersToKana converts other language (English and japanese) to kana
func OthersToKana(
	rawText string,
) string {
	converter := kanatrans.NewAllToKana()
	result := converter.Convert(rawText)
	return result
}
