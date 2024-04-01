package main

import (
	"fmt"
	"log"

	"mrogalski.eu/go/pulseaudio"
)

// don't forget to closeClient()
func openClient() pulseaudio.Client {
	client, err := pulseaudio.NewClient()
	if err != nil {
		panic(err)
	}

	return *client
}

func closeClient(client pulseaudio.Client) {
	defer client.Close()
}

func getVol() float32 {
	client := openClient()

	volume, err := client.Volume()
	if err != nil {
		log.Println("ERROR getVol client.Volume", err)
	}

	closeClient(client)
	return volume
}

func setVol(vol float32) float32 {
	client := openClient()

	err := client.SetVolume(vol)
	if err != nil {
		log.Println("ERROR setVol client.SetVolume", err)
	}

	currVol, err := client.Volume()
	if err != nil {
		log.Println("ERROR setVol client.Volume", err)
	}

	closeClient(client)
	return currVol
}

func toggleMute() float32 {
	client := openClient()

	b, err := client.ToggleMute()
	if err != nil {
		log.Println("ERROR setVol client.SetVolume", err)
	}

	if b {
		fmt.Println("mute on")
	} else {
		fmt.Println("mute off")
	}

	currVol, err := client.Volume()
	if err != nil {
		log.Println("ERROR setVol client.Volume", err)
	}

	closeClient(client)
	return currVol
}
