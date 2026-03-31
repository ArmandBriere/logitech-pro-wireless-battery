package main

import "path/filepath"

func selectIcon(workingDir string, batteryStatus int) string {
	selectedIcon := batteryEmptyPNG
	if batteryStatus >= 75 {
		selectedIcon = batteryFullPNG
	} else if batteryStatus >= 25 {
		selectedIcon = batteryHalfPNG
	}
	return filepath.Join(workingDir, selectedIcon)
}
