package pushswap

import (
	"slices"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		nom  string
		args []string
		want []int
		err  bool
	}{
		{"aucun argument", nil, nil, false},
		{"un seul argument groupé", []string{"2 1 3 6 5 8"}, []int{2, 1, 3, 6, 5, 8}, false},
		{"un argument par nombre", []string{"5", "2", "9"}, []int{5, 2, 9}, false},
		{"mélange", []string{"5 2", "9"}, []int{5, 2, 9}, false},
		{"négatifs", []string{"-3", "10", "0", "7", "-1", "2"}, []int{-3, 10, 0, 7, -1, 2}, false},
		{"plus explicite", []string{"+4"}, []int{4}, false},
		{"zéro non significatif", []string{"05", "5"}, []int{5, 5}, false},
		{"décimal refusé", []string{"3.5"}, nil, true},
		{"texte refusé", []string{"0 one 2 3"}, nil, true},
		{"argument vide refusé", []string{""}, nil, true},
		{"espaces seuls refusés", []string{"   "}, nil, true},
		{"débordement refusé", []string{"999999999999999999999999"}, nil, true},
	}

	for _, tc := range tests {
		t.Run(tc.nom, func(t *testing.T) {
			got, err := ParseArgs(tc.args)
			if tc.err {
				if err == nil {
					t.Fatalf("ParseArgs(%q) aurait dû échouer, obtenu %v", tc.args, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseArgs(%q) : %v", tc.args, err)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("ParseArgs(%q) = %v, attendu %v", tc.args, got, tc.want)
			}
		})
	}
}

func TestCheckDuplicates(t *testing.T) {
	if err := CheckDuplicates([]int{3, 1, 2, -5}); err != nil {
		t.Errorf("aucun doublon attendu : %v", err)
	}
	if err := CheckDuplicates(nil); err != nil {
		t.Errorf("aucun doublon attendu : %v", err)
	}
	if err := CheckDuplicates([]int{3, 1, 3}); err == nil {
		t.Error("le doublon 3 aurait dû être refusé")
	}
}

func TestRanks(t *testing.T) {
	tests := []struct {
		nombres []int
		want    []int
	}{
		{[]int{40, 10, 50, 20, 30}, []int{3, 0, 4, 1, 2}},
		{[]int{-5, 1000, 3}, []int{0, 2, 1}},
		{[]int{5, 2, 9, 1, 7}, []int{2, 1, 4, 0, 3}},
		{[]int{7}, []int{0}},
		{nil, []int{}},
	}
	for _, tc := range tests {
		got := Ranks(tc.nombres)
		if !slices.Equal(got, tc.want) {
			t.Errorf("Ranks(%v) = %v, attendu %v", tc.nombres, got, tc.want)
		}
	}
}

func TestRanksDoNotModifyInput(t *testing.T) {
	nombres := []int{40, 10, 50, 20, 30}
	original := slices.Clone(nombres)
	_ = Ranks(nombres)
	if !slices.Equal(nombres, original) {
		t.Errorf("Ranks a modifié son entrée : %v", nombres)
	}
}
