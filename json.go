package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Audio struct {
	volume float32
	mute   bool
}

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

func (a Audio) MarshalJSON() ([]byte, error) {
	jsonData, err := json.Marshal(map[string]interface{}{
		"volume": a.volume,
		"mute":   a.mute,
	})

	return jsonData, err
}

func a() {
	j := string(marshalAudio(unMute()))
	fmt.Println(j)
}
func t() {
	audioInstance := unMute()
	m, err := json.Marshal(audioInstance)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(m))
}
