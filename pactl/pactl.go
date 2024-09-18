package pactl

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type Sink struct {
	Name   string
	Label  string
	Volume int
	Muted  bool
}

func SetSink(sinkName string, volume string) {
	volumeInPercent := fmt.Sprint(volume) + "%"
	cmd := exec.Command("pactl", "set-sink-volume", sinkName, volumeInPercent)
	_, err := cmd.Output()
	if err != nil {
		log.Println("ERROR [SetSink]", err)
		log.Printf("ERROR [SetSink] SINK_NAME: %s ; VOLUME: %s", sinkName, volumeInPercent)
	}

}

func GetSinks() ([]Sink, error) {
	cmd := exec.Command("pactl", "list", "sinks")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	sinks := []Sink{}
	sinkBlocks := strings.Split(string(output), "Sink #")

	for _, block := range sinkBlocks[1:] {
		sink := Sink{}

		nameRe := regexp.MustCompile(`Name: (.+)`)
		if match := nameRe.FindStringSubmatch(block); len(match) > 1 {
			sink.Name = strings.TrimSpace(match[1])
		}

		volumeRe := regexp.MustCompile(`Volume:.*?(\d+)%`)
		if match := volumeRe.FindStringSubmatch(block); len(match) > 1 {
			fmt.Sscanf(match[1], "%d", &sink.Volume)
		}

		humanNameRe := regexp.MustCompile(`Description: (.+)`)
		if match := humanNameRe.FindStringSubmatch(block); len(match) > 1 {
			sink.Label = strings.TrimSpace(match[1])
		}

		sink.Muted = strings.Contains(block, "Mute: yes")

		sinks = append(sinks, sink)
	}

	return sinks, nil
}

func ListenForChanges(callback func()) {
	cmd := exec.Command("pactl", "subscribe")
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "sink") || strings.Contains(line, "server") {
			callback()
		}
	}
}
