package main

import (
	"log"
	"time"
)

func broadcastUpdates() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		clientsMutex.Lock()
		clientsCount := len(clients)
		clientsMutex.Unlock()

		if clientsCount == 0 {
			if DEBUG {
				log.Println("No clients connected. Skipping volume update.")
			}

			continue
		}

		res := Response{
			Action: string(GetVolume),
		}
		handleGetVolume(&res)

		clientsMutex.Lock()
		updatedClients := 0
		for conn := range clients {
			err := conn.WriteJSON(res)
			if err != nil {
				log.Printf("Error sending volume update to client: %v", err)
				conn.Close()
				delete(clients, conn)
			} else {
				updatedClients++
			}
		}
		clientsMutex.Unlock()

		log.Printf("Volume update sent to %d/%d clients. Value: %v", updatedClients, clientsCount, res.Value)
	}
}
