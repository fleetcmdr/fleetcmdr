package main

import (
	"log"
	"os"
	"testing"

	"howett.net/plist"
)

func TestPowerMetricsParse(T *testing.T) {

	plistBytes, err := os.ReadFile("/Users/bporter/Downloads/powermetrics-output.plist")
	if checkError(err) {
		return
	}

	var pm darwinPowerMetrics

	f, err := plist.Unmarshal(plistBytes, &pm)
	if checkError(err) {
		return
	}

	log.Printf("Format %d returned %#v", f, pm)

}
