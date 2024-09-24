package pactl

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	gen "github.com/undg/go-prapi/pactl/generated"
)

type Sink struct {
	ID     string `json:"id" doc:"The id of the sink. Same  as name"`
	Name   string `json:"name" doc:"The name of the sink. Same as id"`
	Label  string `json:"label" doc:"Human-readable label for the sink"`
	Volume int    `json:"volume" doc:"Current volume level of the sink"`
	Muted  bool   `json:"muted" doc:"Whether the sink is muted"`
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

func adaptSink(ps gen.PactlSinkJSON) Sink {
	frontLeft, err := strconv.Atoi(strings.Trim(ps.Volume.FrontLeft.ValuePercent, "%"))
	if err != nil {
		log.Println("ERROR adaptSink, parse front_left to int", err)
	}

	frontRight, err := strconv.Atoi(strings.Trim(ps.Volume.FrontLeft.ValuePercent, "%"))
	if err != nil {
		log.Println("ERROR adaptSink, parse front_right to int", err)
	}

	return Sink{
		ID:     ps.Name,
		Name:   ps.Name,
		Label:  ps.Description,
		Volume: (frontLeft + frontRight) / 2,
		Muted:  ps.Mute,
	}
}

func GetSinks() ([]Sink, error) {
	cmd := exec.Command("pactl", "--format=json", "list", "sinks")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var pactlSinks []gen.PactlSinkJSON
	err = json.Unmarshal(output, &pactlSinks)
	if err != nil {
		log.Println("ERROR Unmarshal pactlSinks in GetSinks.", err)
	}

	sinks := make([]Sink, len(pactlSinks))
	for i, ps := range pactlSinks {
		sinks[i] = adaptSink(ps)
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
