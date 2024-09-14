package main

import (
	"log"
	"time"
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

		res := Response{
			Action: string(ActionGetVolume),
			Status: StatusSuccess,
		}

		// @TODO (undg) 2024-09-13: Replace with rich data about all sinks and cards.
		handleGetVolume(&res)

		if res == prevRes {
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

		log.Printf("Volume broadcast sent to %d/%d clients. Value: %v", updatedClients, clientsCount, res.Value)
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
			Value: int(timeoutDuration.Seconds()),
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

			log.Printf("Alive status sent to %d/%d clients. Next ping in %ds", updatedClients, clientsCount, res.Value)
		}
}
