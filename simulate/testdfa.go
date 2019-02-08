package main

import (
	"github.com/amandeepm003/scooterdfa"
	"fmt"
	"encoding/json"
)

func main() {



	dfa := voidfa.BuildDFA(voidfa.AvailableTransitions)

	fmt.Println("Automata %+v", dfa)


	fmt.Println("Before_Transition_Trigger_AutomataState %+v", voidfa.ToStateString(dfa.State()))
	err := dfa.Trigger(voidfa.StateTerminated,voidfa.RoleUser)
	fmt.Println("After_Transition_Trigger_AutomataState %+v", voidfa.ToStateString(dfa.State()))

	if err!=nil {
		fmt.Println(err.Error())
		var dfaError voidfa.DFAError
		unmerr := json.Unmarshal([]byte(err.Error()), &dfaError)
		if unmerr!=nil {
			fmt.Println("Err_Json_unmarshalling %+v", unmerr)
		} else {
			fmt.Println("Error JSON Marshalled: %+v",dfaError)
		}
		fmt.Println(dfaError.Status, dfaError.Type)
	}
	
}