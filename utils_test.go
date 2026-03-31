package main

import "testing"

func TestSelectIcon(t *testing.T) {
	workingDir := "/home/user/project"
	var tests = []struct {
		batteryStatus int
		expected      string
	}{
		{80, workingDir + "/icons/battery-full.png"},
		{75, workingDir + "/icons/battery-full.png"},
		{60, workingDir + "/icons/battery-half.png"},
		{50, workingDir + "/icons/battery-half.png"},
		{30, workingDir + "/icons/battery-half.png"},
		{25, workingDir + "/icons/battery-half.png"},
		{20, workingDir + "/icons/battery-empty.png"},
		{10, workingDir + "/icons/battery-empty.png"},
		{00, workingDir + "/icons/battery-empty.png"},
	}
	for _, tc := range tests {
		if result := selectIcon(workingDir, tc.batteryStatus); result != tc.expected {
			t.Errorf("Expected %s, but got %s", tc.expected, result)
		}
	}
}
