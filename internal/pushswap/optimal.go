package pushswap

// maxOptimalSize est la plus grande taille de pile pour laquelle on
// cherche la solution exacte. L'espace d'états reste assez petit pour
// être exploré en largeur, ce qui garantit le nombre minimal
// d'instructions.
const maxOptimalSize = 7

var allInstructions = [...]string{
	"sa", "sb", "ss", "pa", "pb",
	"ra", "rb", "rr",
	"rra", "rrb", "rrr",
}

// searchNode est un état visité par le parcours en largeur.
type searchNode struct {
	a, b   []int
	parent int
	op     string
}

// sortOptimal trouve une plus courte séquence d'instructions par un
// parcours en largeur de tous les états des deux piles, puis la rejoue.
//
// La recherche se fait sur les rangs : un état se code de façon compacte
// et les rangs gardent l'ordre, donc la séquence est valable pour les
// valeurs réelles aussi.
func (s *Stack) sortOptimal() {
	if len(s.A) == 0 {
		return
	}

	start := append([]int(nil), Ranks(s.A)...)
	nodes := []searchNode{{a: start, parent: -1}}
	vus := map[string]struct{}{stateKey(start, nil): {}}

	for head := 0; head < len(nodes); head++ {
		current := nodes[head]
		if len(current.b) == 0 && strictlySorted(current.a) {
			replayOptimal(s, nodes, head)
			return
		}

		for _, op := range allInstructions {
			nextA, nextB := applyRanked(current.a, current.b, op)
			cle := stateKey(nextA, nextB)
			if _, dejaVu := vus[cle]; dejaVu {
				continue
			}
			vus[cle] = struct{}{}
			nodes = append(nodes, searchNode{
				a:      nextA,
				b:      nextB,
				parent: head,
				op:     op,
			})
		}
	}
}

// replayOptimal rejoue les instructions trouvées, de la première à la
// dernière, et les enregistre dans la pile.
func replayOptimal(s *Stack, nodes []searchNode, last int) {
	var instructions []string
	for i := last; nodes[i].parent >= 0; i = nodes[i].parent {
		instructions = append(instructions, nodes[i].op)
	}
	for i := len(instructions) - 1; i >= 0; i-- {
		_ = s.Apply(instructions[i])
	}
}

// stateKey code un état sous forme de chaîne. Les rangs sont plus petits
// que 256 pour toutes les tailles traitées par sortOptimal ; 0xFF sépare
// les deux piles.
func stateKey(a, b []int) string {
	buffer := make([]byte, 0, len(a)+len(b)+1)
	for _, valeur := range a {
		buffer = append(buffer, byte(valeur))
	}
	buffer = append(buffer, 0xFF)
	for _, valeur := range b {
		buffer = append(buffer, byte(valeur))
	}
	return string(buffer)
}

// strictlySorted dit si a est triée strictement (aucun doublon possible
// ici : les rangs sont tous différents).
func strictlySorted(a []int) bool {
	for i := 1; i < len(a); i++ {
		if a[i-1] >= a[i] {
			return false
		}
	}
	return true
}

// applyRanked renvoie le nouvel état après l'instruction op, sans jamais
// modifier ses arguments.
func applyRanked(a, b []int, op string) ([]int, []int) {
	nextA := append([]int(nil), a...)
	nextB := append([]int(nil), b...)

	swap := func(pile []int) {
		if len(pile) >= 2 {
			pile[0], pile[1] = pile[1], pile[0]
		}
	}
	rotateUp := func(pile []int) []int {
		if len(pile) >= 2 {
			return append(pile[1:], pile[0])
		}
		return pile
	}
	rotateDown := func(pile []int) []int {
		if len(pile) >= 2 {
			dernier := len(pile) - 1
			return append([]int{pile[dernier]}, pile[:dernier]...)
		}
		return pile
	}

	switch op {
	case "sa":
		swap(nextA)
	case "sb":
		swap(nextB)
	case "ss":
		swap(nextA)
		swap(nextB)
	case "pa":
		if len(nextB) > 0 {
			nextA = append([]int{nextB[0]}, nextA...)
			nextB = nextB[1:]
		}
	case "pb":
		if len(nextA) > 0 {
			nextB = append([]int{nextA[0]}, nextB...)
			nextA = nextA[1:]
		}
	case "ra":
		nextA = rotateUp(nextA)
	case "rb":
		nextB = rotateUp(nextB)
	case "rr":
		nextA = rotateUp(nextA)
		nextB = rotateUp(nextB)
	case "rra":
		nextA = rotateDown(nextA)
	case "rrb":
		nextB = rotateDown(nextB)
	case "rrr":
		nextA = rotateDown(nextA)
		nextB = rotateDown(nextB)
	}

	return nextA, nextB
}
