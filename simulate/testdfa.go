package main

import (
	"github.com/amandeepm003/scooterdfa"
	"fmt"
)

func main() {



	dfa := voidfa.BuildDFA(voidfa.AvailableTransitions)

	fmt.Printf("Automata %+v", dfa)

    



}