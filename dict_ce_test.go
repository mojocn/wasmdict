package wasmdict

import (
	"testing"
)

func TestCeLookUp(t *testing.T) {
	entryCE := CeLookUp("過", false)
	if entryCE == nil {
		t.Error("not found")
	} else {
		t.Logf("%v", entryCE)
	}

	entryCE = CeLookUp("過", true)
	if entryCE != nil {
		t.Error("not found")
	}

	entryCE = CeLookUp("过", true)
	if entryCE == nil {
		t.Error("not found")
	} else {
		t.Logf("%v", entryCE)
	}

}

func TestCeQueryLike(t *testing.T) {
	words := CeQueryLike("过", true, 10)
	for _, word := range words {
		t.Logf("%v", word)
	}

	words = CeQueryLike("楊", false, 10)
	for _, word := range words {
		t.Logf("%v", word)
	}
}
