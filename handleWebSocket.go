package main

import (
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
	log.Printf("wsEndpoint visited by: %s %s", r.Host, r.RemoteAddr)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("Error splitting host and port: %v", err)
			return false
		}

		ip := net.ParseIP(host)
		if ip == nil {
			log.Printf("Invalid IP: %s", host)
			return false
		}

		return isLocalIP(ip) || strings.HasPrefix(r.Host, "localhost")
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	clientsMutex.Lock()
	clients[conn] = true
	clientCount := len(clients)
	clientsMutex.Unlock()

	log.Printf("New client connected. Total clients: %d", clientCount)

	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		clientCount := len(clients)
		clientsMutex.Unlock()
		conn.Close()
		log.Printf("Client disconnected. Total clients: %d", clientCount)
	}()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading JSON: %v", err)
			}
			break
		}

		res := Response{
			Action: string(msg.Action),
			Status: StatusSuccess,
		}

		switch msg.Action {
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

		handleServerLog(&msg, &res)

		if err := conn.WriteJSON(res); err != nil {
			log.Printf("Error writing JSON: %v", err)
			break
		}
	}
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
