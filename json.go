package main

import (
	"encoding/json"
)

const (
	ActionGet = "get"
	ActionSet = "set"

	TypeCards   = "cards"
	TypeOutputs = "outputs"
	TypeVol     = "vol"
	TypeSchema  = "schema"
	TypeMute    = "mute"
	TypeToggle  = "toggle"
)

type Request struct {
	Action string      `json:"action" doc:"Action to perform: get or set"`
	Type   string      `json:"type" doc:"Type of the action: cards, outputs, vol, schema, mute, toggle"`
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
