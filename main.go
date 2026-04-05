package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	pollingInterval             = 5
	secondsBetweenNotifications = 60 * 5

	notificationTimeout = "5000"

	batteryFullPNG  = "icons/battery-full.png"
	batteryHalfPNG  = "icons/battery-half.png"
	batteryEmptyPNG = "icons/battery-empty.png"
)

type deviceStatus struct {
	batteryStatus         int
	step                  int
	notificationTimestamp int64
}

type HeadsetControlOutput struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	APIVersion    string `json:"api_version"`
	HidapiVersion string `json:"hidapi_version"`
	DeviceCount   int    `json:"device_count"`
	Devices       []struct {
		Status          string   `json:"status"`
		Device          string   `json:"device"`
		Vendor          string   `json:"vendor"`
		Product         string   `json:"product"`
		IDVendor        string   `json:"id_vendor"`
		IDProduct       string   `json:"id_product"`
		Capabilities    []string `json:"capabilities"`
		CapabilitiesStr []string `json:"capabilities_str"`
		Battery         struct {
			Status         string `json:"status"`
			Level          int    `json:"level"`
			VoltageMv      int    `json:"voltage_mv"`
			TimeToEmptyMin int    `json:"time_to_empty_min"`
		} `json:"battery"`
	} `json:"devices"`
}

type monitor struct {
	devices    map[string]deviceStatus
	execCmd    func(name string, arg ...string) IShellCommand
	mouseName  string
	workingDir string
}

func newMonitor(mouseName, workingDir string) *monitor {
	return &monitor{
		devices:    make(map[string]deviceStatus),
		execCmd:    shellCommandWrapper,
		mouseName:  mouseName,
		workingDir: workingDir,
	}
}

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

func (m *monitor) checkHeadsetBattery() {
	output, err := m.execCmd("headsetcontrol", "-b", "-o", "json").Output()
	if err != nil {
		slog.Debug("Headset not found.")
		return
	}

	slog.Debug(string(output))

	jsonOutput := &HeadsetControlOutput{}
	if err := json.Unmarshal(output, jsonOutput); err != nil {
		slog.Debug("Failed to parse headsetcontrol JSON output", "error", err)
		return
	}

	if jsonOutput.DeviceCount == 0 {
		slog.Debug("No headset connected.")
		return
	}

	headsetName := "🎧️ " + jsonOutput.Devices[0].Device
	batteryStatus := jsonOutput.Devices[0].Battery.Level
	m.handleBatteryStatus(headsetName, batteryStatus)
}

func (m *monitor) handleBatteryStatus(deviceName string, batteryStatus int) {
	step := batteryStatus / 5

	slog.Debug("Device name: " + deviceName)
	slog.Debug("Battery status: " + strconv.Itoa(batteryStatus))
	slog.Debug("Step: " + strconv.Itoa(step))

	previousStatus, exists := m.devices[deviceName]
	sec := time.Now().Unix()

	if !exists || (previousStatus.notificationTimestamp+secondsBetweenNotifications < sec &&
		batteryStatus/5 != previousStatus.step) {
		m.devices[deviceName] = deviceStatus{
			batteryStatus:         batteryStatus,
			step:                  step,
			notificationTimestamp: sec,
		}
		m.sendNotification(deviceName, batteryStatus)
	}
}

func (m *monitor) checkMouseBattery() {
	output, err := m.execCmd("solaar", "show").Output()
	if err != nil {
		slog.Debug("Solaar not found.")
		return
	}

	batteryStatus, err := getMouseBattery(string(output), m.mouseName)
	if err != nil {
		slog.Debug("Mouse not found.")
		return
	}

	deviceName := "🖱️ " + m.mouseName
	m.handleBatteryStatus(deviceName, batteryStatus)
}

func (m *monitor) sendNotification(deviceName string, batteryStatus int) {
	slog.Debug("Sending notification", "device", deviceName, "battery", batteryStatus)
	message := deviceName + " - Battery at " + strconv.Itoa(batteryStatus) + "%"

	icon := selectIcon(m.workingDir, batteryStatus)
	cmd := m.execCmd(
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

func getMouseBattery(input string, mouseName string) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	inLightspeedSection := false
	inMouseSection := false
	batteryRegex := regexp.MustCompile(`Battery: (\d+)%`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

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

func main() {
	setLogger()

	mouseName := flag.String("mouse", "PRO X Wireless", "mouse device name")
	flag.Parse()

	workingDir, err := os.Getwd()
	if err != nil {
		slog.Error("Failed to get working directory", "error", err)
		os.Exit(1)
	}

	m := newMonitor(*mouseName, workingDir)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	slog.Info("Starting battery monitor", "mouse", *mouseName)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down battery monitor")
			return
		default:
			slog.Debug("Polling...")
			m.checkHeadsetBattery()
			m.checkMouseBattery()
			time.Sleep(pollingInterval * time.Second)
		}
	}
}
