package movegen

import (
	"github.com/domino14/macondo/alphabet"
)

const (
	// TrivialCrossSet allows every possible letter. It is the default
	// state of a square.
	TrivialCrossSet = (1 << alphabet.MaxAlphabetSize) - 1
)

// A CrossSet is a bit mask of letters that are allowed on a square.
type CrossSet uint64

func (c CrossSet) allowed(letter alphabet.MachineLetter) bool {
	return c&(1<<uint8(letter)) != 0
}

func (c *CrossSet) set(letter alphabet.MachineLetter) {
	*c = *c | (1 << letter)
}

func (c *CrossSet) setAll() {
	*c = TrivialCrossSet
}

func (c *CrossSet) clear() {
	*c = 0
}