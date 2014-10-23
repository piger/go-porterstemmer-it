# Go Porter Stemmer for italian language

A native Go implementation of the Porter Stemmer Algorithm for the italian language.

For more informations see:

http://snowball.tartarus.org/algorithms/italian/stemmer.html

## Usage

A basic example:

``` go
package main

import (
	"fmt"
	"github.com/piger/go-porterstemmer-it"
)

func main() {
	word := "abbandonerò"
	stem := porterstemmer.StemString(word)
	fmt.Printf("The word %q has the stem %q\n", word, stem)
}
```
