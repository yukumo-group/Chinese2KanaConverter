package cpyconverter

import (
	"sync"

	"github.com/go-ego/gpy/phrase"
)

// mapProtector protects the map when doing read/write
var mapProtector sync.RWMutex

// InitGSEDict inits gse dictionary
func InitGSEDict(
	gseDictPath ...string,
) error {
	err := phrase.LoadGseDict(gseDictPath...)
	return err
}

// DumpHeteronymMap dumps map of heternym to the converter
func DumpHeteronymMap(
	heteronymMap map[string]string,
) {
	mapProtector.Lock()
	defer mapProtector.Unlock()
	for chineseText, pinyin := range heteronymMap {
		phrase.AddDict(
			chineseText,
			pinyin,
		)
	}
}
