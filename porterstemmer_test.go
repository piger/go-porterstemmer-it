package porterstemmer

import (
	"bufio"
	"os"
	"testing"
)

func testWord(old, new string, t *testing.T) {
	result := prepareWord(old)
	if result != new {
		t.Fatalf("%s: should be %q, is %q\n", old, new, result)
	}
}

func TestPrepareWord(t *testing.T) {
	testWord("bellí", "bellì", t)
	testWord("aiuola", "aIUola", t)
	testWord("buio", "buIo", t)
	testWord("báccánó", "bàccànò", t)
	testWord("quadro", "qUadro", t)
}

func TestR12(t *testing.T) {
	var words [][]string = [][]string{
		{"beautiful", "iful", "ul"},
		{"beauty", "y", ""},
		{"beau", "", ""},
		{"animadversion", "imadversion", "adversion"},
		{"sprinkled", "kled", ""},
		{"eucharist", "harist", "ist"},
		{"giocatrici", "atrici", "rici"},
	}

	for _, word := range words {
		r1, r2 := findR12(word[0])
		if word[1] != word[0][r1:] {
			t.Fatalf("should be %q, is %q\n", word[1], word[0][r1:])
		}
		if word[2] != word[0][r2:] {
			t.Fatalf("should be %q, is %q\n", word[2], word[0][r2:])
		}
	}
}

func TestStemWord(t *testing.T) {

	var words [][]string = [][]string{
		{"abbandonata", "abbandon"},
		{"abbandonate", "abbandon"},
		{"abbandonati", "abbandon"},
		{"abbandonato", "abbandon"},
		{"abbandonava", "abbandon"},
		{"abbandonerà", "abbandon"},
		{"abbandonerò", "abbandon"},
		{"abbandoneranno", "abbandon"},
		{"abbandono", "abband"},
		{"abbaruffato", "abbaruff"},
		{"abbassamento", "abbass"},
		{"propagarla", "propag"},
		{"propizio", "propiz"},
		{"propio", "prop"},
	}

	for _, word := range words {
		result := StemWord(word[0])
		if result != word[1] {
			t.Fatalf("%q: should be %q, is %q\n", word[0], word[1], result)
		}
	}
}

func TestFiles(t *testing.T) {
	inFile, err := os.Open("voc.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer inFile.Close()

	outFile, err := os.Open("output.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	scannerIn := bufio.NewScanner(inFile)
	scannerOut := bufio.NewScanner(outFile)

	for scannerIn.Scan() && scannerOut.Scan() {
		sIn := scannerIn.Text()
		sOut := scannerOut.Text()

		result := StemWord(sIn)
		if result != sOut {
			t.Fatalf("%q: should be %q, is %q\n", sIn, sOut, result)
		}
	}
}
