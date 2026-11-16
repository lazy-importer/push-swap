package pushswap

import (
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"testing"
)

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

// TestSortOptimalExhaustive vérifie que les petites piles sont triées avec le
// nombre minimal d'instructions. Les totaux attendus ont été calculés par une
// recherche exhaustive indépendante sur tout l'espace des états.
func TestSortOptimalExhaustive(t *testing.T) {
	if testing.Short() {
		t.Skip("recherche exhaustive désactivée en mode court")
	}

	totauxOptimaux := []int{0, 0, 1, 7, 69, 583, 5083}

	for n := 0; n <= 6; n++ {
		total := 0
		permuter(n, func(p []int) {
			s := New(append([]int(nil), p...))
			s.Sort()
			verifierSolution(t, p, s)
			total += len(s.Ops)
		})
		if total != totauxOptimaux[n] {
			t.Errorf("n=%d : total d'instructions = %d, attendu %d (solution non minimale)",
				n, total, totauxOptimaux[n])
		}
	}
}

// TestSortOptimalSeven vérifie la borne connue sur un échantillon de piles de
// sept éléments : toute pile se trie en 13 instructions au plus.
func TestSortOptimalSeven(t *testing.T) {
	if testing.Short() {
		t.Skip("recherche exacte désactivée en mode court")
	}

	rng := rand.New(rand.NewSource(7))
	for i := 0; i < 30; i++ {
		p := rng.Perm(7)
		s := New(append([]int(nil), p...))
		s.Sort()
		verifierSolution(t, p, s)
		if len(s.Ops) > 13 {
			t.Errorf("7 éléments triés en %d instructions, maximum optimal 13 (pile %v)",
				len(s.Ops), p)
		}
	}
}

func TestSortKnownExamples(t *testing.T) {
	tests := []struct {
		nombres []int
		maxOps  int
	}{
		{[]int{2, 1, 3, 6, 5, 8}, 7}, // exemple du sujet
		{[]int{5, 2, 9, 1, 7}, 6},    // optimal pour cette permutation
		{[]int{4, 67, 3, 87, 23}, 6}, // exemple d'utilisation du sujet
		{[]int{3, 2, 1, 0}, 4},       // quatre éléments
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

func TestSortGeneralStrategies(t *testing.T) {
	rng := rand.New(rand.NewSource(123))
	for _, n := range []int{8, 10, 25, 50, 100, 250} {
		for essai := 0; essai < 5; essai++ {
			p := rng.Perm(n)
			s := New(append([]int(nil), p...))
			s.Sort()
			verifierSolution(t, p, s)
		}
	}
}

func TestSortInsertionDirect(t *testing.T) {
	rng := rand.New(rand.NewSource(11))
	for _, n := range []int{8, 30, 100, 300} {
		p := rng.Perm(n)
		s := New(append([]int(nil), p...))
		s.sortInsertion()
		verifierSolution(t, p, s)
	}
}

func TestSortChunksDirect(t *testing.T) {
	rng := rand.New(rand.NewSource(12))
	for _, n := range []int{8, 30, 100, 500} {
		p := rng.Perm(n)
		s := New(append([]int(nil), p...))
		s.sortChunks()
		verifierSolution(t, p, s)
	}
}

// TestRanksGiveSameInstructions vérifie que trier les rangs produit la même
// suite d'instructions que trier les nombres d'origine.
func TestRanksGiveSameInstructions(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	for _, n := range []int{5, 6, 7, 20, 100} {
		bruts := make([]int, n)
		for i := range bruts {
			bruts[i] = rng.Intn(1_000_000) - 500_000
		}
		// Rendre les valeurs uniques.
		for i := range bruts {
			bruts[i] = bruts[i]*n + i
		}

		avecNombres := New(append([]int(nil), bruts...))
		avecNombres.Sort()

		avecRangs := New(Ranks(bruts))
		avecRangs.Sort()

		if !slices.Equal(avecNombres.Ops, avecRangs.Ops) {
			t.Errorf("n=%d : les rangs changent les instructions", n)
		}
	}
}

// TestOptimizationLevels vérifie que les stratégies générales restent sous des
// seuils raisonnables, mesurés sur des piles aléatoires reproductibles.
func TestOptimizationLevels(t *testing.T) {
	rng := rand.New(rand.NewSource(2024))

	for _, tc := range []struct {
		n          int
		essais     int
		moyenneMax int
	}{
		{100, 10, 800},
		{500, 3, 6500},
	} {
		total := 0
		for essai := 0; essai < tc.essais; essai++ {
			p := rng.Perm(tc.n)
			s := New(append([]int(nil), p...))
			s.Sort()
			verifierSolution(t, p, s)
			total += len(s.Ops)
		}
		moyenne := total / tc.essais
		if moyenne > tc.moyenneMax {
			t.Errorf("n=%d : %d instructions en moyenne, maximum accepté %d",
				tc.n, moyenne, tc.moyenneMax)
		}
	}
}

func BenchmarkSort(b *testing.B) {
	for _, n := range []int{5, 100, 500} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			rng := rand.New(rand.NewSource(int64(n)))
			for i := 0; i < b.N; i++ {
				s := New(rng.Perm(n))
				s.Sort()
			}
		})
	}
}
