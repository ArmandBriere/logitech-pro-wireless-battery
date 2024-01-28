package main

import (
	"fmt"
	"testing"
)

type shellCommandMock struct {
	RunFunc func() error
}

// Output implements IShellCommand.
func (shellCommandMock) Output() ([]byte, error) {
	return []byte(""), nil
}

func (sc shellCommandMock) Run() error {
	return sc.RunFunc()
}

func TestSendNotification(t *testing.T) {
	var args []string
	var actualProgramName string

	// Mock shellCommand and extract programName and args
	shellCommand = func(programName string, arg ...string) IShellCommand {
		args = arg
		actualProgramName = programName
		return shellCommandMock{
			RunFunc: func() error {
				return nil
			},
		}
	}

	// Prepare test data
	defaultArgs := []string{
		"-t",
		"5000",
		"-i",
	}
	var tests = []struct {
		status      status
		batteryIcon string
	}{
		{
			status: status{
				batteryStatus: 80,
				headsetName:   "Logitech Pro Wireless",
			},
			batteryIcon: batteryFullPNG,
		},
		{
			status: status{
				batteryStatus: 50,
				headsetName:   "SteelSeries Arctis 7",
			},
			batteryIcon: batteryHalfPNG,
		},
		{
			status: status{
				batteryStatus: 20,
				headsetName:   "Sony WH-1000XM3",
			},
			batteryIcon: batteryEmptyPNG,
		},
	}

	// Run tests
	for _, tc := range tests {
		expectedArgs := append(defaultArgs, []string{
			fmt.Sprintf("/%s", tc.batteryIcon),
			fmt.Sprintf("%s - Battery at %d%%", tc.status.headsetName, tc.status.batteryStatus),
		}...)

		sendNotification(tc.status)

		if actualProgramName != "notify-send" {
			t.Errorf("Bad program name - Expected %s, but got %s", "notify-send", actualProgramName)
		}
		if len(args) != len(expectedArgs) {
			t.Errorf("args len and expectedArgs are not the same, %d != %d", len(expectedArgs), len(args))
		}
		for i, arg := range args {
			if arg != expectedArgs[i] {
				t.Errorf("args[%d] and expectedArgs[%d] are not the same, %s != %s", i, i, expectedArgs[i], arg)
			}
		}
	}
}
