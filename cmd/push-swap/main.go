package main

import (
	"fmt"
	"os"

	"pushswap/internal/pushswap"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		return
	}

	nombres, err := pushswap.ParseArgs(args)
	if err == nil {
		err = pushswap.CheckDuplicates(nombres)
	}
	if err != nil {
		echec()
	}

	s := pushswap.New(nombres)
	s.Sort()

	if !s.Done() {
		echec()
	}
	if err := s.WriteOps(os.Stdout); err != nil {
		echec()
	}
}

// echec affiche le seul message d'erreur autorisé par le sujet et sort
// avec le code 1.
func echec() {
	fmt.Fprintln(os.Stderr, "Error")
	os.Exit(1)
}
