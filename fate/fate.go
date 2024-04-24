package fate

import (
	"math/rand"
)

type diceValue struct {
	label string
	value int
}

var fateDiceValues = map[int]diceValue{
	0: {"-", -1},
	1: {"o", 0},
	2: {"+", 1},
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
