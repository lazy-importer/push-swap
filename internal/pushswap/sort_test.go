package pushswap

import (
	"math/rand/v2"
	"strings"
	"testing"
)

func verifierSolution(t *testing.T, nombres []int, s *Stack) {
	t.Helper()

	if !s.Done() {
		t.Fatalf("le tri a échoué pour %v : a=%v b=%v", nombres, s.A, s.B)
	}

	entree := ""
	if len(s.Ops) > 0 {
		entree = strings.Join(s.Ops, "\n") + "\n"
	}
	ok, err := Check(nombres, strings.NewReader(entree))
	if err != nil {
		t.Fatalf("rejeu impossible pour %v : %v", nombres, err)
	}
	if !ok {
		t.Fatalf("le rejeu des opérations ne trie pas %v (ops=%v)", nombres, s.Ops)
	}
}

// permuter appelle fn sur toutes les permutations de 0..n-1.
func permuter(n int, fn func([]int)) {
	valeurs := make([]int, n)
	for i := range valeurs {
		valeurs[i] = i
	}
	var generer func(int)
	generer = func(k int) {
		if k == n {
			fn(valeurs)
			return
		}
		for i := k; i < n; i++ {
			valeurs[k], valeurs[i] = valeurs[i], valeurs[k]
			generer(k + 1)
			valeurs[k], valeurs[i] = valeurs[i], valeurs[k]
		}
	}
	generer(0)
}

// TestSortThree vérifie les six arrangements de trois éléments.
func TestSortThree(t *testing.T) {
	if testing.Short() {
		t.Skip("recherche exhaustive désactivée en mode court")
	}

	// Total attendu sur les six permutations de trois éléments : 0 + 2 +
	// 1 + 1 + 1 + 2 = 7 instructions.
	totaux := []int{0, 0, 1, 7}

	for n := 0; n <= 3; n++ {
		total := 0
		permuter(n, func(p []int) {
			s := New(append([]int(nil), p...))
			s.Sort()
			verifierSolution(t, p, s)
			total += len(s.Ops)
		})
		if total != totaux[n] {
			t.Errorf("n=%d : total d'instructions = %d, attendu %d", n, total, totaux[n])
		}
	}
}

func TestSortKnownExamples(t *testing.T) {
	tests := []struct {
		nombres []int
		maxOps  int
	}{
		{[]int{2, 1, 3, 6, 5, 8}, 9}, // exemple du sujet
		{[]int{4, 67, 3, 87, 23}, 9}, // exemple d'utilisation du sujet
		{[]int{3, 2, 1, 0}, 5},       // quatre éléments
		{[]int{1, 2, 3, 4, 5}, 0},    // déjà trié
		{[]int{1}, 0},                // un seul élément
		{[]int{2, 1}, 1},             // deux éléments
	}

	for _, tc := range tests {
		s := New(append([]int(nil), tc.nombres...))
		s.Sort()
		verifierSolution(t, tc.nombres, s)
		if len(s.Ops) > tc.maxOps {
			t.Errorf("%v : %d instructions, au plus %d attendues",
				tc.nombres, len(s.Ops), tc.maxOps)
		}
	}
}

func TestSortAlreadySorted(t *testing.T) {
	for n := 0; n <= 50; n++ {
		nombres := make([]int, n)
		for i := range nombres {
			nombres[i] = i
		}
		s := New(nombres)
		s.Sort()
		if len(s.Ops) != 0 {
			t.Fatalf("une pile déjà triée a produit %d instructions : %v", len(s.Ops), s.Ops)
		}
	}
}

// TestSortRandom vérifie le tri sur des piles aléatoires, y compris plus
// grandes que la stratégie d'insertion seule ne le couvre confortablement.
func TestSortRandom(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 123))
	for _, n := range []int{8, 10, 25, 50, 100, 250} {
		for essai := 0; essai < 5; essai++ {
			p := rng.Perm(n)
			s := New(append([]int(nil), p...))
			s.Sort()
			verifierSolution(t, p, s)
		}
	}
}
