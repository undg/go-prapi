package main

import (
	"log"

	"mrogalski.eu/go/pulseaudio"
)

// don't forget to closeClient()
func clientOpen() pulseaudio.Client {
	client, err := pulseaudio.NewClient()
	if err != nil {
		panic(err)
	}

	return *client
}

func clientClose(c pulseaudio.Client) {
	defer c.Close()
}

func clientToggleMute(c pulseaudio.Client) bool {
	mute, err := c.ToggleMute()
	if err != nil {
		log.Println("ERROR clientToggleMute c.ToggleMute", err)
	}
	return mute
}

func clientMute(c pulseaudio.Client) {
	err := c.SetMute(true)
	if err != nil {
		log.Println("ERROR clientMute c.SetMute", err)
	}
}

func clientUnMute(c pulseaudio.Client) {
	err := c.SetMute(false)
	if err != nil {
		log.Println("ERROR clientUnMute c.SetMute", err)
	}
}

func clientMuteStatus(c pulseaudio.Client) bool {
	mute, err := c.Mute()
	if err != nil {
		log.Println("ERROR clientMuteStatus c.Mute", err)
	}
	return mute
}

func clientVolume(c pulseaudio.Client) float32 {
	volume, err := c.Volume()
	if err != nil {
		log.Println("ERROR clientVolume c.Volume", err)
	}
	return volume
}

func getVol() Audio {
	c := clientOpen()

	volume := clientVolume(c)
	mute := clientMuteStatus(c)

	clientClose(c)
	return Audio{volume, mute}
}

func setVol(vol float32) Audio {
	c := clientOpen()

	err := c.SetVolume(vol)
	if err != nil {
		log.Println("ERROR setVol c.SetVolume", err)
	}

	clientUnMute(c)
	volume := clientVolume(c)
	mute := clientMuteStatus(c)

	clientClose(c)
	return Audio{volume, mute}
}

func toggleMute() Audio {
	c := clientOpen()

	volume := clientVolume(c)
	mute := clientToggleMute(c)

	clientClose(c)
	return Audio{volume, mute}
}

func mute() Audio {
	c := clientOpen()

	clientMute(c)
	volume := clientVolume(c)
	mute := clientMuteStatus(c)

	clientClose(c)
	return Audio{volume, mute}
}
func unMute() Audio {
	c := clientOpen()

	clientUnMute(c)
	volume := clientVolume(c)
	mute := clientMuteStatus(c)

	clientClose(c)
	return Audio{volume, mute}
}
