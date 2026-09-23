package converter

import (
	"testing"

	"github.com/yukumo-group/Chinese2KanaConverter/pkg/polyphonic"
)

// TestWithManager tests the converter with manager
func TestWithManager(t *testing.T) {
	t.Parallel()
	newManager := polyphonic.NewManager()
	newManager.AddPolyphonic(
		"都会区",
		"dū huì qū",
	)
	const expectedKanaRes string = "トゥーホイチュイ"
	res, err := SingleChinesePieceToKana(
		"都会区",
		true,
	)
	if err != nil {
		t.Error(err)
	}
	if res != expectedKanaRes {
		t.Errorf(
			"expected %s, got %s",
			expectedKanaRes,
			res,
		)
	}
}
