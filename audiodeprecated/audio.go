package audiodeprecated

import (
	"errors"
	"log"

	"mrogalski.eu/go/pulseaudio"
)

type Audio struct {
	Volume float32
	Mute   bool
}

type CardInfo struct {
	Name  string
	Index uint32
}

type OutputsInfo struct {
	ActiveIndex int
	CardID      string
	CardName    string
	PortName    string
	Available   bool
	PortID      string
}


// don't forget to closeClient()
func ClientOpen() pulseaudio.Client {
	client, err := pulseaudio.NewClient()
	if err != nil {
		panic(err)
	}

	return *client
}

func ClientClose(c pulseaudio.Client) {
	defer c.Close()
}

func ClientToggleMute(c pulseaudio.Client) bool {
	mute, err := c.ToggleMute()
	if err != nil {
		log.Println("ERROR clientToggleMute c.ToggleMute", err)
	}
	return mute
}

func ClientMute(c pulseaudio.Client) {
	err := c.SetMute(true)
	if err != nil {
		log.Println("ERROR clientMute c.SetMute", err)
	}
}

func ClientUnMute(c pulseaudio.Client) {
	err := c.SetMute(false)
	if err != nil {
		log.Println("ERROR clientUnMute c.SetMute", err)
	}
}

func ClientMuteStatus(c pulseaudio.Client) bool {
	mute, err := c.Mute()
	if err != nil {
		log.Println("ERROR clientMuteStatus c.Mute", err)
	}
	return mute
}

func ClientVolume(c pulseaudio.Client) float32 {
	volume, err := c.Volume()
	if err != nil {
		log.Println("ERROR clientVolume c.Volume", err)
	}
	return volume
}

func SetVol(vol float32) Audio {
	c := ClientOpen()

	err := c.SetVolume(vol)
	if err != nil {
		log.Println("ERROR setVol c.SetVolume", err)
	}

	ClientUnMute(c)
	Volume := ClientVolume(c)
	Mute := ClientMuteStatus(c)

	ClientClose(c)
	return Audio{Volume, Mute}
}

func ToggleMute() Audio {
	c := ClientOpen()

	Volume := ClientVolume(c)
	Mute := ClientToggleMute(c)

	ClientClose(c)
	return Audio{Volume, Mute}
}

func Mute(isMuted bool) Audio {
	c := ClientOpen()

	switch {
	case isMuted:
		ClientMute(c)
	case !isMuted:
		ClientUnMute(c)
	}

	Volume := ClientVolume(c)
	Mute := ClientMuteStatus(c)

	ClientClose(c)
	return Audio{Volume, Mute}
}

func GetVol() Audio {
	c := ClientOpen()

	Volume := ClientVolume(c)
	Mute := ClientMuteStatus(c)

	ClientClose(c)
	return Audio{Volume, Mute}
}

func GetCards() ([]CardInfo, error) {
	c := ClientOpen()

	cards, err := c.Cards()
	if err != nil {
		log.Println("ERROR clientVolume c.Volume", err)
		return nil, errors.New("ERROR clientVolume c.Volume")
	}

	ClientClose(c)

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

func GetOutputs() ([]OutputsInfo, error) {
	c := ClientOpen()

	output, activeIndex, err := c.Outputs()
	if err != nil {
		log.Println("ERROR clientVolume c.Volume", err)
		return nil, errors.New("ERROR clientVolume c.Volume")
	}

	ClientClose(c)

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
