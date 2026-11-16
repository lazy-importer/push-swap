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

	ok, err := pushswap.Check(nombres, os.Stdin)
	if err != nil {
		echec()
	}

	if ok {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}

func echec() {
	fmt.Fprintln(os.Stderr, "Error")
	os.Exit(1)
}
