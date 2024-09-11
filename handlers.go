package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

func handleServerLog(msg *Message, res *Response) {
	if msg != nil {
		msgBytes, err := json.MarshalIndent(msg, "", "	")
		if err != nil {
			log.Println("ERROR serverLog json.MarshalIndent", err)
		}
		fmt.Println("request:", string(msgBytes))
	}

	resBytes, err := res.MarshalJSON()
	if err != nil {
		log.Println("ERROR serverLog res.MarshalJson", err)
	}
	fmt.Println("response:", string(resBytes))
}

func handleSetVolume(res *Response, vol float32) {
	audio := setVol(vol)
	res.Value = strconv.FormatFloat(float64(audio.volume), 'f', -1, 32)
	res.Status = StatusSuccess
}

func handleGetVolume(res *Response) {
	audio := getVol()
	res.Value = strconv.FormatFloat(float64(audio.volume), 'f', -1, 32)
	res.Status = StatusSuccess

	handleServerLog(nil, res)
}

func handleGetMute(res *Response) {
	audio := getVol()
	res.Value = strconv.FormatBool(audio.mute)
	res.Status = StatusSuccess
}

func handleGetCards(res *Response) {
	cards, err := getCards()
	if err != nil {
		log.Println("ERROR readerJson GetCards", err)
		res.Error = "ERROR can't get cards information from the system"
		res.Status = StatusError
	}
	b, err := json.Marshal(cards)
	if err != nil {
		log.Println("ERROR readerJson json.Marshal", err)
		res.Error = "ERROR can't pull cards information"
		res.Status = StatusError
	}
	res.Value = string(b)
	res.Status = StatusSuccess
}

func handleGetOutputs(res *Response) {
	outputs, err := getOutputs()
	if err != nil {
		log.Println("ERROR readerJson getOutputs", err)
		res.Error = "ERROR can't get outputs information from the system"
		res.Status = StatusError
	}
	b, err := json.Marshal(outputs)
	if err != nil {
		log.Println("ERROR readerJson json.Marshal", err)
		res.Error = "ERROR can't pull outputs information"
		res.Status = StatusError
	}
	res.Value = string(b)
	res.Status = StatusSuccess
}
