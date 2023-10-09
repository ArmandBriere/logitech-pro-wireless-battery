package main

import (
	"flag"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	PercentageWarning = 30
	PercentageEmpty   = 10
	PollingInterval   = 5

	notificationTimeout = "5000"

	BATTERY_FULL_PNG  = "icons/battery-full.png"
	BATTERY_HALF_PNG  = "icons/battery-half.png"
	BATTERY_EMPTY_PNG = "icons/battery-empty.png"
)

var path string

type status struct {
	headsetName   string
	batteryStatus int
	step          int
}

var headsets map[string]status = make(map[string]status)

// Set logger and logger level
func setLogger() {
	isVerbose := flag.Bool("verbose", false, "display debug logs")
	flag.Parse()

	if *isVerbose {
		log.SetLevel(log.DebugLevel)
	}

	log.SetOutput(os.Stdout)
}

// Run headsetcontrol binary to get headset info
func execHeadsetcontrol() {
	output, err := exec.Command("/usr/bin/headsetcontrol", "-b").Output()
	if err != nil {
		log.Debug("Headset not found.")
		return
	}

	log.Debug(string(output))

	headsetFoundSplit := strings.Split(string(output), "Found ")

	if len(strings.Split(string(output), "Found ")) == 1 {
		log.Debug("Headset not found.")
		return
	}

	headsetName := strings.Split(headsetFoundSplit[1], "!")[0]
	headsetBatterySplit := strings.Split(string(output), "Battery: ")
	if len(headsetBatterySplit) == 1 {
		log.Debug("Headset usb is connected but the headset is in sleep mode.")
		return
	}

	batteryStatus, _ := strconv.Atoi(strings.Split(headsetBatterySplit[1], "%")[0])

	step := batteryStatus / 5

	log.Debug("Headset name: " + headsetName)
	log.Debug("Battery status: " + strconv.Itoa(batteryStatus))
	log.Debug("Step: " + strconv.Itoa(step))

	previousStatus, exists := headsets[headsetName]

	if !exists {
		headsets[headsetName] = status{
			headsetName:   headsetName,
			batteryStatus: batteryStatus,
			step:          step,
		}
	} else if batteryStatus/5 != previousStatus.step {
		headsets[headsetName] = status{
			headsetName:   headsetName,
			batteryStatus: batteryStatus,
			step:          step,
		}
		sendNotification(headsets[headsetName])
	}
}

// Select icon based on battery percentage
func selectIconPng(status status) string {
	var selectedIcon string
	if status.batteryStatus > 75 {
		selectedIcon = BATTERY_FULL_PNG
	} else if status.batteryStatus > 25 {
		selectedIcon = BATTERY_HALF_PNG
	} else {
		selectedIcon = BATTERY_EMPTY_PNG
	}
	return path + "/" + selectedIcon
}

// Send notification using notify-send
func sendNotification(status status) {
	message := status.headsetName + " - Battery at " + strconv.Itoa(status.batteryStatus) + "%"

	icon := selectIconPng(status)
	_, err := exec.Command(
		"/usr/bin/notify-send",
		"-t", notificationTimeout,
		"-i", icon,
		message).Output()
	if err != nil {
		log.Error(err.Error())
	}
}

func main() {
	setLogger()
	path, _ = os.Getwd()
	for {
		execHeadsetcontrol()
		time.Sleep(PollingInterval * time.Second)
	}
}
