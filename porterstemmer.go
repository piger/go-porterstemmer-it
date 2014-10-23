// An implementation of the Porter Stemming algorithm for the italian language.
//
// See: http://snowball.tartarus.org/algorithms/italian/stemmer.html
package porterstemmer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

var s0suff []string = []string{"gliela", "gliele", "glieli", "glielo", "gliene", "sene", "mela", "mele", "meli", "melo", "mene", "tela", "tele", "teli", "telo", "tene", "cela", "cele", "celi", "celo", "cene", "vela", "vele", "veli", "velo", "vene", "gli", "ci", "la", "le", "li", "lo", "mi", "ne", "si", "ti", "vi"}

var verbsuff []string = []string{"irebbero", "erebbero", "issero", "iscono", "iscano", "iresti", "ireste", "iremmo", "irebbe", "iranno", "essero", "eresti", "ereste", "eremmo", "erebbe", "eranno", "assimo", "assero", "ivate", "ivano", "ivamo", "irono", "irete", "iremo", "evate", "evano", "evamo", "erono", "erete", "eremo", "avate", "avano", "avamo", "arono", "isco", "isci", "isce", "isca", "irò", "irà", "irei", "irai", "immo", "iamo", "erò", "erà", "erei", "erai", "endo", "endi", "ende", "enda", "emmo", "assi", "asse", "ando", "ammo", "Yamo", "uto", "uti", "ute", "uta", "ono", "ivo", "ivi", "iva", "ito", "iti", "ite", "ita", "ire", "evo", "evi", "eva", "ete", "ere", "avo", "avi", "ava", "ato", "ati", "ate", "ata", "are", "ano", "ir", "ar"}

func isVowel(c rune) bool {
	switch c {
	case 'a', 'e', 'i', 'o', 'u', 'à', 'è', 'ì', 'ò', 'ù':
		return true
	}
	return false
}

// findRegion returns the start of the first region in word; a region starts from a
// non-vowel charcter followed by a vowel.
func findRegion(word string, start int) int {
	l := len(word)
	var oldr rune

	for i, w := 0, 0; i < l; i += w {
		r, rlen := utf8.DecodeRuneInString(word[i:])
		w = rlen

		if i == 0 || i < start {
			oldr = r
			continue
		}
		if !isVowel(r) && isVowel(oldr) && i+rlen < l {
			return i + rlen
		}
		oldr = r
	}

	return l
}

// findR12 finds the R1 and R2 region for a word.
//   - R1 is the region after the first non-vowel following a vowel, or is the null
//     region at the end of the word if there is no such non-vowel.
//   - R2 is the region after the first non-vowel following a vowel in R1, or is the
//     null region at the end of the word if there is no such non-vowel.
func findR12(word string) (int, int) {
	l := len(word)
	r1 := findRegion(word, 0)
	if r1 == l {
		return l, l
	}

	r2 := findRegion(word, r1)
	return r1, r2
}

// findRV returns the RV region.
// If the second letter is a consonant, RV is the region after the
// next following vowel, or if the first two letters are vowels, RV is
// the region after the next consonant, and otherwise (consonant-vowel
// case) RV is the region after the third letter. But RV is the end of
// the word if these positions cannot be found.
func findRV(word string) int {
	l := len(word)

	ch1, ch1Len := utf8.DecodeRuneInString(word)
	ch2, ch2Len := utf8.DecodeRuneInString(word[ch1Len:])
	_, ch3Len := utf8.DecodeRuneInString(word[ch1Len+ch2Len:])

	if !isVowel(ch2) {
		for i, w := 0, 0; i < l; i += w {
			r, width := utf8.DecodeRuneInString(word[i:])
			w = width

			if i < ch2Len {
				continue
			}
			if isVowel(r) {
				if i+width < l {
					return i + width
				} else {
					return l
				}
			}
		}
		return l
	} else if isVowel(ch1) && isVowel(ch2) {
		for i, w := 0, 0; i < l; i += w {
			r, rlen := utf8.DecodeRuneInString(word[i:])
			w = rlen

			if i < ch1Len+ch2Len {
				continue
			}
			if !isVowel(r) {
				return i + rlen
			}
		}
		return l
	}

	// "and otherwise (consonant-vowel case) RV is the region after the third letter."
	return ch1Len + ch2Len + ch3Len
}

// Step 0) Attached pronoun
func step0(word string) string {
	rv := findRV(word)

	for _, suff2 := range s0suff {
		var p int

		for _, suff1 := range []string{"ando", "endo"} {
			if strings.HasSuffix(word, suff1+suff2) {
				p = strings.LastIndex(word, suff1+suff2)
				if p != -1 && p >= rv {
					return word[0 : p+len(suff1)]
				}
			}
		}

		for _, suff1 := range []string{"ar", "er", "ir"} {
			if strings.HasSuffix(word, suff1+suff2) {
				p = strings.LastIndex(word, suff1+suff2)
				if p != -1 && p >= rv {
					return fmt.Sprintf("%se", word[0:p+len(suff1)])
				}
			}
		}
	}

	return word
}

// Step 1) Standard suffix removal:
// search for the longest among the following suffixes, and perform the action indicated.
func step1(word string) string {
	var p int
	r1, r2 := findR12(word)
	rv := findRV(word)

	if strings.HasSuffix(word, "amente") {
		p = strings.LastIndex(word, "amente")
		if p >= r1 {
			word = word[0:p]

			if strings.HasSuffix(word, "iv") {
				p = strings.LastIndex(word, "iv")
				if p >= r2 {
					word = word[0:p]

					if strings.HasSuffix(word, "at") {
						p = strings.LastIndex(word, "at")
						if p >= r2 {
							word = word[0:p]
						}
					}
				}
			} else {
				word, _ = replaceInRegion(word, []string{"os", "ic", "abil"}, "", r2)
			}
		}
		return word
	}

	for _, s := range []string{"atrici", "atrice", "mente", "istì", "istè", "istà", "ibili", "ibile", "abili", "abile", "isti", "iste", "ista", "ismo", "ismi", "ichi", "iche", "anze", "anza", "anti", "ante", "oso", "osi", "ose", "osa", "ico", "ici", "ice", "ica"} {
		if strings.HasSuffix(word, s) {
			p = strings.LastIndex(word, s)
			if p >= r2 {
				return word[0:p]
			} else {
				break
			}
		}
	}

	for _, s := range []string{"azione", "azioni", "atore", "atori"} {
		if strings.HasSuffix(word, s) {
			p = strings.LastIndex(word, s)
			if p >= r2 {
				word = word[0:p]

				if strings.HasSuffix(word, "ic") {
					p = strings.LastIndex(word, "ic")
					if p >= r2 {
						return word[0:p]
					}
				}
				return word
			}
		}
	}

	var repl bool
	if word, repl = replaceInRegion(word, []string{"logia", "logie"}, "log", r2); repl {
		return word
	}
	if word, repl = replaceInRegion(word, []string{"uzione", "uzioni", "usione", "usioni"}, "u", r2); repl {
		return word
	}
	if word, repl = replaceInRegion(word, []string{"enza", "enze"}, "ente", r2); repl {
		return word
	}
	if word, repl = replaceInRegion(word, []string{"amento", "amenti", "imento", "imenti"}, "", rv); repl {
		return word
	}

	if strings.HasSuffix(word, "ità") {
		p = strings.LastIndex(word, "ità")
		if p >= r2 {
			word = word[0:p]
		}
		word, _ = replaceInRegion(word, []string{"abil", "ic", "iv"}, "", r2)
		return word
	}

	for _, s := range []string{"ivo", "ivi", "iva", "ive"} {
		if strings.HasSuffix(word, s) {
			p = strings.LastIndex(word, s)
			if p >= r2 {
				word = word[0:p]

				if strings.HasSuffix(word, "at") {
					p = strings.LastIndex(word, "at")
					if p >= r2 {
						word = word[0:p]

						if strings.HasSuffix(word, "ic") {
							p = strings.LastIndex(word, "ic")
							if p >= r2 {
								word = word[0:p]
							}
						}
					}
				}
			}
		}
	}

	return word
}

// Step 2) Verb suffixes:
// search for the longest among the following suffixes in *RV*, and if found, delete.
func step2(word string) string {
	rv := findRV(word)

	for _, s := range verbsuff {
		if strings.HasSuffix(word[rv:], s) {
			p := strings.LastIndex(word, s)
			return word[0:p]
		}
	}
	return word
}

func step3a(word string) string {
	rv := findRV(word)

	for _, s := range []string{"a", "e", "i", "o", "à", "è", "ì", "ò"} {
		if !strings.HasSuffix(word, s) {
			continue
		}
		p := strings.LastIndex(word, s)
		if p >= rv {
			word = word[0:p]

			if strings.HasSuffix(word, "i") {
				pp := strings.LastIndex(word, "i")
				if pp >= rv {
					word = word[0:pp]
					return word
				}
			} else {
				return word
			}
		}
	}
	return word
}

func step3b(word string) string {
	rv := findRV(word)
	var p int

	if strings.HasSuffix(word, "ch") {
		p = strings.LastIndex(word, "ch")
		if p >= rv {
			return fmt.Sprintf("%sc", word[0:p])
		}
	}

	if strings.HasSuffix(word, "gh") {
		p = strings.LastIndex(word, "gh")
		if p >= rv {
			return fmt.Sprintf("%sg", word[0:p])
		}
	}

	return word
}

func restoreString(word string) string {
	word = strings.Replace(word, "I", "i", -1)
	word = strings.Replace(word, "U", "u", -1)
	return word
}

// replaceInRegion search for suffixes inside the region defined by rX inside word
// and replaces the suffix with repl; returns the modified word and true when it
// was modified.
func replaceInRegion(word string, suffixes []string, repl string, rX int) (string, bool) {
	var p int

	for _, s := range suffixes {
		if strings.HasSuffix(word, s) {
			p = strings.LastIndex(word, s)
			if p >= rX {
				return fmt.Sprintf("%s%s", word[0:p], repl), true
			}
		}
	}
	return word, false
}

// prepareWord returns a string "prepared" for stemming; it's the first step.
func prepareWord(word string) string {
	word = strings.Replace(word, "á", "à", -1)
	word = strings.Replace(word, "é", "è", -1)
	word = strings.Replace(word, "í", "ì", -1)
	word = strings.Replace(word, "ó", "ò", -1)
	word = strings.Replace(word, "ú", "ù", -1)

	word = strings.Replace(word, "qu", "qU", -1)

	var oldr rune
	var newword []rune
	for i, w := 0, 0; i < len(word); i += w {
		r, rlen := utf8.DecodeRuneInString(word[i:])
		w = rlen

		if i == 0 {
			oldr = r
			newword = append(newword, r)
			continue
		}
		if (r == 'i' || r == 'u') && isVowel(oldr) && i+rlen < len(word) {
			nr, _ := utf8.DecodeRuneInString(word[i+rlen:])
			if isVowel(nr) {
				if r == 'i' {
					newword = append(newword, 'I')
					continue
				} else if r == 'u' {
					newword = append(newword, 'U')
					continue
				}
			}
		}
		oldr = r
		newword = append(newword, r)
	}
	return string(newword)
}

// StemWord returns the italian stemming for the word w.
func StemString(w string) string {
	word := prepareWord(w)
	word0 := step0(word)
	word1 := step1(word0)

	if word0 == word1 {
		word = step2(word1)
	} else {
		word = word1
	}

	word = step3a(word)
	word = step3b(word)
	word = restoreString(word)
	return word
}
