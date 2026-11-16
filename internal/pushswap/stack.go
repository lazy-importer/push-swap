// Package pushswap contient les deux piles, les onze instructions
// Push-Swap, la validation des arguments et les stratégies de tri,
// partagés par push-swap et checker.
package pushswap

import (
	"fmt"
	"io"
)

// Stack regroupe les deux piles et les instructions appliquées.
// Le sommet de chaque pile est le premier élément : le premier entier
// reçu est le sommet de a.
type Stack struct {
	A   []int
	B   []int
	Ops []string
}

// New construit la pile a à partir de values, avec b vide.
func New(values []int) *Stack {
	a := make([]int, len(values))
	copy(a, values)
	return &Stack{A: a}
}

// log enregistre une instruction dans l'historique.
func (s *Stack) log(op string) {
	s.Ops = append(s.Ops, op)
}

// --- opérations élémentaires (sans enregistrement) -------------------------

func (s *Stack) swapA() {
	if len(s.A) >= 2 {
		s.A[0], s.A[1] = s.A[1], s.A[0]
	}
}

func (s *Stack) swapB() {
	if len(s.B) >= 2 {
		s.B[0], s.B[1] = s.B[1], s.B[0]
	}
}

func (s *Stack) rotateA() {
	if len(s.A) >= 2 {
		s.A = append(s.A[1:], s.A[0])
	}
}

func (s *Stack) rotateB() {
	if len(s.B) >= 2 {
		s.B = append(s.B[1:], s.B[0])
	}
}

func (s *Stack) reverseRotateA() {
	if len(s.A) >= 2 {
		dernier := len(s.A) - 1
		s.A = append([]int{s.A[dernier]}, s.A[:dernier]...)
	}
}

func (s *Stack) reverseRotateB() {
	if len(s.B) >= 2 {
		dernier := len(s.B) - 1
		s.B = append([]int{s.B[dernier]}, s.B[:dernier]...)
	}
}

func (s *Stack) pushA() {
	if len(s.B) > 0 {
		s.A = append([]int{s.B[0]}, s.A...)
		s.B = s.B[1:]
	}
}

func (s *Stack) pushB() {
	if len(s.A) > 0 {
		s.B = append([]int{s.A[0]}, s.B...)
		s.A = s.A[1:]
	}
}

// --- instructions publiques ------------------------------------------------

// Sa échange les deux premiers éléments de a.
func (s *Stack) Sa() { s.swapA(); s.log("sa") }

// Sb échange les deux premiers éléments de b.
func (s *Stack) Sb() { s.swapB(); s.log("sb") }

// Ss exécute sa et sb.
func (s *Stack) Ss() { s.swapA(); s.swapB(); s.log("ss") }

// Pa prend le premier de b et le pose sur a.
func (s *Stack) Pa() { s.pushA(); s.log("pa") }

// Pb prend le premier de a et le pose sur b.
func (s *Stack) Pb() { s.pushB(); s.log("pb") }

// Ra décale a vers le haut : le premier devient le dernier.
func (s *Stack) Ra() { s.rotateA(); s.log("ra") }

// Rb décale b vers le haut.
func (s *Stack) Rb() { s.rotateB(); s.log("rb") }

// Rr exécute ra et rb.
func (s *Stack) Rr() { s.rotateA(); s.rotateB(); s.log("rr") }

// Rra décale a vers le bas : le dernier devient le premier.
func (s *Stack) Rra() { s.reverseRotateA(); s.log("rra") }

// Rrb décale b vers le bas.
func (s *Stack) Rrb() { s.reverseRotateB(); s.log("rrb") }

// Rrr exécute rra et rrb.
func (s *Stack) Rrr() { s.reverseRotateA(); s.reverseRotateB(); s.log("rrr") }

// Apply exécute une instruction donnée sous forme de texte.
// Renvoie une erreur si l'instruction n'existe pas.
func (s *Stack) Apply(op string) error {
	switch op {
	case "sa":
		s.Sa()
	case "sb":
		s.Sb()
	case "ss":
		s.Ss()
	case "pa":
		s.Pa()
	case "pb":
		s.Pb()
	case "ra":
		s.Ra()
	case "rb":
		s.Rb()
	case "rr":
		s.Rr()
	case "rra":
		s.Rra()
	case "rrb":
		s.Rrb()
	case "rrr":
		s.Rrr()
	default:
		return ErrUnknownInstruction
	}
	return nil
}

// IsSorted dit si a est triée en ordre croissant.
func (s *Stack) IsSorted() bool {
	for i := 1; i < len(s.A); i++ {
		if s.A[i-1] > s.A[i] {
			return false
		}
	}
	return true
}

// Done dit si a est triée et b vide : c'est l'objectif du sujet.
func (s *Stack) Done() bool {
	return len(s.B) == 0 && s.IsSorted()
}

// WriteOps écrit les instructions enregistrées, une par ligne.
func (s *Stack) WriteOps(w io.Writer) error {
	for _, op := range s.Ops {
		if _, err := fmt.Fprintln(w, op); err != nil {
			return err
		}
	}
	return nil
}

// rotateATo amène l'élément situé à position au sommet de a,
// par le plus court chemin (ra ou rra).
func (s *Stack) rotateATo(position int) {
	if position <= 0 || position >= len(s.A) {
		return
	}
	if position < len(s.A)-position {
		for range position {
			s.Ra()
		}
		return
	}
	for range len(s.A) - position {
		s.Rra()
	}
}

// rotateBTo amène l'élément situé à position au sommet de b,
// par le plus court chemin (rb ou rrb).
func (s *Stack) rotateBTo(position int) {
	if position <= 0 || position >= len(s.B) {
		return
	}
	if position < len(s.B)-position {
		for range position {
			s.Rb()
		}
		return
	}
	for range len(s.B) - position {
		s.Rrb()
	}
}
