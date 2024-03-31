package main

import (
	"fmt"
	"log"

	"mrogalski.eu/go/pulseaudio"
)

func getVol() string {
	client, err := pulseaudio.NewClient()
	if err != nil {
		panic(err)
	}
	vol, err := client.Volume()
	if err != nil {
		log.Println(err)
	}

	defer client.Close()
	// Use `client` to interact with PulseAudio

	return fmt.Sprint(vol)
}
