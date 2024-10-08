package ws

import (
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/undg/go-prapi/buildinfo"
	"github.com/undg/go-prapi/json"
	"github.com/undg/go-prapi/pactl"
	"github.com/undg/go-prapi/utils"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool)
var clientsMutex = &sync.Mutex{}

func upgraderCheckOrigin() {
	Upgrader.CheckOrigin = func(r *http.Request) bool {
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

		return utils.IsLocalIP(ip) || strings.HasPrefix(r.Host, "localhost")
	}
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("wsEndpoint visited by: %s %s", r.Host, r.RemoteAddr)

	upgraderCheckOrigin()

	conn, err := Upgrader.Upgrade(w, r, nil)
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

	initialResponse := json.Response{
		Action:  string(json.ActionGetSinks),
		Status:  json.StatusSuccess,
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
		var msg json.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading JSON: %v", err)
			}
			break
		}

		// Same Action and StatusSuccess if everyting is OK
		res := json.Response{
			Action: string(msg.Action),
			Status: json.StatusSuccess,
		}

		switch msg.Action {

		case json.ActionBroadcastStatus:
			status, _ := pactl.GetStatus()
			res.Payload = status

		case json.ActionGetSinks:
			status, _ := pactl.GetOutputs()
			res.Payload = status

		case json.ActionGetBuildInfo:
			b := buildinfo.Get()
			res.Payload = b

		case json.ActionSetSink:
			handleSetSink(msg, res)
		case json.ActionSetMute:
			handleSetMuted(msg, res)
		case json.ActionGetVolume:
			handleGetVolume(&res)
		case json.ActionGetMute:
			handleGetMute(&res)
		case json.ActionGetCards:
			handleGetCards(&res)
			handleGetOutputs(&res)
		case json.ActionGetSchema:
			handleGetSchema(&res)
		case json.ActionSetVolume:
			handleSetVolume(&res, msg.Payload.(float64))
		default:
			res.Error = "Command not found. Available actions: " + strings.Join(utils.ActionsToStrings(json.AvailableCommands), " ")
			res.Status = json.StatusActionError
		}

		handleServerLog(&msg, &res)

		if err := conn.WriteJSON(res); err != nil {
			log.Printf("Error writing JSON: %v", err)
			break
		}
	}
}
