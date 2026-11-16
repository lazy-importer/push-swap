package pushswap

import "math"

// Sort sorts stack a in ascending order (smallest at the top) and leaves
// stack b empty. Stack b is assumed to be empty when Sort is called.
//
// Up to maxOptimalSize elements, a breadth-first search returns a provably
// shortest sequence. Above that, a cheapest-insertion strategy is used up to
// a moderate size, then the chunk strategy for very large stacks.
func (s *Stack) Sort() {
	if s.Done() || len(s.A) < 2 {
		return
	}

	switch n := len(s.A); {
	case n <= maxOptimalSize:
		s.sortOptimal()
	case n <= 1000:
		s.sortInsertion()
	default:
		s.sortChunks()
	}
}

// sortThree sorts three elements in at most two instructions. Only the
// relative order of the three values matters.
func (s *Stack) sortThree() {
	if len(s.A) < 3 {
		return
	}

	haut, milieu, bas := s.A[0], s.A[1], s.A[2]

	switch {
	case haut < milieu && milieu < bas:
		// 1 2 3: already sorted.
	case haut > milieu && milieu < bas && haut < bas:
		// 2 1 3
		s.Sa()
	case haut > milieu && milieu < bas && haut > bas:
		// 3 1 2
		s.Ra()
	case haut < milieu && milieu > bas && haut > bas:
		// 2 3 1
		s.Rra()
	case haut < milieu && milieu > bas && haut < bas:
		// 1 3 2
		s.Sa()
		s.Ra()
	default:
		// 3 2 1
		s.Sa()
		s.Rra()
	}
}

// sortInsertion keeps the first three elements in a, sorts them, then
// inserts every element of b at the cheapest position.
//
// After the first three elements are sorted, a is always circularly sorted:
// reading it from any starting point wraps around in ascending order. The
// element of b to insert goes just before the smallest value of a that is
// larger than it, or before the minimum when it is the new maximum.
func (s *Stack) sortInsertion() {
	for len(s.A) > 3 {
		s.Pb()
	}
	s.sortThree()

	for len(s.B) > 0 {
		meilleur := -1
		meilleurCout := 0
		meilleureCible := 0

		for i, valeur := range s.B {
			cible := targetPosition(s.A, valeur)
			cout := combinedRotationCost(cible, len(s.A), i, len(s.B))
			if meilleur < 0 || cout < meilleurCout {
				meilleur = i
				meilleurCout = cout
				meilleureCible = cible
			}
		}

		s.rotateAAndBTo(meilleureCible, meilleur)
		s.Pa()
	}

	s.rotateATo(positionOfMinimum(s.A))
}

// targetPosition returns the position where valeur must be inserted in the
// circularly sorted stack a. It runs in O(log n).
func targetPosition(a []int, valeur int) int {
	n := len(a)
	if n == 0 {
		return 0
	}

	// Find the position of the minimum: the array is a sorted sequence
	// rotated around it.
	lo, hi := 0, n-1
	for lo < hi {
		milieu := (lo + hi) / 2
		if a[milieu] > a[hi] {
			lo = milieu + 1
		} else {
			hi = milieu
		}
	}
	pivot := lo

	// Binary search the first element greater than valeur, reading the array
	// in logical order starting at the minimum.
	lo, hi = 0, n
	for lo < hi {
		milieu := (lo + hi) / 2
		if a[(pivot+milieu)%n] > valeur {
			hi = milieu
		} else {
			lo = milieu + 1
		}
	}
	if lo == n {
		return pivot
	}
	return (pivot + lo) % n
}

// rotationPlan chooses the direction (up = ra/rb, down = rra/rrb) and the
// number of rotations needed to bring position to the top.
func rotationPlan(position, length int) (up bool, cost int) {
	if position <= length-position {
		return true, position
	}
	return false, length - position
}

// combinedRotationCost is the number of instructions needed to rotate both
// stacks. When both go in the same direction, rr or rrr does both at once,
// so the cost is the maximum instead of the sum.
func combinedRotationCost(aPosition, aLength, bPosition, bLength int) int {
	aUp, aCost := rotationPlan(aPosition, aLength)
	bUp, bCost := rotationPlan(bPosition, bLength)
	if aUp != bUp {
		return aCost + bCost
	}
	if aCost > bCost {
		return aCost
	}
	return bCost
}

// rotateAAndBTo brings both requested positions to the top, using rr/rrr
// when the rotation directions match.
func (s *Stack) rotateAAndBTo(aPosition, bPosition int) {
	aUp, aCost := rotationPlan(aPosition, len(s.A))
	bUp, bCost := rotationPlan(bPosition, len(s.B))

	if aUp != bUp {
		s.rotateATo(aPosition)
		s.rotateBTo(bPosition)
		return
	}

	commun := aCost
	if bCost < commun {
		commun = bCost
	}
	for i := 0; i < commun; i++ {
		if aUp {
			s.Rr()
		} else {
			s.Rrr()
		}
	}
	for i := commun; i < aCost; i++ {
		if aUp {
			s.Ra()
		} else {
			s.Rra()
		}
	}
	for i := commun; i < bCost; i++ {
		if bUp {
			s.Rb()
		} else {
			s.Rrb()
		}
	}
}

// sortChunks is the general strategy.
//
// Phase 1: the ranks are split into about sqrt(n) blocks. Blocks are pushed
// to b from the smallest ranks to the largest. Inside a block, the lower half
// is sent to the bottom of b with rb, so b ends up almost ordered with the
// largest values on top.
//
// Phase 2: the maximum of b is brought to the top by the shortest path and
// pushed back on a. Since the maximum is usually near the top, this phase is
// cheap, and a is built in ascending order.
func (s *Stack) sortChunks() {
	n := len(s.A)
	rankOf := make(map[int]int, n)
	for i, rang := range Ranks(s.A) {
		rankOf[s.A[i]] = rang
	}

	nbBlocs := int(math.Sqrt(float64(n)))
	if nbBlocs < 2 {
		nbBlocs = 2
	}
	taille := (n + nbBlocs - 1) / nbBlocs

	for lo := 0; lo < n; lo += taille {
		hi := lo + taille
		if hi > n {
			hi = n
		}
		milieu := (lo + hi) / 2

		for {
			position := -1
			meilleureDistance := len(s.A) + 1

			for i, valeur := range s.A {
				rang := rankOf[valeur]
				if rang < lo || rang >= hi {
					continue
				}
				distance := i
				if len(s.A)-i < distance {
					distance = len(s.A) - i
				}
				if distance < meilleureDistance {
					meilleureDistance = distance
					position = i
				}
			}

			if position < 0 {
				break
			}

			s.rotateATo(position)
			valeur := s.A[0]
			s.Pb()
			if rankOf[valeur] < milieu && len(s.B) > 1 {
				s.Rb()
			}
		}
	}

	for len(s.B) > 0 {
		position := positionOfMaximum(s.B)
		s.rotateBTo(position)
		s.Pa()
	}
}

func positionOfMinimum(pile []int) int {
	position := 0
	for i := 1; i < len(pile); i++ {
		if pile[i] < pile[position] {
			position = i
		}
	}
	return position
}

func positionOfMaximum(pile []int) int {
	position := 0
	for i := 1; i < len(pile); i++ {
		if pile[i] > pile[position] {
			position = i
		}
	}
	return position
}
