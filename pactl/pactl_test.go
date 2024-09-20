package pactl

import (
	"regexp"
	"testing"
)

func TestGetSinks(t *testing.T) {
	t.Run("Volume", func(t *testing.T) {
		sink, _ := GetSinks()
		if sink[0].Volume < 0 {
			t.Errorf("Expected volume more than 0, but got %d", sink[0].Volume)
		}
	})
	namePattern := `^(alsa_output|bluez_sink|bluez_output|combined)\..*`
	re := regexp.MustCompile(namePattern)

	tests := []struct {
		Name  string
		Input string
		Want      bool
	}{
		{"ValidSink", "alsa_output.pci-0000_0c_00.4.analog-stereo", true},
		{"InvalidSink", "dupa", false},
		{"InvalidEmptySink", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			sink, _ := GetSinks()
			if got := re.MatchString(tt.Input); got != tt.Want {
				t.Errorf("Expected name should starts with pattern, but got [%s]", sink[0].Name)
			}
		})
	}
}
