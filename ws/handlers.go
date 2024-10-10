package ws

import (
	j "encoding/json"
	"fmt"
	"log"

	"github.com/undg/go-prapi/audiodeprecated"
	"github.com/undg/go-prapi/json"
	"github.com/undg/go-prapi/pactl"
	"github.com/undg/go-prapi/utils"
)

func handleSetSinkVolume(msg json.Message, res json.Response) {
	logPrefix := "Error [handleSetSinkVolume()]:"
	if sinkInfo, ok := msg.Payload.(map[string]interface{}); ok {
		name, ok := sinkInfo["name"].(string)
		if !ok {
			log.Printf("%s sinkInfo['name'].(string) not ok\n", logPrefix)
		}

		volume, ok := sinkInfo["volume"].(string)
		if !ok {
			log.Printf("%s sinkInfo['volume'].(string) not ok\n", logPrefix)
		}

		pactl.SetSinkVolume(name, volume)

		status, err := pactl.GetStatus()
		if err != nil {
			log.Printf("%s pactl.GetStatus(), err: %s\n", logPrefix, err)
		}

		res.Payload = status
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

func handleSetSinkMuted(msg json.Message, res json.Response) {
	logPrefix := "Error [handleSetSinkMuted()]:"
	if sinkInfo, ok := msg.Payload.(map[string]interface{}); ok {
		name, ok := sinkInfo["name"].(string)
		if !ok {
			log.Printf("%s sinkInfo['name'].(string) not ok\n", logPrefix)
		}

		muted, ok := sinkInfo["muted"].(string)
		if !ok {
			log.Printf("%s sinkInfo['muted'].(bool) not ok\n", logPrefix)
		}

		pactl.SetSinkMuted(name, muted)

		volStatus, err := pactl.GetStatus()
		if err != nil {
			log.Printf("%s pactl.GetStatus(), err: %s\n", logPrefix, err)
		}

		res.Payload = volStatus
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

func handleSetSinkInputVolume(msg json.Message, res json.Response) {
	logPrefix := "Error [handleSetSinkInputVolume()]:"
	if sinkInputInfo, ok := msg.Payload.(map[string]interface{}); ok {
		id, ok := sinkInputInfo["id"].(float64)
			log.Printf("\n\n\n%v\n\n\n\n", msg.Payload)
		if !ok {
			log.Printf("%s sinkInfo['id'].(int) not ok\n", logPrefix)
		}

		volume, ok := sinkInputInfo["volume"].(string)
		if !ok {
			log.Printf("%s sinkInfo['volume'].(string) not ok\n", logPrefix)
		}

		pactl.SetSinkInputVolume(fmt.Sprintf("%.0f",id), volume)

		status, err := pactl.GetStatus()

		if err != nil {
			log.Printf("%s pactl.GetStatus(), err: %s\n", logPrefix, err)
		}
		res.Payload = status
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

func handleSetSinkInputMuted(msg json.Message, res json.Response) {
	logPrefix := "Error [handleSetSinkInputMuted()]:"
	if sinkInputInfo, ok := msg.Payload.(map[string]interface{}); ok {
		name, ok := sinkInputInfo["id"].(string)
		if !ok {
			log.Printf("%s sinkInfo['id'].(string) not ok\n", logPrefix)
		}

		muted, ok := sinkInputInfo["muted"].(string)
		if !ok {
			log.Printf("%s sinkInfo['muted'].(bool) not ok\n", logPrefix)
		}

		pactl.SetSinkMuted(name, muted)

		volStatus, err := pactl.GetStatus()
		if err != nil {
			log.Printf("%s pactl.GetStatus(), err: %s\n", logPrefix, err)
		}

		res.Payload = volStatus
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

func handleGetCards(res *json.Response) {
	cards, err := audiodeprecated.GetCards()
	if err != nil {
		log.Printf("ERROR readerJson GetCards %s\n", err)
		res.Error = "ERROR can't get cards information from the system"
		res.Status = json.StatusError
	}
	b, err := j.Marshal(cards)
	if err != nil {
		log.Printf("ERROR readerJson j.Marshal %s\n", err)
		res.Error = "ERROR can't pull cards information"
		res.Status = json.StatusError
	}
	res.Payload = string(b)
	if utils.DEBUG {
		log.Printf("handleGetCards %s\n", res.Payload)
	}
}

func handleGetOutputs(res *json.Response) {
	outputs, err := audiodeprecated.GetOutputs()
	if err != nil {
		log.Printf("ERROR readerJson getOutputs %s\n", err)
		res.Error = "ERROR can't get outputs information from the system"
		res.Status = json.StatusError
	}
	b, err := j.Marshal(outputs)
	if err != nil {
		log.Printf("ERROR readerJson j.Marshal %s\n", err)
		res.Error = "ERROR can't pull outputs information"
		res.Status = json.StatusError
	}
	res.Payload = string(b)
	if utils.DEBUG {
		log.Printf("handleGetOutputs %s\n", res.Payload)
	}
}

func handleGetSchema(res *json.Response) {
	schema := json.GetSchemaJSON()
	res.Payload = schema
	if utils.DEBUG {
		log.Printf("handleGetSchema %s\n", res.Payload)
	}
}

func handleServerLog(msg *json.Message, res *json.Response) {
	log.Println("")
	if msg != nil {
		msgBytes, err := j.MarshalIndent(msg, "", "	")
		if err != nil {
			log.Printf("ERROR serverLog j.MarshalIndent %s\n", err)
		}
		log.Printf("Client message: %s\n", string(msgBytes))
	}

	resBytes, err := res.MarshalJSON()
	if err != nil {
		log.Printf("ERROR serverLog res.MarshalJson %s\n", err)
	}
	fmt.Printf("LOG response: %s\n", string(resBytes))
}
