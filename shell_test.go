package main

import (
	"fmt"
	"testing"
)

type shellCommandMock struct {
	args []string
	name string
}

func (s *shellCommandMock) Output() ([]byte, error) {
	return []byte(""), nil
}

func (s *shellCommandMock) Run() error {
	return nil
}

func TestSendNotification(t *testing.T) {
	workingDir := "/home/user/project"

	var tests = []struct {
		name          string
		deviceName    string
		batteryStatus int
		batteryIcon   string
	}{
		{
			name:          "full battery",
			deviceName:    "Logitech Pro Wireless",
			batteryStatus: 80,
			batteryIcon:   workingDir + "/" + batteryFullPNG,
		},
		{
			name:          "half battery",
			deviceName:    "SteelSeries Arctis 7",
			batteryStatus: 50,
			batteryIcon:   workingDir + "/" + batteryHalfPNG,
		},
		{
			name:          "low battery",
			deviceName:    "Sony WH-1000XM3",
			batteryStatus: 20,
			batteryIcon:   workingDir + "/" + batteryEmptyPNG,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := &shellCommandMock{}
			m := &monitor{
				devices:    make(map[string]deviceStatus),
				execCmd:    func(name string, arg ...string) IShellCommand { mock.name = name; mock.args = arg; return mock },
				workingDir: workingDir,
			}

			m.sendNotification(tc.deviceName, tc.batteryStatus)

			if mock.name != "notify-send" {
				t.Errorf("Bad program name - Expected %s, but got %s", "notify-send", mock.name)
			}

			expectedArgs := []string{
				"-t",
				"5000",
				"-i",
				tc.batteryIcon,
				fmt.Sprintf("%s - Battery at %d%%", tc.deviceName, tc.batteryStatus),
			}

			if len(mock.args) != len(expectedArgs) {
				t.Errorf("args len and expectedArgs are not the same, %d != %d", len(expectedArgs), len(mock.args))
			}
			for i, arg := range mock.args {
				if arg != expectedArgs[i] {
					t.Errorf("args[%d] and expectedArgs[%d] are not the same, %s != %s", i, i, expectedArgs[i], arg)
				}
			}
		})
	}
}
