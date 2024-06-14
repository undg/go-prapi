package main

import (
	"encoding/json"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	response := Response{
		Action: "get",
		Type:   "cards",
		Status: StatusSuccess,
		Value:  "test value",
	}

	expected := `{"action":"get","status":1000,"type":"cards","value":"test value"}`

	assertJSON(t, response, expected)
}

func TestMarshalJSONWithError(t *testing.T) {
	response := Response{
		Action: "get",
		Type:   "cards",
		Status: StatusError,
		Error:  "test error",
	}

	expected := `{"action":"get","error":"test error","status":1001,"type":"cards"}`

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
