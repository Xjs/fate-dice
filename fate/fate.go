package fate

import (
	"math/rand"
	"time"
)

type diceValue struct {
	label string
	value int
}

var fateDiceValues = map[int]diceValue{
	0: diceValue{"-", -1},
	1: diceValue{"o", 0},
	2: diceValue{"+", 1},
}

// Seed seeds the math/rand random number generator with the current time.
func Seed() {
	rand.Seed(time.Now().UnixNano())
}

// Fate throws n fate dice, returning a result string showing the die faces
// and the sum of the die values.
func Fate(n int) (string, int) {
	resultString := ""
	result := 0
	for i := 0; i < n; i++ {
		v := fateDiceValues[rand.Intn(len(fateDiceValues))]
		resultString += v.label
		result += v.value
	}

	if result == n { // critical hit
		resultString = "🎯 " + resultString
	} else if result == -n { // critical miss
		resultString = "💩 " + resultString
	}

	return resultString, result
}
