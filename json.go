package main

import (
	"encoding/json"
)

// @TODO (undg) 2024-09-11: fiddle with Request and Message and decide the wichone is better.
const (
	// Request Actions and Types
	ActionGet = "get"
	ActionSet = "set"

	TypeCards   = "cards"
	TypeOutputs = "outputs"
	TypeVol     = "vol"
	TypeSchema  = "schema"
	TypeMute    = "mute"
	TypeToggle  = "toggle"

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


type Request struct {
	Action string      `json:"action" doc:"Action to perform: get or set"`
	Type   string      `json:"type" doc:"Type of the action: cards, outputs, vol, schema, mute, toggle"`
	Value  interface{} `json:"value,omitempty" doc:"Optional value for set actions"`
}

type Action string

// Message is an request from the client
type Message struct {
	Action Action      `json:"action" doc:"Action to perform fe. GetVolume, SetVolume, SetMute..."`
	Value  interface{} `json:"value,omitempty" doc:"Optional value for set actions"`
}

type Response struct {
	Action string      `json:"action" doc:"Action performed"`
	Type   string      `json:"type" doc:"Type of the action"`
	Status int16       `json:"status" doc:"Status code"`
	Value  interface{} `json:"value" doc:"Resulting value"`
	Error  string      `json:"error,omitempty" doc:"Error message"`
}

const (
	StatusSuccess int16 = 1000
	StatusError   int16 = 1001
)

func (r Response) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"action": r.Action,
		"status": r.Status,
		"type":   r.Type,
	}

	if r.Value != nil {
		data["value"] = r.Value
	}

	if r.Error != "" {
		data["error"] = r.Error
	}

	return json.Marshal(data)
}
