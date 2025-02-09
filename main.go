package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	pollingInterval           = 5
	secondBetweenNotification = 60 * 5

	notificationTimeout = "5000"

	batteryFullPNG  = "icons/battery-full.png"
	batteryHalfPNG  = "icons/battery-half.png"
	batteryEmptyPNG = "icons/battery-empty.png"
)

type status struct {
	deviceName            string
	batteryStatus         int
	step                  int
	notificationTimestamp int64
}

var path string
var devices map[string]status = make(map[string]status)

// Set logger and logger level
func setLogger() {
	isVerbose := flag.Bool("verbose", false, "display debug logs")
	flag.Parse()

	var logLevel slog.LevelVar
	logLevel.Set(slog.LevelInfo)

	if *isVerbose {
		logLevel.Set(slog.LevelDebug)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: &logLevel,
	}))

	slog.SetDefault(logger)
}

// checkHeadsetBattery check the battery status of the headset
func checkHeadsetBattery() {
	output, err := exec.Command("headsetcontrol", "-b").Output()
	if err != nil {
		slog.Debug("Headset not found.")
		return
	}

	slog.Debug(string(output))

	headsetFoundSplit := strings.Split(string(output), "Found ")

	if len(strings.Split(string(output), "Found ")) == 1 {
		slog.Debug("Headset not found.")
		return
	}

	headsetName := strings.Split(headsetFoundSplit[1], "!")[0]
	headsetBatterySplit := strings.Split(string(output), "Battery: ")
	if len(headsetBatterySplit) == 1 {
		slog.Debug("Headset usb is connected but the headset is in sleep mode.")
		return
	}

	batteryStatus, _ := strconv.Atoi(strings.Split(headsetBatterySplit[1], "%")[0])

	headsetName = "🎧️ " + headsetName
	handleBatteryStatus(headsetName, batteryStatus)
}

// handleBatteryStatus decides to send notification or not based on battery level
func handleBatteryStatus(deviceName string, batteryStatus int) {
	step := batteryStatus / 5

	slog.Debug("Device name: " + deviceName)
	slog.Debug("Battery status: " + strconv.Itoa(batteryStatus))
	slog.Debug("Step: " + strconv.Itoa(step))

	previousStatus, exists := devices[deviceName]
	sec := time.Now().Unix()

	if !exists || (previousStatus.notificationTimestamp+secondBetweenNotification < sec &&
		batteryStatus/5 != previousStatus.step) {
		devices[deviceName] = status{
			deviceName:            deviceName,
			batteryStatus:         batteryStatus,
			step:                  step,
			notificationTimestamp: sec,
		}
		sendNotification(devices[deviceName])
	}
}

// checkMouseBattery check the battery status of the mouse
func checkMouseBattery() {
	output, err := exec.Command("solaar", "show").Output()
	if err != nil {
		slog.Debug("Solaar not found.")
		return
	}

	deviceName := "PRO X Wireless"
	batteryStatus, err := getMouseBattery(string(output), deviceName)
	if err != nil {
		slog.Debug("Mouse not found.")
		return
	}

	deviceName = "🖱️ " + deviceName
	handleBatteryStatus(deviceName, batteryStatus)
}

// getMouseBattery extracts the battery percentage of the mouse
func getMouseBattery(input string, mouseName string) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	inLightspeedSection := false
	inMouseSection := false
	batteryRegex := regexp.MustCompile(`Battery: (\d+)%`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Detect sections
		if line == "Lightspeed Receiver" {
			inLightspeedSection = true
			continue
		}

		if inLightspeedSection && strings.HasSuffix(line, mouseName) {
			inMouseSection = true
			continue
		}

		if inLightspeedSection && inMouseSection {
			if match := batteryRegex.FindStringSubmatch(line); match != nil {
				slog.Debug("Mouse battery: " + match[1])
				return strconv.Atoi(match[1])
			}
		}
	}

	slog.Error("mouse battery not found")
	return 0, fmt.Errorf("mouse battery not found")
}

// Send notification using notify-send
func sendNotification(status status) {
	slog.Debug("Sending notification", "status", status)
	message := status.deviceName + " - Battery at " + strconv.Itoa(status.batteryStatus) + "%"

	icon := selectIconPng(status)
	cmd := shellCommand(
		"notify-send",
		"-t", notificationTimeout,
		"-i", icon,
		message,
	)

	data, err := cmd.Output()
	slog.Debug("Notification sent", "data", string(data))
	if err != nil {
		slog.Debug("Error sending notification")
		slog.Error(err.Error())
	}
}

func main() {
	setLogger()
	path, _ = os.Getwd()
	for {
		slog.Debug("Polling...")
		checkHeadsetBattery()
		checkMouseBattery()
		time.Sleep(pollingInterval * time.Second)
	}
}
