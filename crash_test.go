package main

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
)

func TestMultipleWebSocketConnections(t *testing.T) {
    url := "ws://192.168.1.110:8448/api/v1/ws"

    clientCount := 10
    messageCount := 1000

    clients := make([]*websocket.Conn, clientCount)
    var wg sync.WaitGroup

    // Connect all clients first
    for i := 0; i < clientCount; i++ {
        ws, _, err := websocket.DefaultDialer.Dial(url, nil)
        if err != nil {
            t.Fatalf("WebSocket connection failed: %v", err)
        }
        clients[i] = ws
        defer ws.Close()
    }

    // Send messages concurrently
    for j := 0; j < messageCount; j++ {
        wg.Add(clientCount)
        for i := 0; i < clientCount; i++ {
            go func(ws *websocket.Conn) {
                defer wg.Done()
                volume := fmt.Sprintf("%d", rand.Intn(101))
                err := ws.WriteMessage(websocket.TextMessage, []byte(`{"action": "SetSinkVolume", "payload": {"name": "alsa_output.platform-snd_aloop.0.analog-stereo", "volume": "`+volume+`"}}`))
                if err != nil {
                    t.Errorf("Failed to write message: %v", err)
                    return
                }

                _, _, err = ws.ReadMessage()
                if err != nil {
                    t.Errorf("Failed to read message: %v", err)
                    return
                }
            }(clients[i])
        }
        wg.Wait()
    }
}
