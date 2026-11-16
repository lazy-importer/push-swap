package pushswap

import (
	"io"
	"strings"
)

// Check exécute les instructions lues dans r sur une copie de numbers,
// puis dit si a est triée et b vide à la fin.
// Une instruction inconnue ou mal formatée renvoie une erreur.
func Check(numbers []int, r io.Reader) (bool, error) {
	s := New(numbers)

	// On lit tout, puis on découpe en lignes : une instruction par ligne.
	contenu, err := io.ReadAll(r)
	if err != nil {
		return false, err
	}
	lignes := strings.Split(string(contenu), "\n")

	// Le dernier élément est vide quand l'entrée se termine par \n.
	if len(lignes) > 0 && lignes[len(lignes)-1] == "" {
		lignes = lignes[:len(lignes)-1]
	}

	for _, ligne := range lignes {
		// Une ligne peut se terminer par \r : entrée venue de Windows.
		if err := s.Apply(strings.TrimSuffix(ligne, "\r")); err != nil {
			return false, err
		}
	}

	return s.Done(), nil
}
