package main

import (
	"encoding/json"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	response := Response{
		Action: string(ActionGetCards),
		Status: StatusSuccess,
		Payload:  "test payload",
	}

	expected := `{"action":"GetCards","payload":"test payload","status":4000}`

	assertJSON(t, response, expected)
}

func TestMarshalJSONWithError(t *testing.T) {
	response := Response{
		Action: string(ActionGetCards),
		Status: StatusError,
		Error:  "test error",
	}

	expected := `{"action":"GetCards","error":"test error","status":4001}`

	assertJSON(t, response, expected)
}

func assertJSON(t *testing.T, response Response, expected string) {
	result, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Error marshaling JSON: %v", err)
	}

	if string(result) != expected {
		t.Errorf("\nExpected %s\nGot      %s", expected, result)
	}
}
