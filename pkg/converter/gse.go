package converter

import (
	"github.com/yukumo-group/Chinese2KanaConverter/internal/cpyconverter"
)

// InitGSEDict initializes the gse dict
func InitGSEDict(
	paths ...string,
) error {
	return cpyconverter.InitGSEDict(
		paths...,
	)
}
