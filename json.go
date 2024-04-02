package main

import (
	"encoding/json"
	"log"
)

type Audio struct {
	volume float32
	mute   bool
}

type Result struct {
	Audio
	schema  string
	message string
	error   string
}

func marshalResult(a Result) []byte {
	jsonData, err := json.Marshal(map[string]interface{}{
		"volume": a.volume,
		"mute":   a.mute,
		"schema": a.schema,
	})

	if err != nil {
		log.Println(err)
	}

	return jsonData
}

func (a Result) MarshalJSON() ([]byte, error) {
	jsonData, err := json.Marshal(map[string]interface{}{
		"volume": a.volume,
		"mute":   a.mute,
		"schema": a.schema,
	})

	return jsonData, err
}
