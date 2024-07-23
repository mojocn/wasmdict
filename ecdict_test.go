package wasmecdict

import (
	"testing"
)

func TestLookUp(t *testing.T) {
	for _, s := range []string{"awesome", "America", "Europe", "China", "book", "joker", "polish", "Polish", "china", "China"} {
		word := LookUp(s)
		if word == nil {
			t.Errorf("Word %s not found", s)
		} else {
			t.Logf("Word %s found: %s, %s", s, word.Definition, word.Translation)
		}
	}
}
