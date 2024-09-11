package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	once       sync.Once
	stopTicker = make(chan struct{})
)

func readerJSON(conn *websocket.Conn) {

	defer func() {
		conn.Close()
		once.Do(func() {
			close(stopTicker)
		})
	}()

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

		handleServerLog(&msg, &res)

		if err := conn.WriteJSON(res); err != nil {
			log.Println(err)
			break
		}

		select {
		case <-stopTicker:
			return
		default:
		}
	}
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

	go tickerVolume(stopTicker)

	go readerJSON(ws)

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
