package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type msg struct {
	Set        *float32
	Vol        *bool
	Mute       *bool
	UnMute     *bool
	ToggleMute *bool
	Msg        *string
}

func readerJson(conn *websocket.Conn) {
	for {
		m := msg{}

		err := conn.ReadJSON(&m)
		if err != nil {
			log.Println("ERROR conn.ReadJSON", err)
		}

		bytes, err := json.MarshalIndent(m, "", "	")
		if err != nil {
			log.Println("ERROR readerJson json.MarshalIndent", err)
		}

		fmt.Printf("Got JSON: \n%v\n", string(bytes))

		if err := conn.WriteJSON(m); err != nil {
			log.Println(err)
		}
	}
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}

		msg := string(p)
		msgOut := string(p)

		if msg == "get" {
			msgOut = string(marshalAudio(getVol()))
		}

		if msg == "set" {
			msgOut = string(marshalAudio(setVol(0.5)))
		}

		if msg == "muteToggle" {
			msgOut = string(marshalAudio(toggleMute()))
		}

		if msg == "mute" {
			msgOut = string(marshalAudio(mute()))
		}

		if msg == "unMute" {
			msgOut = string(marshalAudio(unMute()))
		}

		if err := conn.WriteMessage(messageType, []byte(msgOut)); err != nil {
			log.Println(err)
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("wsEndpoint visited by:", r.Host, r.RemoteAddr)

	// upgrader.CheckOrigin = func(r *http.Request) bool {
	// 	// @TODO (undg) 2024-04-01: r.Host bla bla bla
	// 	return true
	// }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	go readerJson(ws)
}
