// An implementation of the Porter Stemming algorithm for the italian language.
//
// See: http://snowball.tartarus.org/algorithms/italian/stemmer.html
package porterstemmer

import (
	"unicode"
)

var s0suff [][]rune = [][]rune{
	[]rune("gliela"),
	[]rune("gliele"),
	[]rune("glieli"),
	[]rune("glielo"),
	[]rune("gliene"),
	[]rune("sene"),
	[]rune("mela"),
	[]rune("mele"),
	[]rune("meli"),
	[]rune("melo"),
	[]rune("mene"),
	[]rune("tela"),
	[]rune("tele"),
	[]rune("teli"),
	[]rune("telo"),
	[]rune("tene"),
	[]rune("cela"),
	[]rune("cele"),
	[]rune("celi"),
	[]rune("celo"),
	[]rune("cene"),
	[]rune("vela"),
	[]rune("vele"),
	[]rune("veli"),
	[]rune("velo"),
	[]rune("vene"),
	[]rune("gli"),
	[]rune("ci"),
	[]rune("la"),
	[]rune("le"),
	[]rune("li"),
	[]rune("lo"),
	[]rune("mi"),
	[]rune("ne"),
	[]rune("si"),
	[]rune("ti"),
	[]rune("vi"),
}

var step1suffs [][]rune = [][]rune{
	[]rune("atrici"),
	[]rune("atrice"),
	[]rune("mente"),
	[]rune("ibili"),
	[]rune("ibile"),
	[]rune("abili"),
	[]rune("abile"),
	[]rune("istì"),
	[]rune("istè"),
	[]rune("istà"),
	[]rune("isti"),
	[]rune("iste"),
	[]rune("ista"),
	[]rune("ismo"),
	[]rune("ismi"),
	[]rune("ichi"),
	[]rune("iche"),
	[]rune("anze"),
	[]rune("anza"),
	[]rune("anti"),
	[]rune("ante"),
	[]rune("oso"),
	[]rune("osi"),
	[]rune("ose"),
	[]rune("osa"),
	[]rune("ico"),
	[]rune("ici"),
	[]rune("ice"),
	[]rune("ica"),
}

var verbsuff [][]rune = [][]rune{
	[]rune("irebbero"),
	[]rune("erebbero"),
	[]rune("issero"),
	[]rune("iscono"),
	[]rune("iscano"),
	[]rune("iresti"),
	[]rune("ireste"),
	[]rune("iremmo"),
	[]rune("irebbe"),
	[]rune("iranno"),
	[]rune("essero"),
	[]rune("eresti"),
	[]rune("ereste"),
	[]rune("eremmo"),
	[]rune("erebbe"),
	[]rune("eranno"),
	[]rune("assimo"),
	[]rune("assero"),
	[]rune("ivate"),
	[]rune("ivano"),
	[]rune("ivamo"),
	[]rune("irono"),
	[]rune("irete"),
	[]rune("iremo"),
	[]rune("evate"),
	[]rune("evano"),
	[]rune("evamo"),
	[]rune("erono"),
	[]rune("erete"),
	[]rune("eremo"),
	[]rune("avate"),
	[]rune("avano"),
	[]rune("avamo"),
	[]rune("arono"),
	[]rune("isco"),
	[]rune("isci"),
	[]rune("isce"),
	[]rune("isca"),
	[]rune("irei"),
	[]rune("irai"),
	[]rune("immo"),
	[]rune("iamo"),
	[]rune("erei"),
	[]rune("erai"),
	[]rune("endo"),
	[]rune("endi"),
	[]rune("ende"),
	[]rune("enda"),
	[]rune("emmo"),
	[]rune("assi"),
	[]rune("asse"),
	[]rune("ando"),
	[]rune("ammo"),
	[]rune("Yamo"),
	[]rune("irò"),
	[]rune("irà"),
	[]rune("erò"),
	[]rune("erà"),
	[]rune("uto"),
	[]rune("uti"),
	[]rune("ute"),
	[]rune("uta"),
	[]rune("ono"),
	[]rune("ivo"),
	[]rune("ivi"),
	[]rune("iva"),
	[]rune("ito"),
	[]rune("iti"),
	[]rune("ite"),
	[]rune("ita"),
	[]rune("ire"),
	[]rune("evo"),
	[]rune("evi"),
	[]rune("eva"),
	[]rune("ete"),
	[]rune("ere"),
	[]rune("avo"),
	[]rune("avi"),
	[]rune("ava"),
	[]rune("ato"),
	[]rune("ati"),
	[]rune("ate"),
	[]rune("ata"),
	[]rune("are"),
	[]rune("ano"),
	[]rune("ir"),
	[]rune("ar"),
}

func Join(a, b []rune) []rune {
	result := make([]rune, len(a)+len(b))
	bp := copy(result, a)
	copy(result[bp:], b)
	return result
}

func isVowel(r rune) bool {
	switch r {
	case 'a', 'e', 'i', 'o', 'u', 'à', 'è', 'ì', 'ò', 'ù':
		return true
	}
	return false
}

// findRegion returns the start of the first region in word; a region starts from a
// non-vowel charcter followed by a vowel.
func findRegion(word []rune, start int) int {
	l := len(word)
	var oldr rune

	for i, runeValue := range word {
		if i == 0 || i < start {
			oldr = runeValue
			continue
		}
		if !isVowel(runeValue) && isVowel(oldr) && i < l {
			return i + 1
		}
		oldr = runeValue
	}

	return l
}

// findR12 finds the R1 and R2 region for a word.
//   - R1 is the region after the first non-vowel following a vowel, or is the null
//     region at the end of the word if there is no such non-vowel.
//   - R2 is the region after the first non-vowel following a vowel in R1, or is the
//     null region at the end of the word if there is no such non-vowel.
func findR12(word []rune) (int, int) {
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
func findRV(word []rune) int {
	l := len(word)

	// when word[1] cannot be found
	if l <= 2 {
		return l
	}
	ch1 := word[0]
	ch2 := word[1]

	if !isVowel(ch2) {
		for i, r := range word {
			if i < 1 {
				continue
			}
			if isVowel(r) {
				if i+1 < l {
					return i + 1
				} else {
					return l
				}
			}
		}
		return l
	} else if isVowel(ch1) && isVowel(ch2) {
		for i, r := range word {
			if i < 2 {
				continue
			}
			if !isVowel(r) {
				return i + 1
			}
		}
		return l
	}

	// "and otherwise (consonant-vowel case) RV is the region after the third letter."
	return 3
}

// replaceInRegion search for suffixes inside the region defined by rX inside word
// and replaces the suffix with repl; returns the modified word and true when it
// was modified.
func replaceInRegion(word []rune, suffixes [][]rune, repl []rune, rX int) ([]rune, bool) {
	var p int

	for _, s := range suffixes {
		if hasSuffix(word, s) {
			p = lastIndex(word, s)
			if p >= rX {
				return Join(word[0:p], repl), true
			}
		}
	}
	return word, false
}

// Step 0) Attached pronoun
func step0(word []rune) []rune {
	rv := findRV(word)

	for _, suff2 := range s0suff {
		var p int

		for _, suff1 := range [][]rune{[]rune("ando"), []rune("endo")} {
			suffj := Join(suff1, suff2)
			if hasSuffix(word, suffj) {
				p = lastIndex(word, suffj)
				if p != -1 && p >= rv {
					return word[0 : p+len(suff1)]
				}
			}
		}

		for _, suff1 := range [][]rune{[]rune("ar"), []rune("er"), []rune("ir")} {
			suffj := Join(suff1, suff2)
			if hasSuffix(word, suffj) {
				p = lastIndex(word, suffj)
				if p != -1 && p >= rv {
					return Join(word[0:p+len(suff1)], []rune{'e'})
				}
			}
		}
	}

	return word
}

// Step 1) Standard suffix removal:
// search for the longest among the following suffixes, and perform the action indicated.
func step1(word []rune) []rune {
	var p int
	r1, r2 := findR12(word)
	rv := findRV(word)

	if hasSuffix(word, []rune("amente")) {
		p := lastIndex(word, []rune("amente"))
		if p >= r1 {
			word = word[0:p]

			if hasSuffix(word, []rune("iv")) {
				p = lastIndex(word, []rune("iv"))
				if p >= r2 {
					word = word[0:p]

					if hasSuffix(word, []rune("at")) {
						p = lastIndex(word, []rune("at"))
						if p >= r2 {
							word = word[0:p]
						}
					}
				}
			} else {
				word, _ = replaceInRegion(word, [][]rune{[]rune("os"), []rune("ic"), []rune("abil")}, []rune(""), r2)
			}
		}
		return word
	}

	for _, s := range step1suffs {
		if hasSuffix(word, s) {
			p = lastIndex(word, s)
			if p >= r2 {
				return word[0:p]
			} else {
				break
			}
		}
	}

	for _, s := range [][]rune{[]rune("azione"), []rune("azioni"), []rune("atore"), []rune("atori")} {
		if hasSuffix(word, s) {
			p = lastIndex(word, s)
			if p >= r2 {
				word = word[0:p]

				if hasSuffix(word, []rune("ic")) {
					p = lastIndex(word, []rune("ic"))
					if p >= r2 {
						return word[0:p]
					}
				}
				return word
			}
		}
	}

	var repl bool
	if word, repl = replaceInRegion(word, [][]rune{[]rune("logia"), []rune("logie")}, []rune("log"), r2); repl {
		return word
	}

	if word, repl = replaceInRegion(word, [][]rune{[]rune("uzione"), []rune("uzioni"), []rune("usione"), []rune("usioni")}, []rune("u"), r2); repl {
		return word
	}

	if word, repl = replaceInRegion(word, [][]rune{[]rune("enza"), []rune("enze")}, []rune("ente"), r2); repl {
		return word
	}

	if word, repl = replaceInRegion(word, [][]rune{[]rune("amento"), []rune("amenti"), []rune("imento"), []rune("imenti")}, []rune(""), rv); repl {
		return word
	}

	if hasSuffix(word, []rune("ità")) {
		p = lastIndex(word, []rune("ità"))
		if p >= r2 {
			word = word[0:p]
		}
		word, _ = replaceInRegion(word, [][]rune{[]rune("abil"), []rune("ic"), []rune("iv")}, []rune(""), r2)
		return word
	}

	for _, s := range [][]rune{[]rune("ivo"), []rune("ivi"), []rune("iva"), []rune("ive")} {
		if hasSuffix(word, s) {
			p = lastIndex(word, s)
			if p >= r2 {
				word = word[0:p]

				if hasSuffix(word, []rune("at")) {
					p = lastIndex(word, []rune("at"))
					if p >= r2 {
						word = word[0:p]

						if hasSuffix(word, []rune("ic")) {
							p = lastIndex(word, []rune("ic"))
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
func step2(word []rune) []rune {
	rv := findRV(word)

	for _, s := range verbsuff {
		if hasSuffix(word[rv:], s) {
			p := lastIndex(word, s)
			return word[0:p]
		}
	}
	return word
}

func step3a(word []rune) []rune {
	rv := findRV(word)

	for _, s := range [][]rune{[]rune("a"), []rune("e"), []rune("i"), []rune("o"), []rune("à"), []rune("è"), []rune("ì"), []rune("ò")} {
		if !hasSuffix(word, []rune(s)) {
			continue
		}
		p := lastIndex(word, []rune(s))
		if p >= rv {
			word = word[0:p]

			if hasSuffix(word, []rune("i")) {
				pp := lastIndex(word, []rune("i"))
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

func step3b(word []rune) []rune {
	rv := findRV(word)
	var p int

	if hasSuffix(word, []rune("ch")) {
		p = lastIndex(word, []rune("ch"))
		if p >= rv {
			return Join(word[0:p], []rune("c"))
		}
	}

	if hasSuffix(word, []rune("gh")) {
		p = lastIndex(word, []rune("gh"))
		if p >= rv {
			return Join(word[0:p], []rune("g"))
		}
	}

	return word
}

func restoreString(word []rune) []rune {
	for i, r := range word {
		switch r {
		case 'I':
			word[i] = 'i'
		case 'U':
			word[i] = 'u'
		}
	}
	return word
}

// prepareWord returns a string "prepared" for stemming; it's the first step.
func prepareWord(word []rune) []rune {
	for i, r := range word {
		switch r {
		case 'á':
			word[i] = 'à'
		case 'é':
			word[i] = 'è'
		case 'í':
			word[i] = 'ì'
		case 'ó':
			word[i] = 'ò'
		case 'ú':
			word[i] = 'ù'
		}
	}
	word = replace(word, []rune("qu"), []rune("qU"), -1)

	var oldr rune
	var newword []rune
	for i, r := range word {
		if i == 0 {
			oldr = r
			newword = append(newword, r)
			continue
		}
		if (r == 'i' || r == 'u') && isVowel(oldr) && i+1 < len(word) {
			nr := word[i+1]
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
	return newword
}

// Stem returns the italian stemming for the word s, lowercased.
func Stem(s []rune) []rune {
	l := len(s)
	if l == 0 {
		return s
	}
	for i := 0; i < l; i++ {
		s[i] = unicode.ToLower(s[i])
	}

	result := StemWithoutLowerCasing(s)
	return result
}

// StemWithoutLowerCasing returns the italian stemming for the word s.
func StemWithoutLowerCasing(s []rune) []rune {
	word := prepareWord(s)
	word0 := step0(word)
	word1 := step1(word0)

	if equal(word0, word1) {
		word = step2(word1)
	} else {
		word = word1
	}

	word = step3a(word)
	word = step3b(word)
	word = restoreString(word)
	return word
}
