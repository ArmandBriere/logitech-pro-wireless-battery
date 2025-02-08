package main

import (
	"testing"
)

const sonaarOutput = `
PRO X Wireless Gaming Headset
     Device path  : /dev/hidraw1
     USB id       : 046d:0ABA
     Codename     : PRO Headset
     Kind         : headset
     Battery: 60% 3897mV , 0.
Lightspeed Receiver
  Device path  : /dev/hidraw4
  USB id       : 046d:C547
  Serial       : F961A462
  Has 1 paired device(s) out of a maximum of 2.
  Notifications: wireless, software present (0x000900)
  Device activity counters: 1=97

  1: PRO X Wireless
     Device path  : None
     Battery: 86%, 0.
`
const mouseName = "PRO X Wireless"

func TestGetBatteryPercentage(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput int
		expectError    bool
	}{
		{
			input: `Lightspeed Receiver
1: PRO X Wireless
Battery: 85%`,
			expectedOutput: 85,
			expectError:    false,
		},
		{
			input: `Lightspeed Receiver
1: PRO X Wireless
Battery: 50%`,
			expectedOutput: 50,
			expectError:    false,
		},
		{
			input: `Lightspeed Receiver
1: PRO X Wireless
Battery: 10%`,
			expectedOutput: 10,
			expectError:    false,
		},
		{
			input: `Lightspeed Receiver
1: PRO X Wireless
Battery: not a number`,
			expectedOutput: 0,
			expectError:    true,
		},
		{
			input: `Lightspeed Receiver
1: Another Device
Battery: 75%`,
			expectedOutput: 0,
			expectError:    true,
		},
		{
			input: `PRO X Wireless Gaming Headset
Whatever
Battery: 100%`,
			expectedOutput: 0,
			expectError:    true,
		},
		{
			input:          sonaarOutput,
			expectedOutput: 86,
			expectError:    false,
		},
	}

	for _, test := range tests {
		output, err := getMouseBattery(test.input, mouseName)
		if (err != nil) != test.expectError {
			t.Errorf("getBatteryPercentage(%q) returned error %v, expected error: %v", test.input, err, test.expectError)
		}
		if output != test.expectedOutput {
			t.Errorf("getBatteryPercentage(%q) = %d, expected %d", test.input, output, test.expectedOutput)
		}
	}
}
