package main

import "testing"

func TestSelectIconPng(t *testing.T) {
	var tests = []struct {
		batteryStatus int
		expected      string
	}{
		{80, batteryFullPNG},
		{75, batteryFullPNG},
		{60, batteryHalfPNG},
		{50, batteryHalfPNG},
		{30, batteryHalfPNG},
		{25, batteryHalfPNG},
		{20, batteryEmptyPNG},
		{10, batteryEmptyPNG},
		{00, batteryEmptyPNG},
	}
	for _, tc := range tests {
		status := status{
			batteryStatus: tc.batteryStatus,
		}
		tc.expected = path + "/" + tc.expected
		if result := selectIconPng(status); result != tc.expected {
			t.Errorf("Expected %s, but got %s", tc.expected, result)
		}
	}
}
