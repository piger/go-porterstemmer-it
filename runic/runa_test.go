package runic

import (
	"testing"
)

func TestEqual(t *testing.T) {
	a := []rune("ciao")
	b := []rune("pippo")

	if Equal(a, b) {
		t.Fail()
	}

	if !Equal([]rune("peto"), []rune("peto")) {
		t.Fail()
	}
}

func TestHasSuffix(t *testing.T) {
	a := []rune("cippalippa")
	if !HasSuffix(a, []rune("lippa")) {
		t.Fail()
	}

	b := []rune("petofono")
	if HasSuffix(b, []rune("cazzy")) {
		t.Fail()
	}
}

type StemTest struct {
	Word []rune
	Stem []rune
}

func TestStemWord(t *testing.T) {
	tests := []StemTest{
		StemTest{[]rune("abbandonata"), []rune("abbandon")},
		StemTest{[]rune("abbandonate"), []rune("abbandon")},
		StemTest{[]rune("abbandonati"), []rune("abbandon")},
		StemTest{[]rune("abbandonato"), []rune("abbandon")},
		StemTest{[]rune("abbandonava"), []rune("abbandon")},
		StemTest{[]rune("abbandonerà"), []rune("abbandon")},
		StemTest{[]rune("abbandonerò"), []rune("abbandon")},
		StemTest{[]rune("abbandoneranno"), []rune("abbandon")},
		StemTest{[]rune("abbandono"), []rune("abband")},
		StemTest{[]rune("abbaruffato"), []rune("abbaruff")},
		StemTest{[]rune("abbassamento"), []rune("abbass")},
		StemTest{[]rune("propagarla"), []rune("propag")},
		StemTest{[]rune("propizio"), []rune("propiz")},
		StemTest{[]rune("propio"), []rune("prop")},
	}

	var rv []rune
	for _, test := range tests {
		rv = StemWord(test.Word)
		if !Equal(rv, test.Stem) {
			t.Fatalf("Stem failed: '%s'->'%s' (should be: '%s')\n", string(test.Word), string(rv), string(test.Stem))
		}
	}
}
