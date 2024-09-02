package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var stopTicker = make(chan struct{})

func readerJSON(conn *websocket.Conn) {

	defer close(stopTicker)

	for {
		msg := Request{}
		res := Response{}

		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("ERROR conn.ReadJSON", err)
			return
		}

		if msg.Action == ActionGet {
			switch msg.Type {
			case TypeCards:
				handleGetCards(&res)

			case TypeOutputs:
				handleGetOutputs(&res)

			case TypeVol:
				handleGetVolume(&res)

			case TypeMute:
				handleGetMute(&res)
			}
		}

		if msg.Action == ActionSet {
			switch msg.Type {
			case TypeVol:
				handleSetVolume(&res, msg.Value.(float32))
			}
		}

		serverLog(&msg, &res)

		if err := conn.WriteJSON(res); err != nil {
			log.Println(err)
		}

		select {
		case <-stopTicker:
			return
		default:
		}

	}
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

	serverLog(nil, res)
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

func serverLog(msg *Request, res *Response) {
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("wsEndpoint visited by:", r.Host, r.RemoteAddr)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		switch {
		case r.Host == "localhost"+PORT:
			return true
		default:
			return false
		}
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	go readerJSON(ws)

	go tickerVolume(stopTicker)
}

func tickerVolume(stop <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			handleGetVolume(&Response{})
		case <-stop:
			return
		}
	}
}
