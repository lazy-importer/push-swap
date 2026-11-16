package pushswap

import (
	"bytes"
	"slices"
	"testing"
)

func TestInstructions(t *testing.T) {
	tests := []struct {
		nom    string
		op     string
		a, b   []int
		wantA  []int
		wantB  []int
		record string
	}{
		{"sa", "sa", []int{1, 2, 3}, []int{9, 8}, []int{2, 1, 3}, []int{9, 8}, "sa"},
		{"sb", "sb", []int{1, 2}, []int{3, 4, 5}, []int{1, 2}, []int{4, 3, 5}, "sb"},
		{"ss", "ss", []int{1, 2}, []int{3, 4}, []int{2, 1}, []int{4, 3}, "ss"},
		{"pa", "pa", []int{1, 2}, []int{3, 4}, []int{3, 1, 2}, []int{4}, "pa"},
		{"pb", "pb", []int{1, 2}, []int{3, 4}, []int{2}, []int{1, 3, 4}, "pb"},
		{"ra", "ra", []int{1, 2, 3}, []int{9, 8}, []int{2, 3, 1}, []int{9, 8}, "ra"},
		{"rb", "rb", []int{1}, []int{3, 4, 5}, []int{1}, []int{4, 5, 3}, "rb"},
		{"rr", "rr", []int{1, 2, 3}, []int{4, 5, 6}, []int{2, 3, 1}, []int{5, 6, 4}, "rr"},
		{"rra", "rra", []int{1, 2, 3}, []int{9}, []int{3, 1, 2}, []int{9}, "rra"},
		{"rrb", "rrb", []int{1}, []int{3, 4, 5}, []int{1}, []int{5, 3, 4}, "rrb"},
		{"rrr", "rrr", []int{1, 2, 3}, []int{4, 5, 6}, []int{3, 1, 2}, []int{6, 4, 5}, "rrr"},
	}

	for _, tc := range tests {
		t.Run(tc.nom, func(t *testing.T) {
			s := &Stack{A: slices.Clone(tc.a), B: slices.Clone(tc.b)}
			if err := s.Apply(tc.op); err != nil {
				t.Fatalf("Apply(%q) : %v", tc.op, err)
			}
			if !slices.Equal(s.A, tc.wantA) {
				t.Errorf("a = %v, attendu %v", s.A, tc.wantA)
			}
			if !slices.Equal(s.B, tc.wantB) {
				t.Errorf("b = %v, attendu %v", s.B, tc.wantB)
			}
			if !slices.Equal(s.Ops, []string{tc.record}) {
				t.Errorf("opérations = %v, attendu [%s]", s.Ops, tc.record)
			}
		})
	}
}

func TestInstructionsOnSmallStacks(t *testing.T) {
	for _, op := range allInstructions {
		s := &Stack{A: []int{42}}
		if err := s.Apply(op); err != nil {
			t.Fatalf("Apply(%q) : %v", op, err)
		}
		if len(s.A)+len(s.B) != 1 {
			t.Fatalf("Apply(%q) a perdu un élément : a=%v b=%v", op, s.A, s.B)
		}
		if op != "pb" && !slices.Equal(s.A, []int{42}) {
			t.Fatalf("Apply(%q) a modifié une pile trop petite : a=%v", op, s.A)
		}
	}
}

func TestUnknownInstruction(t *testing.T) {
	s := New([]int{3, 1, 2})
	for _, op := range []string{"SA", "saa", " sa", "sa ", "", "rotate"} {
		if err := s.Apply(op); err == nil {
			t.Errorf("Apply(%q) aurait dû renvoyer une erreur", op)
		}
	}
	if !slices.Equal(s.A, []int{3, 1, 2}) {
		t.Fatalf("une instruction inconnue a modifié la pile : %v", s.A)
	}
}

func TestDone(t *testing.T) {
	tests := []struct {
		a, b []int
		want bool
	}{
		{[]int{1, 2, 3}, nil, true},
		{[]int{1, 1, 2}, nil, true},
		{nil, nil, true},
		{[]int{3, 2, 1}, nil, false},
		{[]int{1, 2, 3}, []int{4}, false},
		{[]int{2, 3}, []int{1}, false},
	}
	for _, tc := range tests {
		s := &Stack{A: tc.a, B: tc.b}
		if got := s.Done(); got != tc.want {
			t.Errorf("Done(a=%v, b=%v) = %v, attendu %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestWriteOps(t *testing.T) {
	s := New([]int{2, 1})
	s.Sa()
	s.Pb()
	s.Pa()

	var buffer bytes.Buffer
	if err := s.WriteOps(&buffer); err != nil {
		t.Fatalf("WriteOps : %v", err)
	}
	if got, want := buffer.String(), "sa\npb\npa\n"; got != want {
		t.Errorf("WriteOps = %q, attendu %q", got, want)
	}
}

func TestRotateHelpers(t *testing.T) {
	tests := []struct {
		position int
		wantA    []int
	}{
		{0, []int{1, 2, 3, 4, 5}},
		{2, []int{3, 4, 5, 1, 2}}, // two ra
		{4, []int{5, 1, 2, 3, 4}}, // one rra
	}
	for _, tc := range tests {
		s := &Stack{A: []int{1, 2, 3, 4, 5}}
		s.rotateATo(tc.position)
		if !slices.Equal(s.A, tc.wantA) {
			t.Errorf("rotateATo(%d) : a = %v, attendu %v", tc.position, s.A, tc.wantA)
		}
	}
}
