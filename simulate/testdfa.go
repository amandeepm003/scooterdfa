package main

import (
	"github.com/amandeepm003/scooterdfa"
	"bufio"
	"os"
	"fmt"
	"strings"
	"encoding/json"
)

func main() {

/*

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

*/

	/*
		Uncomment above block to explore core basic working inside a main function
		However, Below is actual rule interpreter for input commands and responses

	   Notes:
	  1. Concerns like (a)Post 9.30pm, and (b) No Status change past 48 hours

	 */

	newDFA := voidfa.BuildDFA(voidfa.AvailableTransitions)

	bufRead := bufio.NewReader(os.Stdin)


	for {
		// Print
		fmt.Printf("VOI DFA Current State = %s ", voidfa.ToStateString(newDFA.State()))
		fmt.Printf("\n Enter Trigger as '<State>:<Role>' (e.g. RIDING:HUNTER ).Else type 'help' ->")

		// Read

		trigger,err  := bufRead.ReadString('\n')
		if err!=nil {
			panic(err)
		}

		if trigger == "help" {
			fmt.Println("Available States: READY, RIDING, BATTLOW, BOUNTY, COLLECTED, DROPPED, UNKNOWN, TERMINATED")
			fmt.Println("Available Roles: USER, HUNTER, ADMIN, SYSINT")
			fmt.Println("Example: Use 'BOUNTY:SYSINT'(without quotes) as trigger to validate transition to Bounty State")
			continue
		}

		commandrules := strings.SplitN(trigger, ":", 2)
		if len(commandrules) != 2 {
			fmt.Printf("\n ERROR: Enter Trigger as '<State>:<Role>' (e.g. RIDING,HUNTER ).Else type 'help' ->")
			continue
		}

		commandrules[0] = strings.TrimSpace(commandrules[0])
		commandrules[1] = strings.TrimSpace(commandrules[1])

		state:= voidfa.StateUnknown
		role:= voidfa.RoleSysInternal

		fmt.Printf("Debug--> State :%s : Role: %s \n", commandrules[0],commandrules[1])

		switch (commandrules[0]) {
		case "READY": state = voidfa.StateReady
		case "RIDING": state =voidfa.StateRiding
		case "BATTLOW": state =voidfa.StateBatteryLow
		case "BOUNTY": state =voidfa.StateBounty
		case "COLLECTED": state =voidfa.StateCollected
		case "DROPPED": state =voidfa.StateDropped
		case "UNKNOWN": state =voidfa.StateUnknown
		case "TERMINATED": state =voidfa.StateTerminated
		default: {
			fmt.Println("\n Invalid State Input. Available States: READY, RIDING, BATTLOW, BOUNTY, COLLECTED, DROPPED, UNKNOWN, TERMINATED ")
			continue
		}

		}

		switch (commandrules[1]) {
		case "USER": role = voidfa.RoleUser
		case "HUNTER": role =voidfa.RoleHunter
		case "ADMIN": role =voidfa.RoleAdmin
		case "SYSINT": role =voidfa.RoleSysInternal
		default: {
			fmt.Println("\n Invalid Role Input. Available Roles: USER, HUNTER, ADMIN, SYSINT")
			continue
		}
		}

		err = newDFA.Trigger(state,role)

		if err!=nil {
			fmt.Println(err.Error())
			var dfaError voidfa.DFAError
			unmerr := json.Unmarshal([]byte(err.Error()), &dfaError)
			if unmerr!=nil {
				fmt.Println("Error Internal Error with DFA (Err_Json_unmarshalling) %+v", unmerr)
			}
			fmt.Printf("Trigger_Error %d %s \n", dfaError.Status, dfaError.Type)
		} else {
			fmt.Printf("Trigger_Success (New State = %s) \n", voidfa.ToStateString(newDFA.State()))
		}

	}

}