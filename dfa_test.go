package voidfa

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestBuildDFA(t *testing.T) {
	dfa := BuildDFA(AvailableTransitions)
	assert.Equal(t, dfa.state, StateReady)
}

func TestUserCanStartRidingFromReady(t *testing.T) {
	dfa := BuildDFA(AvailableTransitions)
	err := dfa.Trigger(StateRiding,RoleUser)
	assert.Nil(t, err)
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
