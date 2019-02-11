package voidfa

import (
	"sync"
	"time"
	"errors"
	"encoding/json"
	"fmt"
)

/*
To Model Mealy's DFA here

Rationale:
In the Mealy model, the output is a function of both the present state and the input (event + role input).
In the Moore model, the output is a function of only the present state
 */


//Mealy machine by its nature would like to fuse some of incoming states, but for now lets keep that for later :)

type VoiDFA struct {
	triggers map[State]map[State][]Role
	mutex sync.Mutex //semaphore protecting dfa to act on trigger
	state   State
	lastStateChange time.Time //To measure "inactive since 48 hours", will be internally updated on every state change success
	batteryStatus float64 //Unused now, but actual battery status ( e.g. can be used for validations when Collected->Dropped)
}

func (dfa *VoiDFA) State() State {
	return dfa.state
}


func BuildDFA(inputTransitions []DFATransition) *VoiDFA {

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



func (dfa *VoiDFA) Trigger(destState State, role Role) error {

	if destState==dfa.state { //State is always reachable by itself, irrespective of role, without counting transition or obtaining lock
		return nil
	}

	dfa.mutex.Lock()
	defer dfa.mutex.Unlock()


	if role == RoleAdmin { //Super users who can do everything.
		dfa.state = destState
		dfa.lastStateChange = time.Now()
		return nil
	}

	//Admin override is kept before "lastStateChange based check & automatic transition"(below).
	//This is done to ensure that admin action is achieved within one clock cycle, and no extra state transition flow

	sinceLastStateChange := time.Since(dfa.lastStateChange)
	if (sinceLastStateChange.Hours() >= 48 ) {
		dfa.state = StateUnknown
		dfaError := DFAError{Type: "Invalid Transition", Detail: "Role: " + ToRoleString(role)+ " CurrState: "+ ToStateString(dfa.state) + " DestState: "+ ToStateString(destState),
			Status: 400, TimeStamp: time.Now().Format("2006-01-02T15:04:05Z")}
		dfErrByes,_ :=json.Marshal(dfaError)
		return errors.New(string(dfErrByes))
	}

	loc, _ := time.LoadLocation("Europe/Copenhagen")
	//NOTE: LIMITATION here for time being, this central time is base for all VOIs across time zones!
	// Rather trigger request API caller's  IP address location should be loaded here for such validation
	timeNow := time.Now().In(loc)

	//Voi goes to sleep at 9.30pm
	startSleepTime := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 21, 30, 0, 0, loc)

	fmt.Printf("DEBUG ** timeNow",timeNow)
	fmt.Printf("DEBUG ** startSleepTime",startSleepTime)
	//Will reach here only for non-admin roles, basically admin can even ride in midnight :)
	if (dfa.state == StateReady && timeNow.After(startSleepTime) ){
		dfa.state = StateBounty
		dfa.lastStateChange = time.Now()
		return nil
	}

	roles, valid := dfa.triggers[dfa.state][destState]
	if !valid {
		dfaError := DFAError{Type: "Invalid Transition", Detail: "Role: " + ToRoleString(role)+ " CurrState: "+ ToStateString(dfa.state) + " DestState: "+ ToStateString(destState),
			Status: 400, TimeStamp: time.Now().Format("2006-01-02T15:04:05Z")}
		dfErrByes,_ :=json.Marshal(dfaError)
		return errors.New(string(dfErrByes))
	}

	// Check if permissions are valid
	if !rolePermitted(roles,role) {
		dfaError := DFAError{Type: "Access Denied", Detail: "Role:  " + ToRoleString(role) + " CurrState: "+ ToStateString(dfa.state) + " DestState: "+ ToStateString(destState),
			Status: 403, TimeStamp: time.Now().Format("2006-01-02T15:04:05Z")}
		dfErrByes,_ :=json.Marshal(dfaError)
		return errors.New(string(dfErrByes))
	}

	dfa.state = destState //Reached here after validation, so set DFA to this new state
	dfa.lastStateChange = time.Now()
	return nil
}






