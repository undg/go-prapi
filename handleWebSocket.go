package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/undg/go-prapi/pactl"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("wsEndpoint visited by: %s %s", r.Host, r.RemoteAddr)

	upgraderCheckOrigin()

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

	// Execute ActionGetStatus when a new client connects
	status, _ := pactl.GetStatus()

	initialResponse := Response{
		Action:  string(ActionGetSinks),
		Status:  StatusSuccess,
		Payload: status,
	}
	if err := conn.WriteJSON(initialResponse); err != nil {
		log.Printf("Error sending initial sinks data: %v", err)
	}

	// Cleanup after client is disconnected
	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		clientCount := len(clients)
		clientsMutex.Unlock()
		conn.Close()
		log.Printf("Client disconnected. Total clients: %d", clientCount)
	}()

	// Messaging system with client
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading JSON: %v", err)
			}
			break
		}

		// Same Action and StatusSuccess if everyting is OK
		res := Response{
			Action: string(msg.Action),
			Status: StatusSuccess,
		}

		switch msg.Action {

		case ActionBroadcastStatus:
			status, _ := pactl.GetStatus()
			res.Payload = status

		case ActionGetSinks:
			status, _ := pactl.GetOutputs()
			res.Payload = status

		case ActionSetSink:
			handleSetSink(msg, res)
		case ActionSetMute:
			handleSetMuted(msg, res)
		case ActionGetVolume:
			handleGetVolume(&res)
		case ActionGetMute:
			handleGetMute(&res)
		case ActionGetCards:
			handleGetCards(&res)
			handleGetOutputs(&res)
		case ActionGetSchema:
			handleGetSchema(&res)
		case ActionSetVolume:
			handleSetVolume(&res, msg.Payload.(float64))
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
