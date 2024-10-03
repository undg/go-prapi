package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/undg/go-prapi/pactl"
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

func handleSetMuted(msg Message, res Response) {
	if sinkInfo, ok := msg.Payload.(map[string]interface{}); ok {
		name, ok := sinkInfo["name"].(string)
		if !ok {
			log.Printf("Error in [handleSetMuted]: sinkInfo['name'].(string) not ok")
		}
		muted, ok := sinkInfo["muted"].(bool)
		if !ok {
			log.Printf("Error in [handleSetMuted]: sinkInfo['muted'].(bool) not ok")
		}
		pactl.SetSinkMuted(name, muted)
		volStatus, err := pactl.GetStatus()
		if err != nil {
			log.Printf("Error in [handleSetMuted]: pactl.GetStatus(), err: %v", err)
		}
		res.Payload = volStatus
	} else {
		res.Error = "Invalid sink information format"
		res.Status = StatusActionError
	}
}

func handleSetMutedToggle(msg Message, res Response) {
	if sinkInfo, ok := msg.Payload.(map[string]interface{}); ok {
		name, ok := sinkInfo["name"].(string)
		if !ok {
			log.Printf("Error in HandleWebSocket [ActionSetMute]: sinkInfo['name'].(string) not ok")
		}
		muted, ok := sinkInfo["muted"].(bool)
		if !ok {
			log.Printf("Error in HandleWebSocket [ActionSetMute]: sinkInfo['muted'].(bool) not ok")
		}
		pactl.SetSinkMuted(name, muted)
		volStatus, err := pactl.GetStatus()
		if err != nil {
			log.Printf("Error in HandleWebSocket [ActionSetMute]: pactl.GetStatus(), err: %v", err)
		}
		res.Payload = volStatus
	} else {
		res.Error = "Invalid sink information format"
		res.Status = StatusActionError
	}
}

func handleSetSink(msg Message, res Response) {
	if sinkInfo, ok := msg.Payload.(map[string]interface{}); ok {
		name, ok := sinkInfo["name"].(string)
		if !ok {
			log.Printf("Error in HandleWebSocket [ActionSetSink]: sinkInfo['name'].(string) not ok")
		}
		volume, ok := sinkInfo["volume"].(string)
		if !ok {
			log.Printf("Error in HandleWebSocket [ActionSetSink]: sinkInfo['volume'].(string) not ok")
		}
		pactl.SetSink(name, volume)
		status, err := pactl.GetStatus()
		if err != nil {
			log.Printf("Error in HandleWebSocket [ActionSetSink]: pactl.GetStatus(), err: %v", err)
		}
		res.Payload = status
	} else {
		res.Error = "Invalid sink information format"
		res.Status = StatusActionError
	}
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
