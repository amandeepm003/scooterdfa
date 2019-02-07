package main

import (
	"github.com/amandeepm003/scooterdfa"
	"fmt"
)

func main() {



	dfa := voidfa.BuildDFA(voidfa.AvailableTransitions)

	fmt.Println("Automata %+v", dfa)


	fmt.Println("AutomataState %+v", dfa.State())
	err := dfa.Trigger(voidfa.StateDropped,voidfa.RoleAdmin)
	fmt.Println("AutomataState %+v", voidfa.ToStateString(dfa.State()))

	fmt.Println("Error %+v", err)


}