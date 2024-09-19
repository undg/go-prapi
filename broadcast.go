package main

import (
	"log"
	"reflect"
	"time"

	"github.com/undg/go-prapi/pactl"
)

var prevRes Response

func broadcastUpdates() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		clientsMutex.Lock()
		clientsCount := len(clients)
		clientsMutex.Unlock()

		if clientsCount == 0 {
			if DEBUG {
				log.Println("No clients connected. Skipping VOLUME update.")
			}

			continue
		}

		// Same Action and StatusSuccess if everyting is OK
		res := Response{
			Action: string(ActionGetSinks),
			Status: StatusSuccess,
		}

		sinks, _ := pactl.GetSinks()
		res.Payload = sinks

		equal := reflect.DeepEqual(res, prevRes)
		if equal {
			continue
		}

		prevRes = res

		clientsMutex.Lock()
		updatedClients := 0
		for conn := range clients {
			err := conn.WriteJSON(res)
			if err != nil {
				log.Printf("Error broadcast VOLUME update to client: %v", err)
				conn.Close()
				delete(clients, conn)
			} else {
				updatedClients++
			}
		}
		clientsMutex.Unlock()

		log.Printf("Volume broadcast sent to %d/%d clients. Value: %v", updatedClients, clientsCount, res.Payload)
	}
}

func broadcastImAlive() {
	timeoutDuration := 30 * time.Second
	ticker := time.NewTicker(timeoutDuration)
	defer ticker.Stop()
	for range ticker.C {
		clientsMutex.Lock()
		clientsCount := len(clients)
		clientsMutex.Unlock()

		if clientsCount == 0 {
			if DEBUG {
				log.Println("No clients connected. Skipping ALIVE ping.")
			}

			continue
		}

		res := Response{
			Action: string(ActionImAlive),
			Status: StatusSuccess,
			Payload: int(timeoutDuration.Seconds()),
		}

		clientsMutex.Lock()
		updatedClients := 0
		for conn := range clients {
			err := conn.WriteJSON(res)
			if err != nil {
				log.Printf("Error sending alive status to client: %v", err)
					conn.Close()
					delete(clients, conn)
				} else {
					updatedClients++
				}
			}
			clientsMutex.Unlock()

			log.Printf("Alive status sent to %d/%d clients. Next ping in %ds", updatedClients, clientsCount, res.Payload)
		}
}
