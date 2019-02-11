package voidfa

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"time"
)


func TestBuildDFA(t *testing.T) {
	dfa := BuildDFA(AvailableTransitions)
	assert.Equal(t, dfa.state, StateReady)
}

func TestUserCanStartRidingFromReadyWithinAllowedTimes(t *testing.T) {
	dfa := BuildDFA(AvailableTransitions)
	err := dfa.Trigger(StateRiding,RoleUser)

	timeNow := time.Now()
	if timeNow.After(time.Date(timeNow.Year(),timeNow.Month(),timeNow.Day(),7,0,0,0, time.UTC)) &&
		timeNow.Before(time.Date(timeNow.Year(),timeNow.Month(),timeNow.Day(),21,30,0,0, time.UTC)) {
		assert.Nil(t, err)
	} else {
		assert.NotNil(t, err)
	}

}

func TestUserCannotTerminate(t *testing.T) {
	dfa := BuildDFA(AvailableTransitions)
	err := dfa.Trigger(StateTerminated,RoleUser)
	assert.NotNil(t, err)
	assert.NotEqual(t,StateTerminated,dfa.state )
}

func TestAdminCanTerminate(t *testing.T) {
	dfa := BuildDFA(AvailableTransitions)
	err := dfa.Trigger(StateTerminated,RoleAdmin)
	assert.Nil(t, err)
	assert.Equal(t, StateTerminated, dfa.state)
}

func TestUserCannotRideBountyVehicle(t *testing.T) {
	dfa := BuildDFA(AvailableTransitions)
	dfa.Trigger(StateBounty,RoleAdmin)
	assert.Equal(t, StateBounty, dfa.state)

	err:= dfa.Trigger(StateRiding,RoleUser)
	assert.NotNil(t, err)
	assert.NotEqual(t,StateRiding,dfa.state )
}

func TestHunterCannotRideBountyVehicle(t *testing.T) {
	dfa := BuildDFA(AvailableTransitions)
	dfa.Trigger(StateBounty,RoleAdmin)
	assert.Equal(t, StateBounty, dfa.state)

	err:= dfa.Trigger(StateRiding,RoleHunter)
	assert.NotNil(t, err)
	assert.NotEqual(t,StateRiding,dfa.state )
}

func TestHunterCanCollectBountyVehicle(t *testing.T) {
	dfa := BuildDFA(AvailableTransitions)
	dfa.Trigger(StateBounty,RoleAdmin)
	assert.Equal(t, StateBounty, dfa.state)

	err:= dfa.Trigger(StateCollected,RoleHunter)
	assert.Nil(t, err)
	assert.Equal(t,StateCollected,dfa.state )
}


func TestUserCannotClaimServiceModeVehicle(t *testing.T) {
	dfa := BuildDFA(AvailableTransitions)
	dfa.Trigger(StateServiceMode,RoleAdmin)
	assert.Equal(t, StateServiceMode, dfa.state)

	err:= dfa.Trigger(StateRiding,RoleUser)
	assert.NotNil(t, err)
	assert.Equal(t,StateServiceMode,dfa.state )
	//fmt.Printf("Err %v",err)

	var dfaError DFAError
	er := json.Unmarshal([]byte(err.Error()), &dfaError)
	assert.Nil(t,er,nil)
	assert.Equal(t,"Invalid Transition",dfaError.Type)
}


func TestHunterCannotClaimServiceModeVehicle(t *testing.T) {
	dfa := BuildDFA(AvailableTransitions)
	dfa.Trigger(StateServiceMode,RoleAdmin)
	assert.Equal(t, StateServiceMode, dfa.state)

	err:= dfa.Trigger(StateCollected,RoleHunter)
	assert.NotNil(t, err)
	assert.Equal(t,StateServiceMode,dfa.state )
	//fmt.Printf("Err %v",err)

	var dfaError DFAError
	er := json.Unmarshal([]byte(err.Error()), &dfaError)
	assert.Nil(t,er,nil)
	assert.Equal(t,"Invalid Transition",dfaError.Type)
	assert.Equal(t,400,dfaError.Status)
}


func TestUserCannotClaimBountyVehicle(t *testing.T) {
	dfa := BuildDFA(AvailableTransitions)
	dfa.Trigger(StateBounty,RoleAdmin)
	assert.Equal(t, StateBounty, dfa.state)

	err:= dfa.Trigger(StateCollected,RoleUser)
	assert.NotNil(t, err)
	assert.Equal(t,StateBounty,dfa.state )
	//fmt.Printf("Err %v",err)

	var dfaError DFAError
	er := json.Unmarshal([]byte(err.Error()), &dfaError)
	assert.Nil(t,er,nil)
	assert.Equal(t,"Access Denied",dfaError.Type)
	assert.Equal(t,403,dfaError.Status)
}
