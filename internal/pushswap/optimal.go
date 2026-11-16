package pushswap

// maxOptimalSize is the largest stack for which the exact shortest solution
// is searched. The whole state space is small enough to explore with a
// breadth-first search, which guarantees a minimal number of instructions.
const maxOptimalSize = 7

var allInstructions = [...]string{
	"sa", "sb", "ss", "pa", "pb",
	"ra", "rb", "rr",
	"rra", "rrb", "rrr",
}

// searchNode is one state visited by the breadth-first search.
type searchNode struct {
	a, b   []int
	parent int
	op     string
}

// sortOptimal finds a shortest sequence of instructions with a breadth-first
// search over all the states of the two stacks, then replays it.
//
// The search is done on ranks so that a state can be encoded compactly.
// Ranks preserve the order, so the sequence is valid for the real values too.
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
			if _, already := vus[cle]; already {
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

// replayOptimal applies the instructions found by the search, from the first
// move to the last, and records them in the stack.
func replayOptimal(s *Stack, nodes []searchNode, last int) {
	var instructions []string
	for i := last; nodes[i].parent >= 0; i = nodes[i].parent {
		instructions = append(instructions, nodes[i].op)
	}
	for i := len(instructions) - 1; i >= 0; i-- {
		_ = s.Apply(instructions[i])
	}
}

// stateKey encodes a state as a string. Ranks are smaller than 256 for every
// size handled by sortOptimal; 0xFF separates the two stacks.
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

func strictlySorted(a []int) bool {
	for i := 1; i < len(a); i++ {
		if a[i-1] >= a[i] {
			return false
		}
	}
	return true
}

// applyRanked returns a new state after applying op. It never modifies its
// arguments.
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
			last := len(pile) - 1
			return append([]int{pile[last]}, pile[:last]...)
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
