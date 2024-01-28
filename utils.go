package main

// Select icon based on battery percentage
func selectIconPng(status status) string {
	var selectedIcon string = batteryEmptyPNG
	if status.batteryStatus >= 75 {
		selectedIcon = batteryFullPNG
	} else if status.batteryStatus >= 25 {
		selectedIcon = batteryHalfPNG
	}
	return path + "/" + selectedIcon
}
