package main

import (
	"errors"
	"os/exec"
	"strings"
)

type sink struct {
	name         string
	friendlyName string
	isDefault    bool
}

func (sink *sink) setDefault() error {
	cmd := exec.Command("pactl", "set-default-sink", sink.name)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func ensurePactlOrExit() error {
	cmd := exec.Command("pactl", "info")
	if err := cmd.Run(); err != nil {
		return errors.New("pactl was not found or no PulseAudio-compatible server is running")
	}
	return nil
}

func sinks() ([]sink, error) {
	cmd := exec.Command("bash", "-c", "pactl list sinks | grep -E 'Name|Description' | cut -d ':' -f2-")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.New("unable to retrieve sinks")
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines)%2 != 0 {
		return nil, errors.New("unexpected output format, uneven number of lines")
	}

	defaultSinkName, err := defaultSinkName()
	if (err) != nil {
		return nil, errors.New("unable to retrieve default sink")
	}

	sinks := make([]sink, 0, len(lines)/2)
	for i := 0; i < len(lines)-1; i += 2 {
		name := strings.TrimSpace(lines[i])
		sinks = append(sinks, sink{
			name:         name,
			friendlyName: strings.TrimSpace(lines[i+1]),
			isDefault:    name == defaultSinkName,
		})
	}
	return sinks, nil
}

func defaultSinkName() (string, error) {
	cmd := exec.Command("bash", "-c", "pactl info | grep 'Default Sink:' | cut -d ':' -f2-")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
