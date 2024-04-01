package main

import (
	"encoding/json"
	"log"
)

func marshalAudio(a Audio) []byte {
	jsonData, err := json.Marshal(map[string]interface{}{
		"volume": a.volume,
		"mute":   a.mute,
	})
	if err != nil {
		log.Println(err)
	}

	return jsonData
}
