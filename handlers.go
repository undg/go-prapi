package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

func handleServerLog(msg *Message, res *Response) {
	log.Println("")
	if msg != nil {
		msgBytes, err := json.MarshalIndent(msg, "", "	")
		if err != nil {
			log.Printf("ERROR serverLog json.MarshalIndent %s", err)
		}
		log.Printf("Client message: %s", string(msgBytes))
	}

	resBytes, err := res.MarshalJSON()
	if err != nil {
		log.Printf("ERROR serverLog res.MarshalJson %s", err)
	}
	fmt.Printf("LOG response: %s", string(resBytes))
}

func handleSetVolume(res *Response, vol float64) {
	audioValue := setVol(float32(vol))
	res.Value = audioValue.volume
	log.Printf("handleSetVolume %s", res.Value)
}

func handleGetVolume(res *Response) {
	audio := getVol()
	res.Value = strconv.FormatFloat(float64(audio.volume), 'f', -1, 32)
	log.Printf("handleGetVolume %s", res.Value)
}

func handleGetMute(res *Response) {
	audio := getVol()
	res.Value = strconv.FormatBool(audio.mute)
	log.Printf("handleGetMute %s", res.Value)
}

func handleGetCards(res *Response) {
	cards, err := getCards()
	if err != nil {
		log.Printf("ERROR readerJson GetCards %s", err)
		res.Error = "ERROR can't get cards information from the system"
		res.Status = StatusError
	}
	b, err := json.Marshal(cards)
	if err != nil {
		log.Printf("ERROR readerJson json.Marshal %s", err)
		res.Error = "ERROR can't pull cards information"
		res.Status = StatusError
	}
	res.Value = string(b)
	log.Printf("handleGetCards %s", res.Value)
}

func handleGetOutputs(res *Response) {
	outputs, err := getOutputs()
	if err != nil {
		log.Printf("ERROR readerJson getOutputs %s", err)
		res.Error = "ERROR can't get outputs information from the system"
		res.Status = StatusError
	}
	b, err := json.Marshal(outputs)
	if err != nil {
		log.Printf("ERROR readerJson json.Marshal %s", err)
		res.Error = "ERROR can't pull outputs information"
		res.Status = StatusError
	}
	res.Value = string(b)
	log.Printf("handleGetOutputs %s", res.Value)
}

func handleGetSchema(res *Response) {
	schema := GetSchemaJSON()
	res.Value = schema
	log.Printf("handleGetSchema %s", res.Value)
}
