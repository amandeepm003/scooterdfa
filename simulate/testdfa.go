package main

import (
	"github.com/amandeepm003/scooterdfa"
	"fmt"
)

func main() {



	dfa := voidfa.BuildDFA(voidfa.AvailableTransitions)

	fmt.Println("Automata %+v", dfa)


	fmt.Println("AutomataState %+v", voidfa.ToStateString(dfa.State()))
	err := dfa.Trigger(voidfa.StateReady,voidfa.RoleUser)
	fmt.Println("AutomataState %+v", voidfa.ToStateString(dfa.State()))

	fmt.Println("Error %+v", err)


}