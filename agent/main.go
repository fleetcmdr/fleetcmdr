package main

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/kardianos/service"
)

type daemon struct {
	ID                    int
	daemonCfg             *service.Config
	service               service.Service
	hc                    http.Client
	programUrl            url.URL
	installPath           string
	version               semver
	debug                 bool
	cmdr                  string
	lastSystemDataCheckin time.Time
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	d := newDaemon()
	d.daemonCfg = getPlatformAgentConfig()
	d.cmdr = "http://localhost:2213"

	if service.Interactive() {
		d.debug = true
	} else {
		d.runAgent()
	}

}

func (d *daemon) runAgent() {

	d.checkinProcessor()
}
