package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/danielgtaylor/huma/schema"
	"github.com/gorilla/websocket"
)

type Request struct {
	GetCards      bool     `json:"get_cards" doc:"NOT_IMPLEMENTED"`
	UseCard       string   `json:"use_card" doc:"NOT_IMPLEMENTED"`
	GetSinks      bool     `json:"get_sinks" doc:"NOT_IMPLEMENTED"`
	UseSink       string   `json:"use_sink" doc:"NOT_IMPLEMENTED"`
	GetVol        bool     `json:"get_vol" doc:"Request volume level (true)"`
	SetVol        *float32 `json:"set_vol" doc:"Set volume value as an float between 0.0 and 2.0"`
	Mute          *bool    `json:"mute" doc:"Mute or unMute (true/false)"`
	ToggleMute    bool     `json:"toggle_mute (true)"`
	GetJsonSchema bool     `json:"get_json_schema" doc:"NOT_IMPLEMENTED get this JSON schema (true)"`
}

func readerJson(conn *websocket.Conn) {
	for {
		msg := Request{}
		res := Result{}

		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("ERROR conn.ReadJSON", err)
			return
		}

		switch {
		case msg.GetCards:
			b, err := json.Marshal(getCards())
			if err != nil {
				log.Println("ERROR readerJson json.Marshal", err)
				res.error = "ERROR can't pull cards informations"
				break
			}
			res.response = string(b)
		case msg.UseCard != "":
			res.response = "NOT_IMPLEMENTED"
		case msg.GetSinks:
			res.response = "NOT_IMPLEMENTED"
		case msg.UseSink != "":
			res.response = "NOT_IMPLEMENTED"
		case msg.GetVol:
			res.Audio = getVol()
		case msg.SetVol != nil && *msg.SetVol >= 0 && *msg.SetVol < 2.0:
			res.Audio = setVol(*msg.SetVol)
		case msg.Mute != nil:
			res.Audio = mute(*msg.Mute)
		case msg.ToggleMute:
			res.Audio = toggleMute()
		case msg.GetJsonSchema:
			s, err := schema.Generate(reflect.TypeOf(Request{}))
			if err != nil {
				log.Println("ERROR readerJson schema.Generate", err)
				res.error = "ERROR can't generate schema"
				break
			}
			bytes, err := json.Marshal(s)
			if err != nil {
				log.Println("ERROR readerJson json.Marshal", err)
				res.error = "ERROR can't marshal JSON"
				break
			}
			res.schema = string(bytes)
		}

		serverLog(msg, res)

		if err := conn.WriteJSON(res); err != nil {
			log.Println(err)
		}
	}
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
