package pushswap

import (
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		nom     string
		nombres []int
		entree  string
		want    bool
		err     bool
	}{
		{"déjà trié sans instruction", []int{1, 2, 3}, "", true, false},
		{"non trié sans instruction", []int{3, 2, 1}, "", false, false},
		{"exemple du sujet", []int{3, 2, 1, 0}, "rra\npb\nsa\nrra\npa\n", true, false},
		{"séquence incomplète", []int{3, 2, 1, 0}, "rra\npb\nsa\n", false, false},
		{"b non vide", []int{1, 2}, "pb\n", false, false},
		{"combinaisons", []int{2, 1, 3}, "ss\n", true, false},
		{"rr", []int{3, 1, 2}, "rr\n", true, false},
		{"rrr", []int{2, 3, 1}, "rrr\n", true, false},
		{"retour à l'état initial", []int{3, 1, 2}, "ra\nrra\n", false, false},
		{"fins de ligne Windows", []int{2, 1}, "sa\r\n", true, false},
		{"instruction inconnue", []int{2, 1}, "rotate\n", false, true},
		{"instruction mal formatée", []int{2, 1}, " sa\n", false, true},
		{"espace final", []int{2, 1}, "sa \n", false, true},
		{"ligne vide", []int{2, 1}, "\n", false, true},
		{"majuscules refusées", []int{2, 1}, "SA\n", false, true},
	}

	for _, tc := range tests {
		t.Run(tc.nom, func(t *testing.T) {
			got, err := Check(tc.nombres, strings.NewReader(tc.entree))
			if tc.err {
				if err == nil {
					t.Fatalf("Check aurait dû renvoyer une erreur, obtenu %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Check : %v", err)
			}
			if got != tc.want {
				t.Errorf("Check = %v, attendu %v", got, tc.want)
			}
		})
	}
}

func TestCheckDoesNotModifyInput(t *testing.T) {
	nombres := []int{3, 1, 2}
	original := append([]int(nil), nombres...)
	if _, err := Check(nombres, strings.NewReader("sa\n")); err != nil {
		t.Fatalf("Check : %v", err)
	}
	for i := range nombres {
		if nombres[i] != original[i] {
			t.Fatalf("Check a modifié son entrée : %v", nombres)
		}
	}
}

func TestCheckEveryInstruction(t *testing.T) {
	// Chaque instruction, appliquée une fois, doit être acceptée par le
	// checker même sans effet sur une pile de deux éléments.
	for _, op := range []string{
		"sa", "sb", "ss", "pa", "pb",
		"ra", "rb", "rr", "rra", "rrb", "rrr",
	} {
		if _, err := Check([]int{1, 2}, strings.NewReader(op+"\n")); err != nil {
			t.Errorf("Check a refusé %q : %v", op, err)
		}
	}
}
