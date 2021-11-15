package porterstemmer

import (
	"bufio"
	"os"
	"testing"
)

type StemTest struct {
	Word []rune
	Stem []rune
}

func TestPrepareWord(t *testing.T) {
	tests := []StemTest{
		{[]rune("bellí"), []rune("bellì")},
		{[]rune("aiuola"), []rune("aIUola")},
		{[]rune("buio"), []rune("buIo")},
		{[]rune("báccánó"), []rune("bàccànò")},
		{[]rune("quadro"), []rune("qUadro")},
	}

	for _, test := range tests {
		result := prepareWord(test.Word)
		if !equal(result, test.Stem) {
			t.Fatalf("'%s' should be %q, is %q\n", string(test.Word), string(test.Stem), string(result))
		}
	}

}

func TestEqual(t *testing.T) {
	a := []rune("ciao")
	b := []rune("pippo")

	if equal(a, b) {
		t.Fail()
	}

	if !equal([]rune("peto"), []rune("peto")) {
		t.Fail()
	}
}

func TestHasSuffix(t *testing.T) {
	a := []rune("cippalippa")
	if !hasSuffix(a, []rune("lippa")) {
		t.Fail()
	}

	b := []rune("petofono")
	if hasSuffix(b, []rune("cazzy")) {
		t.Fail()
	}
}

func TestStemWord(t *testing.T) {
	tests := []StemTest{
		{[]rune("abbandonata"), []rune("abbandon")},
		{[]rune("abbandonate"), []rune("abbandon")},
		{[]rune("abbandonati"), []rune("abbandon")},
		{[]rune("abbandonato"), []rune("abbandon")},
		{[]rune("abbandonava"), []rune("abbandon")},
		{[]rune("abbandonerà"), []rune("abbandon")},
		{[]rune("abbandonerò"), []rune("abbandon")},
		{[]rune("abbandoneranno"), []rune("abbandon")},
		{[]rune("abbandono"), []rune("abband")},
		{[]rune("abbaruffato"), []rune("abbaruff")},
		{[]rune("abbassamento"), []rune("abbass")},
		{[]rune("propagarla"), []rune("propag")},
		{[]rune("propizio"), []rune("propiz")},
		{[]rune("propio"), []rune("prop")},
	}

	var rv []rune
	for _, test := range tests {
		rv = StemWithoutLowerCasing(test.Word)
		if !equal(rv, test.Stem) {
			t.Fatalf("Stem failed: '%s'->'%s' (should be: '%s')\n", string(test.Word), string(rv), string(test.Stem))
		}
	}
}

func TestFiles(t *testing.T) {
	inFile, err := os.Open("./testdata/voc.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer inFile.Close()

	outFile, err := os.Open("./testdata/output.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	scannerIn := bufio.NewScanner(inFile)
	scannerOut := bufio.NewScanner(outFile)

	for scannerIn.Scan() && scannerOut.Scan() {
		sIn := scannerIn.Text()
		sOut := scannerOut.Text()

		result := StemWithoutLowerCasing([]rune(sIn))
		if !equal(result, []rune(sOut)) {
			t.Fatalf("%q: should be %q, is %q\n", string(sIn), string(sOut), string(result))
		}
	}
}
