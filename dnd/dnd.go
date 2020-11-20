// package dnd implements dice throws with the D notation (see
// e. g. https://nethackwiki.com/wiki/D_notation)
package dnd

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
)

// A Throw represents a throw of Dice dice with Faces faces each.
type Throw struct {
	Dice  int
	Faces int
}

// Parse parses a string of the form XdY, where X and Y are integer numbers greater than zero,
// or X may be empty.
func Parse(s string) (Throw, error) {
	var result Throw
	words := strings.Split(s, "d")
	if len(words) != 2 {
		return Throw{}, errors.New("must contain exactly one d")
	}

	faces, err := strconv.Atoi(words[1])
	if err != nil {
		return Throw{}, err
	} else if faces < 1 {
		return Throw{}, errors.New("number of faces must be > 0")
	}
	result.Faces = faces

	if words[0] == "" {
		result.Dice = 1
		return result, nil
	}

	dice, err := strconv.Atoi(words[0])
	if err != nil {
		return Throw{}, err
	} else if dice < 1 {
		return Throw{}, errors.New("number of dice must be > 0")
	}
	result.Dice = dice

	return result, nil
}

// Throw throws with the number of dice with faces as specified in t.
// It uses math/rand.
func (t Throw) Throw() int {
	var result int
	for i := 0; i < t.Dice; i++ {
		result += (1 + rand.Intn(t.Faces))
	}
	return result
}

func (t Throw) Emoji() string {
	var result string
	for i := 0; i < t.Dice; i++ {
		result += "🎲"
	}
	return result
}

// D is a convenience function that Parses the string and if successful, Throws the dice.
func D(s string) (int, error) {
	t, err := Parse(s)
	if err != nil {
		return 0, err
	}
	return t.Throw(), nil
}
