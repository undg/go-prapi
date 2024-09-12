package main

import (
	"encoding/json"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	response := Response{
		Action: string(GetCards),
		Status: StatusSuccess,
		Value:  "test value",
	}

	expected := `{"action":"GetCards","status":4000,"value":"test value"}`

	assertJSON(t, response, expected)
}

func TestMarshalJSONWithError(t *testing.T) {
	response := Response{
		Action: string(GetCards),
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
