package main

import (
	"errors"
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

func mute(isMuted bool) Audio {
	c := clientOpen()

	switch {
	case isMuted:
		clientMute(c)
	case !isMuted:
		clientUnMute(c)
	}

	volume := clientVolume(c)
	mute := clientMuteStatus(c)

	clientClose(c)
	return Audio{volume, mute}
}

type CardInfo struct {
	Name  string
	Index uint32
}

func getCards() ([]CardInfo, error) {
	c := clientOpen()

	cards, err := c.Cards()
	if err != nil {
		log.Println("ERROR clientVolume c.Volume", err)
		return nil, errors.New("ERROR clientVolume c.Volume")
	}

	clientClose(c)

	cardsInfo := []CardInfo{}

	for _, card := range cards {
		cardInfo := CardInfo{
			Name:  card.Name,
			Index: card.Index,
		}
		cardsInfo = append(cardsInfo, cardInfo)
	}

	return cardsInfo, nil
}

type OutputsInfo struct {
	ActiveIndex int
	CardID      string
	CardName    string
	PortName    string
	Available   bool
	PortID      string
}

func getOutputs() ([]OutputsInfo, error) {
	c := clientOpen()

	output, activeIndex, err := c.Outputs()
	if err != nil {
		log.Println("ERROR clientVolume c.Volume", err)
		return nil, errors.New("ERROR clientVolume c.Volume")
	}

	clientClose(c)

	outputsInfo := []OutputsInfo{}
	for _, output := range output {
		cardInfo := OutputsInfo{
			ActiveIndex: activeIndex,
			CardID:      output.CardID,
			CardName:    output.CardName,
			PortName:    output.PortName,
			Available:   output.Available,
			PortID:      output.PortID,
		}
		outputsInfo = append(outputsInfo, cardInfo)
	}

	return outputsInfo, nil
}
