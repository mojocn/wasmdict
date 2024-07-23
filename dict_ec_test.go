package wasmdict

import (
	"testing"
)

func TestEcLookUp(t *testing.T) {
	for _, s := range []string{"awesome", "America", "Europe", "China", "book", "joker", "polish", "Polish", "china", "China"} {
		word := EcLookUp(s)
		if word == nil {
			t.Errorf("Word %s not found", s)
		} else {
			t.Logf("Word %s found: %s, %s", s, word.Definition, word.Translation)
		}
	}
}
