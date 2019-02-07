package voidfa

import (
	"sync"
	"time"
)

/*
To Model Mealy's DFA here

Rationale:
In the Mealy model, the output is a function of both the present state and the input (event + role input).
In the Moore model, the output is a function of only the present state
 */



var transitions = availableTransitions //load transitions for this DFA from available data dictionary


//Mealy machine by its nature would like to fuse some of incoming states, but for now lets keep that for later :)

type VoiDFA struct {
	triggers map[State]map[State][]Role
	mutex sync.Mutex //semaphore protecting dfa to act on trigger
	state   State
}

func (dfa *VoiDFA) State() State {
	return dfa.state
}


func buildDFA(inputTransitions []DFATransition) *VoiDFA {

	if len(inputTransitions) < 1 {
		//Not a DFA: Reason -> https://stackoverflow.com/questions/13791205/can-a-dfa-have-epsilon-lambda-transitions
		return nil
	}

	availableTriggers := map[State]map[State][]Role{}

	for _, t := range inputTransitions {
		if _, found := availableTriggers[t.PrevState]; !found {
			availableTriggers[t.PrevState] = map[State][]Role{}
		}
		availableTriggers[t.PrevState][t.NewState] = t.Roles
	}

	return &VoiDFA{state: inputTransitions[0].PrevState, triggers: availableTriggers}
}



func (dfa *VoiDFA) Trigger(destState State, role Role) *DFAError {
	dfa.mutex.Lock()
	defer dfa.mutex.Unlock()

	if role == RoleAdmin { //Super users who can do everything.
		dfa.state = destState
		return nil
	}

	roles, valid := dfa.triggers[dfa.state][destState]
	if !valid {
		return &DFAError{Type: "Invalid Transition", Detail: "Role: " + toRoleString(role)+ " CurrState: "+ toStateString(dfa.state) + " DestState: "+ toStateString(destState), Status: 400, TimeStamp: time.Now()}
	}

	// Check if permissions are valid
	if !rolePermitted(roles,role) {
		return &DFAError{Type: "Access Denied", Detail: "Role:  " + toRoleString(role) + " CurrState: "+ toStateString(dfa.state) + " DestState: "+ toStateString(destState), Status: 403, TimeStamp:time.Now()}
	}

	dfa.state = destState //Reached here after validation, so set DFA to this new state
	return nil
}






