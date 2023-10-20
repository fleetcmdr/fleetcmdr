package main

import (
	"log"
	"time"
)

const (
	checkinURL    string = "/api/v1/checkin"
	systemDataURL        = "/api/v1/sendSystemData"
)

func (d *agentDaemon) checkinProcessor() {

	t := time.NewTicker(time.Minute)
	systemDataCheckinTicker := time.NewTimer(time.Until(d.lastSystemDataCheckin.Add(time.Hour)))

	for {
		select {
		case <-t.C:
			//log.Printf("Checking in with id (%d) and serial '%s', with version %s", d.ID, d.getSystemData().SPHardwareDataType[0].SerialNumber, d.version.string())
			d.checkin()
		case <-systemDataCheckinTicker.C:
			log.Printf("Sending in system data...")
			systemDataCheckinTicker = time.NewTimer(time.Until(time.Now().Add(time.Hour)))
			d.lastSystemDataCheckin = time.Now()
			go d.sendSystemData()
		}

	}
}

type checkinData struct {
	ID      int
	Serial  string
	Version semver
}

type systemData struct {
	ID      int
	OS      string
	Payload any
}
