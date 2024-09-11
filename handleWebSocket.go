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

var (
	once       sync.Once
	stopTicker = make(chan struct{})
)

func readerJSON(conn *websocket.Conn) {

	// defer func() {
	// 	// conn.Close()
	// 	once.Do(func() {
	// 		close(stopTicker)
	// 	})
	// }()

	for {
		msg := Message{}
		res := Response{}

		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("ERROR conn.ReadJSON", err)
			return
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
			value, ok := msg.Value.(float64)
			if !ok {
				res.Error = "Value is not a float"
				res.Status = StatusValueError
				break
			}
			res.Status = StatusSuccess
			res.Value = setVol(float32(value))
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

		// select {
		// case <-stopTicker:
		// 	return
		// default:
		// }
	}
}

func actionsToStrings(actions []Action) []string {
	strs := make([]string, len(actions))
	for i, action := range actions {
		strs[i] = string(action)
	}
	return strs
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("wsEndpoint visited by:", r.Host, r.RemoteAddr)

	// upgrader.CheckOrigin = func(r *http.Request) bool {
	// 	switch {
	// 	case r.Host == "localhost"+PORT:
	// 		return true
	// 	default:
	// 		return false
	// 	}
	// }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	defer conn.Close()

	readerJSON(conn)

	// tickerVolume(stopTicker)
}

