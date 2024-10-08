package ws

import (
	"log"
	"reflect"
	"time"

	"github.com/undg/go-prapi/json"
	"github.com/undg/go-prapi/pactl"
	"github.com/undg/go-prapi/utils"
)

var prevRes json.Response

const writeWait = 10 * time.Second

func BroadcastUpdates() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		clientsMutex.Lock()
		clientsCount := len(clients)
		clientsMutex.Unlock()

		if clientsCount == 0 {
			if utils.DEBUG {
				log.Println("No clients connected. Skipping VOLUME update.")
			}

			continue
		}

		// Same Action and StatusSuccess if everyting is OK
		res := json.Response{
			Action: string(json.ActionBroadcastStatus),
			Status: json.StatusSuccess,
		}

		status, err := pactl.GetStatus()
		if err != nil {
			log.Println("ERROR pactl.GetStatus() in broadcastUpdates()", err)
		}
		res.Payload = status

		equal := reflect.DeepEqual(res, prevRes)
		if equal {
			continue
		}

		prevRes = res

		clientsMutex.Lock()
		updatedClients := 0
		for conn := range clients {
			conn.SetWriteDeadline(time.Now().Add(writeWait))
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
