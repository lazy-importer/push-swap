package pushswap

func (s *Stack) Sort() {
	if s.Done() || len(s.A) < 2 {
		return
	}
	s.sortInsertion()
}

// sortThree trie trois éléments en deux instructions au plus.
// Seul l'ordre relatif des trois valeurs compte.
func (s *Stack) sortThree() {
	if len(s.A) < 3 {
		return
	}

	haut, milieu, bas := s.A[0], s.A[1], s.A[2]

	switch {
	case haut < milieu && milieu < bas:
		// 1 2 3 : déjà trié.
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

// sortInsertion garde les trois premiers éléments dans a, les trie, puis
// insère chaque élément de b au poste le moins cher.
//
// Une fois les trois premiers triés, a est toujours triée de façon
// circulaire : lue depuis n'importe quel point de départ, elle remonte
// en ordre croissant puis repart au début. L'élément de b s'insère juste
// avant la plus petite valeur de a qui lui est supérieure, ou avant le
// minimum quand il devient le nouveau maximum.
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

// targetPosition renvoie le poste où insérer valeur dans a, triée de
// façon circulaire. Complexité : O(log n).
func targetPosition(a []int, valeur int) int {
	n := len(a)
	if n == 0 {
		return 0
	}

	// Position du minimum : le tableau est une suite triée tournée autour.
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

	// Recherche du premier élément supérieur à valeur, en lisant le
	// tableau dans l'ordre logique qui part du minimum.
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

// rotationPlan choisit la direction (up = ra/rb, down = rra/rrb) et le
// nombre de rotations pour amener position au sommet.
func rotationPlan(position, length int) (up bool, cost int) {
	if position <= length-position {
		return true, position
	}
	return false, length - position
}

// combinedRotationCost compte les instructions pour tourner les deux
// piles. Quand elles tournent dans le même sens, rr ou rrr fait les deux
// d'un coup : le coût est le maximum, pas la somme.
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

// rotateAAndBTo amène les deux postes demandés au sommet, en utilisant
// rr ou rrr quand les directions de rotation sont les mêmes.
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
	for range commun {
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

// positionOfMinimum renvoie la position du plus petit élément de pile.
func positionOfMinimum(pile []int) int {
	position := 0
	for i := 1; i < len(pile); i++ {
		if pile[i] < pile[position] {
			position = i
		}
	}
	return position
}
