package pushswap

import (
	"errors"
	"fmt"
	"sort"
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

// Ranks remplace chaque nombre par sa position dans l'ordre croissant
// (0 pour le plus petit, n-1 pour le plus grand). Pas de doublons.
//
// Les rangs gardent l'ordre : trier les rangs produit exactement les
// mêmes instructions que trier les nombres, avec des valeurs toutes
// comprises entre 0 et n-1.
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
