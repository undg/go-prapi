package ws

import (
	j "encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/undg/go-prapi/audiodeprecated"
	"github.com/undg/go-prapi/json"
	"github.com/undg/go-prapi/pactl"
	"github.com/undg/go-prapi/utils"
)

func handleServerLog(msg *json.Message, res *json.Response) {
	log.Println("")
	if msg != nil {
		msgBytes, err := j.MarshalIndent(msg, "", "	")
		if err != nil {
			log.Printf("ERROR serverLog j.MarshalIndent %s", err)
		}
		log.Printf("Client message: %s", string(msgBytes))
	}

	resBytes, err := res.MarshalJSON()
	if err != nil {
		log.Printf("ERROR serverLog res.MarshalJson %s", err)
	}
	fmt.Printf("LOG response: %s", string(resBytes))
}

func handleSetMuted(msg json.Message, res json.Response) {
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
		res.Status = json.StatusActionError
	}
}

func handleSetMutedToggle(msg json.Message, res json.Response) {
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
		res.Status = json.StatusActionError
	}
}

func handleSetSink(msg json.Message, res json.Response) {
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
		res.Status = json.StatusActionError
	}
}

func handleSetVolume(res *json.Response, vol float64) {
	audioValue := audiodeprecated.SetVol(float32(vol))
	res.Payload = audioValue.Volume
	if utils.DEBUG {
		log.Printf("handleSetVolume %s", res.Payload)
	}
}

func handleGetVolume(res *json.Response) {
	audio := audiodeprecated.GetVol()
	res.Payload = strconv.FormatFloat(float64(audio.Volume), 'f', -1, 32)
	if utils.DEBUG {
		log.Printf("handleGetVolume %s", res.Payload)
	}
}

func handleGetMute(res *json.Response) {
	audio := audiodeprecated.GetVol()
	res.Payload = strconv.FormatBool(audio.Mute)
	if utils.DEBUG {
		log.Printf("handleGetMute %s", res.Payload)
	}
}

func handleGetCards(res *json.Response) {
	cards, err := audiodeprecated.GetCards()
	if err != nil {
		log.Printf("ERROR readerJson GetCards %s", err)
		res.Error = "ERROR can't get cards information from the system"
		res.Status = json.StatusError
	}
	b, err := j.Marshal(cards)
	if err != nil {
		log.Printf("ERROR readerJson j.Marshal %s", err)
		res.Error = "ERROR can't pull cards information"
		res.Status = json.StatusError
	}
	res.Payload = string(b)
	if utils.DEBUG {
		log.Printf("handleGetCards %s", res.Payload)
	}
}

func handleGetOutputs(res *json.Response) {
	outputs, err := audiodeprecated.GetOutputs()
	if err != nil {
		log.Printf("ERROR readerJson getOutputs %s", err)
		res.Error = "ERROR can't get outputs information from the system"
		res.Status = json.StatusError
	}
	b, err := j.Marshal(outputs)
	if err != nil {
		log.Printf("ERROR readerJson j.Marshal %s", err)
		res.Error = "ERROR can't pull outputs information"
		res.Status = json.StatusError
	}
	res.Payload = string(b)
	if utils.DEBUG {
		log.Printf("handleGetOutputs %s", res.Payload)
	}
}

func handleGetSchema(res *json.Response) {
	schema := json.GetSchemaJSON()
	res.Payload = schema
	if utils.DEBUG {
		log.Printf("handleGetSchema %s", res.Payload)
	}
}
