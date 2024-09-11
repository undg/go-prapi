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
	Action Action      `json:"action" doc:"Action to perform fe. GetVolume, SetVolume, SetMute..."`
	Value  interface{} `json:"value,omitempty" doc:"Optional value for set actions"`
}

type Response struct {
	Action string      `json:"action" doc:"Action performed"`
	Status int16       `json:"status" doc:"Status code"`
	Value  interface{} `json:"value" doc:"Resulting value"`
	Error  string      `json:"error,omitempty" doc:"Error message"`
}

const (
	StatusSuccess int16 = 1000
	StatusError   int16 = 1001
	StatusActionError   int16 = 1002
	StatusValueError   int16 = 1003
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
