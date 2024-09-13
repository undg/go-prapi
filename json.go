package main

import (
	"encoding/json"
)

type Action string

const (
	// Message Actions
	ActionGetCards   Action = "GetCards"
	ActionGetOutputs Action = "GetOutputs"
	ActionGetVolume  Action = "GetVolume"
	ActionGetSchema  Action = "GetSchema"
	ActionGetMute    Action = "GetMute"

	ActionSetVolume  Action = "SetVolume"
	ActionSetMute    Action = "SetMute"
	ActionToggleMute Action = "ToggleMute"

	ActionImAlive Action = "ImAlive"
)

var availableCommands = []Action{
	ActionGetCards,
	ActionGetOutputs,
	ActionGetVolume,
	ActionGetSchema,
	ActionGetMute,

	ActionSetVolume,
	ActionSetMute,
	ActionToggleMute,

	ActionImAlive,
}

// Message is an request from the client
type Message struct {
	// Actions listed in availableCommands slice
	Action Action `json:"action" doc:"Action to perform fe. GetVolume, SetVolume, SetMute..."`
	// Paylod send with Set* actions if necessary
	Value interface{} `json:"value,omitempty" doc:"Paylod send with Set* actions if necessary"`
}

type Response struct {
	// Action performed by API
	Action string `json:"action" doc:"Action performed by API"`
	// Status code
	Status int16 `json:"status" doc:"Status code"`
	// Response payload
	Value interface{} `json:"value" doc:"Response payload"`
	// Error description if any
	Error string `json:"error,omitempty" doc:"Error description if any"`
}

const (
	StatusSuccess          int16 = 4000
	StatusError            int16 = 4001
	StatusActionError      int16 = 4002
	StatusValueError       int16 = 4003
	StatusErrorInvalidJSON int16 = 4004
)

func (r Response) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"action": r.Action,
		"status": r.Status,
	}

	if r.Value != nil {
		data["value"] = r.Value
	}

	if r.Error != "" {
		data["error"] = r.Error
	}

	return json.Marshal(data)
}
