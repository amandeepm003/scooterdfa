package voidfa

//States of DFA involved

type State int

const (
	StateReady State = iota
	StateBatteryLow
	StateBounty
	StateRiding
	StateCollected
	StateDropped

	// Not in service, only admin can claim.
	StateServiceMode //Workshop
	StateTerminated  //Dismantled/Destroyed
	StateUnknown     //After 48 hours of staying ready, could be false status/stolen i.e. tagged Unknown in system
)

//Roles in scene affecting DFA

type Role int
const (
	//External User Roles
	RoleUser Role = iota
	RoleHunter
	RoleAdmin
	//Internal Role
	RoleSysInternal
)

// Enumerating all events that can happen in DFA, and by whom (and above)
type Event struct {
	OutputState State
	Role  Role
}

var (
	EventDeploy = Event {OutputState:StateReady, Role: RoleHunter} //voi is made available in street
	EventBatteryLowTrigger = Event {OutputState: StateBatteryLow, Role: RoleSysInternal} //not just riding, lying around idle also will decay battery
	EventHunterAlertedForBounty = Event {OutputState: StateBounty} //zero state transition
	EventHunterCollected = Event {OutputState:StateCollected, Role:RoleHunter}
	EventHunterDropped = Event {OutputState:StateDropped, Role:RoleHunter}
	EventRideStart = Event {OutputState:StateRiding}
	EventRideComplete = Event {OutputState: StateReady}
	EventSleep = Event{OutputState: StateBounty, Role: RoleSysInternal}  //At 9.30pm it will go to sleep/bounty mode
	EventUnknown = Event {OutputState: StateUnknown, Role: RoleSysInternal} //Inactive 48hrs, Tampered, Stolen, Drowned, Anything here
	EventStandby = Event {OutputState: StateServiceMode, Role: RoleAdmin}
	EventTerminate = Event {OutputState:StateTerminated, Role:RoleAdmin}
)

//Mealy Machine where every state transition can be justified back on input event
type DFATransition struct {
	PrevState State
	NewState State
	Event Event //Place to document/log a transition, furthermore (for non-admin) DFA transition must match to specific event, unlike NFA
	Roles []Role
}

var AvailableTransitions = []DFATransition{
	//NewState is resultant state, event "could be" most probable cause
	{PrevState: StateReady, NewState: StateRiding, Event: EventRideStart, Roles: []Role{RoleUser, RoleHunter, RoleAdmin}},
	{PrevState: StateRiding, NewState: StateReady, Event: EventRideStart, Roles: []Role{RoleUser, RoleHunter, RoleAdmin}},
	{PrevState: StateRiding, NewState: StateBatteryLow, Event: EventBatteryLowTrigger, Roles: []Role{RoleSysInternal}}, //
	{PrevState: StateBatteryLow, NewState: StateBounty, Event:EventHunterAlertedForBounty, Roles: []Role{RoleSysInternal}},
	{PrevState: StateBounty, NewState: StateCollected, Event: EventHunterCollected, Roles: []Role{RoleHunter, RoleAdmin}},
	{PrevState: StateCollected, NewState: StateDropped, Event: EventHunterDropped, Roles: []Role{RoleHunter, RoleAdmin}},

	{PrevState: StateReady, NewState: StateUnknown, Event: EventUnknown, Roles: []Role{RoleSysInternal}},
	{PrevState: StateDropped, NewState: StateReady, Event: EventDeploy, Roles: []Role{RoleHunter, RoleAdmin}},
}



/*
//Questions/Assumptions

1. There is one impossible state to go from Ready to Unknown state after 48 hours.
As from Ready, every day at 9.30PM, vehicle goes to Bounty state.
So there will always be a state change within 1 day and hence the 48 hours will never be achieved ?

 */
