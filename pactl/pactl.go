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

func SetSinkVolume(sinkName string, volume string) {
	errPrefix := "ERROR [SetSinkVolume()]"
	volumeInPercent := fmt.Sprint(volume) + "%"

	cmd := exec.Command("pactl", "set-sink-volume", sinkName, volumeInPercent)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s pactl set-sink-volume: %s\n", errPrefix, err)
		log.Printf("%s pactl set-sink-volume: {SINK_NAME: %s ; VOLUME: %s}\n", errPrefix, sinkName, volumeInPercent)
	}
}

func SetSinkMuted(sinkName string, mutedCmd string) {
	errPrefix := "ERROR [SetSinkMuted()]"

	cmd := exec.Command("pactl", "set-sink-mute", sinkName, mutedCmd)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s pactl set-sink-mute: %s\n", errPrefix, err)
		log.Printf("%s pactl set-sink-mute: {SINK_NAME: %s ; MUTED: %s}\n", errPrefix, sinkName, mutedCmd)
	}
}

func SetSinkInputVolume(sinkInputID string, volume string) {
	errPrefix := "ERROR [SetSinkInputVolume()]"
	volumeInPercent := volume + "%"

	cmd := exec.Command("pactl", "set-sink-input-volume", sinkInputID, volumeInPercent)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s pactl set-sink-input-volume: %s\n", errPrefix, err)
		log.Printf("%s pactl set-sink-input-volume: {SINK_INPUT_ID: %s ; VOLUME: %s}\n", errPrefix, sinkInputID, volumeInPercent)
	}
}

func SetSinkInputMuted(sinkInputID string, mutedCmd string) {
	errPrefix := "ERROR [SetSinkInputMuted()]"

	cmd := exec.Command("pactl", "set-sink-mute", sinkInputID, mutedCmd)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s pactl set-sink-mute: %s\n", errPrefix, err)
		log.Printf("%s pactl set-sink-mute: {SINK_INPUT_ID: %s ; MUTED: %s}\n", errPrefix, sinkInputID, mutedCmd)
	}
}

func adaptOutputs(p gen.PactlSinkJSON) Output {
	errPrefix := "ERROR [adaptOutputs()]"

	frontLeft, err := strconv.Atoi(strings.Trim(p.Volume.FrontLeft.ValuePercent, "%"))
	if err != nil {
		log.Printf("%s parse FRONT_LEFT to int: %s\n", errPrefix, err)
	}

	frontRight, err := strconv.Atoi(strings.Trim(p.Volume.FrontLeft.ValuePercent, "%"))
	if err != nil {
		log.Printf("%s parse FRONT_RIGHT to int: %s\n", errPrefix, err)
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
	errPrefix := "ERROR [GetOutputs()]"

	cmd := exec.Command("pactl", "--format=json", "list", "sinks")
	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var pactlSinks []gen.PactlSinkJSON
	err = json.Unmarshal(cmdOutput, &pactlSinks)
	if err != nil {
		log.Printf("%s json.Unmarshal: %s\n", errPrefix, err)
	}

	sinks := make([]Output, len(pactlSinks))
	for i, ps := range pactlSinks {
		sinks[i] = adaptOutputs(ps)
	}

	return sinks, nil
}

func adaptApps(p gen.PactlAppsJSON) App {
	errPrefix := "ERROR [adaptApps()]"

	frontLeft, err := strconv.Atoi(strings.Trim(p.Volume.FrontLeft.ValuePercent, "%"))
	if err != nil {
		log.Printf("%s parse FRONT_LEFT to INT: %s\n", errPrefix, err)
	}

	frontRight, err := strconv.Atoi(strings.Trim(p.Volume.FrontLeft.ValuePercent, "%"))
	if err != nil {
		log.Printf("%s parse FRONT_RIGHT to INT: %s\n", errPrefix, err)
	}

	return App{
		ID:       int(p.Index),
		OutputID: int(p.Sink),
		Label:    p.Properties.Application_Name,
		Volume:   (frontLeft + frontRight) / 2,
		Muted:    p.Mute,
	}
}

func GetApps() []App {
	errPrefix := "ERROR [pactl.GetApps()]"

	cmd := exec.Command("pactl", "--format=json", "list", "sink-inputs")
	cmdOutput, err := cmd.Output()
	if err != nil {
		log.Printf("%s cmd.Output(): %s", errPrefix, err)
	}

	var pactlApps []gen.PactlAppsJSON
	err = json.Unmarshal(cmdOutput, &pactlApps)
	if err != nil {
		log.Printf("%s json.Unmarshal(): %s", errPrefix, err)
	}

	apps := make([]App, len(pactlApps))
	for i, ps := range pactlApps {
		apps[i] = adaptApps(ps)
	}

	return apps
}

func ListenForChanges(callback func()) {
	cmd := exec.Command("pactl", "subscribe")
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "on sink-input") || strings.Contains(line, "on sink") {
		log.Println(line)
			callback()
		}
	}
}

func GetStatus() Status {
	errPrefix := "ERROR [GetStatus()]"

	outputs, err := GetOutputs()
	if err != nil {
		log.Printf("%s GetOutputs(): %s", errPrefix, err)
	}

	apps := GetApps()

	bi := buildinfo.Get()

	return Status{
		Outputs:   outputs,
		Apps:      apps,
		BuildInfo: *bi,
	}
}
