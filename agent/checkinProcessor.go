package main

import (
	"time"
)

func (d *agentDaemon) checkinProcessor() {

	t := time.NewTicker(time.Minute)
	systemDataCheckinTicker := time.NewTimer(time.Until(d.lastSystemDataCheckin.Add(time.Hour)))

	for {
		select {
		case <-t.C:
			d.checkin()
		case <-systemDataCheckinTicker.C:
			systemDataCheckinTicker = time.NewTimer(time.Until(time.Now().Add(time.Hour)))
			d.lastSystemDataCheckin = time.Now()
			d.sendSystemData()
		}

	}
}

type checkinData struct {
	ID int
}

type systemData struct {
	ID      int
	OS      string
	Payload any
}
