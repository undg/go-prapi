package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type Request struct {
	Action string      `json:"action" doc:"Action to perform: get or set"`
	Type   string      `json:"type" doc:"Type of the action: cards, outputs, vol, schema, mute, toggle"`
	Value  interface{} `json:"value,omitempty" doc:"Optional value for set actions"`
}

type Response struct {
	Action string      `json:"action" doc:"Action performed"`
	Type   string      `json:"type" doc:"Type of the action"`
	Value  interface{} `json:"value,omitempty" doc:"Resulting value"`
}

const (
	ActionGet = "get"
	ActionSet = "set"

	TypeCards   = "cards"
	TypeOutputs = "outputs"
	TypeVol     = "vol"
	TypeSchema  = "schema"
	TypeMute    = "mute"
	TypeToggle  = "toggle"
)

type GetRequest = string
type SetRequest = string

func readerJson(conn *websocket.Conn) {
	for {
		msg := Request{}
		res := Result{}

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

		}

		serverLog(msg, res)

		if err := conn.WriteJSON(res); err != nil {
			log.Println(err)
		}

	}
}

func handleGetVolume(res *Result) {
	audio := getVol()
	res.response = strconv.FormatFloat(float64(audio.volume), 'f', -1, 32)
}

func handleGetMute(res *Result) {
	audio := getVol()
	res.response = strconv.FormatBool(audio.mute)
}

func handleGetCards(res *Result) {
	cards, err := getCards()
	if err != nil {
		log.Println("ERROR readerJson GetCards", err)
		res.error = "ERROR can't get cards information from the system"
	}
	b, err := json.Marshal(cards)
	if err != nil {
		log.Println("ERROR readerJson json.Marshal", err)
		res.error = "ERROR can't pull cards information"
	}
	res.response = string(b)
}

func handleGetOutputs(res *Result) {
	outputs, err := getOutputs()
	if err != nil {
		log.Println("ERROR readerJson getOutputs", err)
		res.error = "ERROR can't get outputs information from the system"
	}
	b, err := json.Marshal(outputs)
	if err != nil {
		log.Println("ERROR readerJson json.Marshal", err)
		res.error = "ERROR can't pull outputs information"
	}
	res.response = string(b)
}

func serverLog(msg Request, res Result) {
	bytes, err := json.MarshalIndent(msg, "", "	")
	if err != nil {
		log.Println("ERROR readerJson json.MarshalIndent", err)
	}
	fmt.Println("request:", string(bytes))
	fmt.Println("response:", string(marshalResult(res)))
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
	go readerJson(ws)
}
