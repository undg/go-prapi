package main

import (
	"encoding/json"
)

type Action string

const (
	// Message Actions
	GetCards   Action = "GetCards"
	GetOutputs Action = "GetOutputs"
	GetVolume  Action = "GetVolume"
	GetSchema  Action = "GetSchema"
	GetMute    Action = "GetMute"

	SetVolume  Action = "SetVolume"
	SetMute    Action = "SetMute"
	ToggleMute Action = "ToggleMute"
)

var availableCommands = []Action{
	GetCards,
	GetOutputs,
	GetVolume,
	GetSchema,
	GetMute,

	SetVolume,
	SetMute,
	ToggleMute,
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
	Action string `json:"actionIn" doc:"Action performed by API"`
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
