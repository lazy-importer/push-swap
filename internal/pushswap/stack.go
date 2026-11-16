// Package pushswap implements the two stacks, the eleven Push-Swap
// instructions, argument validation, the sorting strategies and the checker
// logic shared by the push-swap and checker executables.
package pushswap

import (
	"bufio"
	"io"
)

// Stack holds both stacks and every instruction applied so far.
//
// The top of each stack is the first element, as required by the subject:
// the first integer received is the top of stack a.
type Stack struct {
	A   []int
	B   []int
	Ops []string
}

// New builds a stack a from values, with an empty stack b.
func New(values []int) *Stack {
	a := make([]int, len(values))
	copy(a, values)
	return &Stack{A: a}
}

func (s *Stack) log(op string) {
	s.Ops = append(s.Ops, op)
}

// --- elementary helpers (no logging) ---------------------------------------

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
		last := len(s.A) - 1
		s.A = append([]int{s.A[last]}, s.A[:last]...)
	}
}

func (s *Stack) reverseRotateB() {
	if len(s.B) >= 2 {
		last := len(s.B) - 1
		s.B = append([]int{s.B[last]}, s.B[:last]...)
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

// --- public instructions ---------------------------------------------------

// Sa swaps the first two elements of a.
func (s *Stack) Sa() { s.swapA(); s.log("sa") }

// Sb swaps the first two elements of b.
func (s *Stack) Sb() { s.swapB(); s.log("sb") }

// Ss executes sa and sb.
func (s *Stack) Ss() { s.swapA(); s.swapB(); s.log("ss") }

// Pa moves the top of b to the top of a.
func (s *Stack) Pa() { s.pushA(); s.log("pa") }

// Pb moves the top of a to the top of b.
func (s *Stack) Pb() { s.pushB(); s.log("pb") }

// Ra shifts a up: the first element becomes the last one.
func (s *Stack) Ra() { s.rotateA(); s.log("ra") }

// Rb shifts b up.
func (s *Stack) Rb() { s.rotateB(); s.log("rb") }

// Rr executes ra and rb.
func (s *Stack) Rr() { s.rotateA(); s.rotateB(); s.log("rr") }

// Rra shifts a down: the last element becomes the first one.
func (s *Stack) Rra() { s.reverseRotateA(); s.log("rra") }

// Rrb shifts b down.
func (s *Stack) Rrb() { s.reverseRotateB(); s.log("rrb") }

// Rrr executes rra and rrb.
func (s *Stack) Rrr() { s.reverseRotateA(); s.reverseRotateB(); s.log("rrr") }

// Apply executes one instruction given as text. It returns an error when the
// instruction does not exist.
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

// IsSorted reports whether a is sorted in ascending order, top first.
func (s *Stack) IsSorted() bool {
	for i := 1; i < len(s.A); i++ {
		if s.A[i-1] > s.A[i] {
			return false
		}
	}
	return true
}

// Done reports whether a is sorted and b is empty.
func (s *Stack) Done() bool {
	return len(s.B) == 0 && s.IsSorted()
}

// WriteOps writes the recorded instructions, one per line, to w.
func (s *Stack) WriteOps(w io.Writer) error {
	buffered := bufio.NewWriter(w)
	for _, op := range s.Ops {
		if _, err := buffered.WriteString(op); err != nil {
			return err
		}
		if err := buffered.WriteByte('\n'); err != nil {
			return err
		}
	}
	return buffered.Flush()
}

// rotateATo brings the element at position to the top of a using the shortest
// path (ra or rra).
func (s *Stack) rotateATo(position int) {
	if position <= 0 || position >= len(s.A) {
		return
	}
	if position < len(s.A)-position {
		for i := 0; i < position; i++ {
			s.Ra()
		}
		return
	}
	for i := 0; i < len(s.A)-position; i++ {
		s.Rra()
	}
}

// rotateBTo brings the element at position to the top of b using the shortest
// path (rb or rrb).
func (s *Stack) rotateBTo(position int) {
	if position <= 0 || position >= len(s.B) {
		return
	}
	if position < len(s.B)-position {
		for i := 0; i < position; i++ {
			s.Rb()
		}
		return
	}
	for i := 0; i < len(s.B)-position; i++ {
		s.Rrb()
	}
}
