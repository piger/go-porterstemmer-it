package porterstemmer

// A bunch of functions copied from the Go standard library:
// https://golang.org/src/bytes/bytes.go

func equal(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i, c := range a {
		if c != b[i] {
			return false
		}
	}
	return true
}

func hasSuffix(s, suffix []rune) bool {
	return len(s) >= len(suffix) && equal(s[len(s)-len(suffix):], suffix)
}

func lastIndex(s, sep []rune) (result int) {
	n := len(sep)
	if n == 0 {
		return len(s)
	}
	c := sep[0]
	for i := len(s) - n; i >= 0; i-- {
		if s[i] == c && (n == 1 || equal(s[i:i+n], sep)) {
			return i
		}
	}
	return -1
}

func indexRune(s []rune, c rune) int {
	for i, r := range s {
		if r == c {
			return i
		}
	}
	return -1
}

func index(s, sep []rune) int {
	n := len(sep)
	if n == 0 {
		return 0
	}
	if n > len(s) {
		return -1
	}
	c := sep[0]
	if n == 1 {
		return indexRune(s, c)
	}
	i := 0
	t := s[:len(s)-n+1]
	for i < len(t) {
		if t[i] != c {
			o := indexRune(t[i:], c)
			if o < 0 {
				break
			}
			i += o
		}
		if equal(s[i:i+n], sep) {
			return i
		}
		i++
	}
	return -1
}

func count(s, sep []rune) int {
	n := len(sep)
	if n == 0 {
		return len(s) + 1
	}
	if n > len(s) {
		return 0
	}
	count := 0
	c := sep[0]
	i := 0
	t := s[:len(s)-n+1]
	for i < len(t) {
		if t[i] != c {
			o := indexRune(t[i:], c)
			if o < 0 {
				break
			}
			i += o
		}
		if n == 1 || equal(s[i:i+n], sep) {
			count++
			i += n
			continue
		}
		i++
	}
	return count
}

func replace(s, old, new []rune, n int) []rune {
	m := 0
	if n != 0 {
		m = count(s, old)
	}
	if m == 0 {
		// XXX why []rune(nil)? they return: append([]byte(nil), s...)
		return append([]rune(nil), s...)
	}
	if n < 0 || m < n {
		n = m
	}

	t := make([]rune, len(s)+n*(len(new)-len(old)))
	w := 0
	start := 0
	for i := 0; i < n; i++ {
		j := start
		if len(old) == 0 {
			if i > 0 {
				j++
			}
		} else {
			j += index(s[start:], old)
		}
		w += copy(t[w:], s[start:j])
		w += copy(t[w:], new)
		start = j + len(old)

	}
	w += copy(t[w:], s[start:])
	return t[0:w]
}
