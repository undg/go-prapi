package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("wsEndpoint visited by:", r.Host, r.RemoteAddr)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return false
		}

		ip := net.ParseIP(host)
		if ip == nil {
			return false
		}

		return isLocalIP(ip) || strings.HasPrefix(r.Host, "localhost")
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		conn.Close()
	}()

	go readerJSON(conn)
}

func readerJSON(conn *websocket.Conn) {
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

		case GetSchema:
			handleGetSchema(&res)

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

	}
}

func isLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	if ip4 := ip.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		}
	}

	return false
}
