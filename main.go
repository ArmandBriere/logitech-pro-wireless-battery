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

	BATTERY_FULL_PNG  = "./icons/battery-full.png"
	BATTERY_HALF_PNG  = "./icons/battery-half.png"
	BATTERY_EMPTY_PNG = "./icons/battery-empty.png"
)

var headsetName string
var batteryStatus int

// setLogger and logger level
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

	headsetName = strings.Split(headsetFoundSplit[1], "!")[0]
	headsetBatterySplit := strings.Split(string(output), "Battery: ")
	if len(headsetBatterySplit) == 1 {
		log.Debug("Headset usb is connected but the headset is in sleep mode.")
		return
	}

	batteryStatus, _ = strconv.Atoi(strings.Split(headsetBatterySplit[1], "%")[0])
	log.Debug("Headset name: " + headsetName)
	log.Debug("Battery status: " + strconv.Itoa(batteryStatus))

	sendNotification()
}

func sendNotification() {

	message := headsetName + " - Battery at " + strconv.Itoa(batteryStatus) + "%"

	_, err := exec.Command("/usr/bin/notify-send", "-t", notificationTimeout, message).Output()
	if err != nil {
		log.Error(err.Error())
	}
}

func main() {
	setLogger()

	for {
		execHeadsetcontrol()
		time.Sleep(PollingInterval * time.Second)
	}
}
