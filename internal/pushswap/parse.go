package pushswap

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrUnknownInstruction est renvoyée quand une instruction ne fait pas
// partie des onze instructions Push-Swap.
var ErrUnknownInstruction = errors.New("unknown instruction")

// ParseArgs convertit les arguments du programme en entiers.
//
// Le sujet accepte un seul argument contenant toute la liste
// (push-swap "2 1 3 6 5 8") aussi bien qu'un argument par nombre
// (push-swap 2 1 3). Chaque argument est donc découpé sur les espaces.
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

// CheckDuplicates renvoie une erreur dès qu'un nombre apparaît deux fois.
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
