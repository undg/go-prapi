package pactl

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/undg/go-prapi/buildinfo"
	gen "github.com/undg/go-prapi/pactl/generated"
)

type Status = struct {
	Outputs   []Output            `json:"outputs" doc:"List of output devices"`
	Apps      []App               `json:"apps" doc:"List of applications"`
	BuildInfo buildinfo.BuildInfo `json:"buildInfo" doc:"Build information"`
}

type Output struct {
	ID     int    `json:"id" doc:"The id of the sink. Same  as name"`
	Name   string `json:"name" doc:"The name of the sink. Same as id"`
	Label  string `json:"label" doc:"Human-readable label for the sink"`
	Volume int    `json:"volume" doc:"Current volume level of the sink"`
	Muted  bool   `json:"muted" doc:"Whether the sink is muted"`
}

type App struct {
	ID       int    `json:"id" doc:"The id of the sink. Same  as name"`
	OutputID int    `json:"outputId" doc:"Id of parrent device, same as output.id"`
	Label    string `json:"label" doc:"Human-readable label for the sink"`
	Volume   int    `json:"volume" doc:"Current volume level of the sink"`
	Muted    bool   `json:"muted" doc:"Whether the sink is muted"`
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

func SetSinkMuted(sinkName string, muted bool) {
	mutedCmd := "0"
	if muted {
		mutedCmd = "1"
	}

	cmd := exec.Command("pactl", "set-sink-mute", sinkName, mutedCmd)
	_, err := cmd.Output()
	if err != nil {
		log.Println("ERROR [SetSinkMuted]", err)
		log.Printf("ERROR [SetSinkMuted] SINK_NAME: %s ; muted: %s", sinkName, mutedCmd)
	}
}

func adaptOutputs(p gen.PactlSinkJSON) Output {
	frontLeft, err := strconv.Atoi(strings.Trim(p.Volume.FrontLeft.ValuePercent, "%"))
	if err != nil {
		log.Println("ERROR adaptSink, parse front_left to int", err)
	}

	frontRight, err := strconv.Atoi(strings.Trim(p.Volume.FrontLeft.ValuePercent, "%"))
	if err != nil {
		log.Println("ERROR adaptSink, parse front_right to int", err)
	}

	return Output{
		ID:     int(p.Index),
		Name:   p.Name,
		Label:  p.Description,
		Volume: (frontLeft + frontRight) / 2,
		Muted:  p.Mute,
	}
}

func GetOutputs() ([]Output, error) {
	cmd := exec.Command("pactl", "--format=json", "list", "sinks")
	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var pactlSinks []gen.PactlSinkJSON
	err = json.Unmarshal(cmdOutput, &pactlSinks)
	if err != nil {
		log.Println("ERROR Unmarshal pactlSinks in GetSinks.", err)
	}

	sinks := make([]Output, len(pactlSinks))
	for i, ps := range pactlSinks {
		sinks[i] = adaptOutputs(ps)
	}

	return sinks, nil
}

func adaptApps(p gen.PactlAppsJSON) App {
	frontLeft, err := strconv.Atoi(strings.Trim(p.Volume.FrontLeft.ValuePercent, "%"))
	if err != nil {
		log.Println("ERROR adaptSink, parse front_left to int", err)
	}

	frontRight, err := strconv.Atoi(strings.Trim(p.Volume.FrontLeft.ValuePercent, "%"))
	if err != nil {
		log.Println("ERROR adaptSink, parse front_right to int", err)
	}

	return App{
		ID:       int(p.Index),
		OutputID: int(p.Sink),
		Label:    p.Properties.Application_Name,
		Volume:   (frontLeft + frontRight) / 2,
		Muted:    p.Mute,
	}
}

func GetApps() ([]App, error) {
	cmd := exec.Command("pactl", "--format=json", "list", "sink-inputs")
	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var pactlApps []gen.PactlAppsJSON
	err = json.Unmarshal(cmdOutput, &pactlApps)
	if err != nil {
		log.Println("ERROR Unmarshal pactlApps in GetApps.", err)
	}

	apps := make([]App, len(pactlApps))
	for i, ps := range pactlApps {
		apps[i] = adaptApps(ps)
	}

	return apps, nil
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

func GetStatus() (Status, error) {
	outputs, err := GetOutputs()
	if err != nil {
		log.Println("ERROR GetOutputs() in GetStatus", err)
	}

	apps, err := GetApps()
	if err != nil {
		log.Println("ERROR GetApps() in GetStatus()", err)
	}

	bi := buildinfo.Get()

	return Status{
		Outputs:   outputs,
		Apps:      apps,
		BuildInfo: *bi,
	}, nil
}
