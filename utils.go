package voidfa

import "time"

type DFAError struct { //RFC 7807 compliant error struct
	Type      string    `json:"type"`  //Invalid Transition, Access Denied
	Title     string    `json:"error"` //Transition Failed
	Status    int       `json:"status"` // n.a. here but e.g. 400 - invalid request, 403 - access denied
	Detail    string    `json:"detail"` // Role, PrevState, CurrState
	TimeStamp time.Time `json:"timestamp"` //Timestamp at attempt
}


func rolePermitted (roles []Role, matchRole Role) bool {
	for _, r := range roles {
		if matchRole == r {
			return true
		}
	}
	return false
}


func toStateString(state State) string {
	var stringval string
	switch state {
	case StateReady: stringval = "StateReady"
	case StateBatteryLow: stringval = "StateBatteryLow"
	case StateBounty: stringval = "StateBounty"
	case StateRiding: stringval = "StateRiding"
	case StateCollected: stringval = "StateCollected"
	case StateDropped: stringval = "StateDropped"
	case StateServiceMode: stringval = "StateServiceMode"
	case StateTerminated: stringval = "StateTerminated"
	case StateUnknown: stringval = "StateUnknown"
	default: stringval = ""
	}
	return stringval
}

func toRoleString(role Role) string {
	var roleval string
	switch role {
	case RoleUser: roleval = "RoleUser"
	case RoleHunter: roleval = "RoleHunter"
	case RoleAdmin: roleval = "RoleAdmin"
	case RoleSysInternal: roleval = "RoleSysInternal"
	default: roleval = ""
	}
	return roleval
}
