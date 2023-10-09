package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	PercentageWarning = 30
	PercentageEmpty   = 10
	PollingInterval   = 3
)

// setLogger and logger level
func setLogger() {
	isVerbose := flag.Bool("verbose", false, "display debug logs")
	flag.Parse()

	if *isVerbose {
		log.SetLevel(log.DebugLevel)
	}

	log.SetOutput(os.Stdout)
}

func main() {
	setLogger()

	log.Debug("hi")
}
