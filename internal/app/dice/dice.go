package dice

import (
	"math/rand"
	"regexp"
	"strconv"
)

type DiceRoll struct {
	rolls int
	sides int
}

// Parses `ndn` format string (e.g. `2d6`)
func parseToDiceRoll(s string) *DiceRoll {
	re := regexp.MustCompile(`^(\d+)d(\d+)$`)
	match := re.FindStringSubmatch(s)
	if match == nil {
		return nil
	}

	rolls, err := strconv.Atoi(match[1])
	if err != nil {
		return nil
	}

	sides, err := strconv.Atoi(match[2])
	if err != nil {
		return nil
	}

	return &DiceRoll{
		rolls: rolls,
		sides: sides,
	}
}

func (d DiceRoll) Roll() int {
	result := 0
	for i := 0; i < d.rolls; i++ {
		result += rand.Intn(d.sides) + 1
	}

	return result
}
