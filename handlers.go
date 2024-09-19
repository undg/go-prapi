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
	res.Payload = audioValue.volume
	if DEBUG {
		log.Printf("handleSetVolume %s", res.Payload)
	}
}

func handleGetVolume(res *Response) {
	audio := getVol()
	res.Payload = strconv.FormatFloat(float64(audio.volume), 'f', -1, 32)
	if DEBUG {
		log.Printf("handleGetVolume %s", res.Payload)
	}
}

func handleGetMute(res *Response) {
	audio := getVol()
	res.Payload = strconv.FormatBool(audio.mute)
	if DEBUG {
		log.Printf("handleGetMute %s", res.Payload)
	}
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
	res.Payload = string(b)
	if DEBUG {
		log.Printf("handleGetCards %s", res.Payload)
	}
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
	res.Payload = string(b)
	if DEBUG {
		log.Printf("handleGetOutputs %s", res.Payload)
	}
}

func handleGetSchema(res *Response) {
	schema := GetSchemaJSON()
	res.Payload = schema
	if DEBUG {
		log.Printf("handleGetSchema %s", res.Payload)
	}
}
