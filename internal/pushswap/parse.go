package pushswap

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ErrUnknownInstruction is returned when an instruction is not one of the
// eleven Push-Swap instructions.
var ErrUnknownInstruction = errors.New("unknown instruction")

// ParseArgs converts the program arguments into integers.
//
// The subject allows a single argument containing the whole list
// (push-swap "2 1 3 6 5 8") as well as one argument per number
// (push-swap 2 1 3). Every argument is therefore split on spaces.
func ParseArgs(args []string) ([]int, error) {
	if len(args) == 0 {
		return nil, nil
	}

	nombres := make([]int, 0, len(args))
	for _, arg := range args {
		if strings.TrimSpace(arg) == "" {
			return nil, fmt.Errorf("%q is not an integer", arg)
		}
		for _, champ := range strings.Fields(arg) {
			n, err := strconv.Atoi(champ)
			if err != nil {
				return nil, fmt.Errorf("%q is not an integer", champ)
			}
			nombres = append(nombres, n)
		}
	}
	return nombres, nil
}

// CheckDuplicates returns an error as soon as a number appears twice.
func CheckDuplicates(nombres []int) error {
	vus := make(map[int]struct{}, len(nombres))
	for _, n := range nombres {
		if _, ok := vus[n]; ok {
			return fmt.Errorf("the number %d appears more than once", n)
		}
		vus[n] = struct{}{}
	}
	return nil
}

// Ranks replaces every number by its position in ascending order
// (0 for the smallest, n-1 for the largest). Duplicates are not allowed.
//
// Ranks keep the order: sorting the ranks produces exactly the same
// instructions as sorting the original numbers, while every rank stays
// between 0 and n-1.
func Ranks(nombres []int) []int {
	tries := append([]int(nil), nombres...)
	sort.Ints(tries)

	rangDe := make(map[int]int, len(tries))
	for rang, valeur := range tries {
		rangDe[valeur] = rang
	}

	rangs := make([]int, len(nombres))
	for i, valeur := range nombres {
		rangs[i] = rangDe[valeur]
	}
	return rangs
}
