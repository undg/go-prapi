package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	once             sync.Once
	globalStopTicker = make(chan struct{})
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("wsEndpoint visited by:", r.Host, r.RemoteAddr)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		switch {
		case r.Host == "localhost"+PORT:
			return true
		default:
			return false
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	stopTicker := make(chan struct{})

	defer func() {
		conn.Close()
		close(stopTicker)
	}()

	go readerJSON(conn, stopTicker)

	// Wait for the connection to close
	<-stopTicker
}

func readerJSON(conn *websocket.Conn, stopTicker chan struct{}) {

	for {
		msg := Message{}
		res := Response{}

		res.Status = StatusSuccess

		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("ERROR conn.ReadJSON", err)
			res.Status = StatusErrorInvalidJSON
			res.Error = "Invalid JSON"
		}

		switch msg.Action {
		case GetCards:
			handleGetCards(&res)

		case GetOutputs:
			handleGetOutputs(&res)

		case GetVolume:
			handleGetVolume(&res)

		case GetMute:
			handleGetMute(&res)

		case SetVolume:
			handleSetVolume(&res, msg.Value.(float64))

		default:
			res.Error = "Command not found. Available actions: " + strings.Join(actionsToStrings(availableCommands), " ")
			res.Status = StatusActionError
		}

		res.Action = string(msg.Action)

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
