package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/kardianos/service"
)

const (
	versionMajor = 0
	versionMinor = 0
	versionPatch = 1
)

func (v semver) string() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

type agentDaemon struct {
	ID                    int
	hostname              string
	daemonCfg             *service.Config
	daemon                service.Service
	hc                    http.Client
	programUrl            url.URL
	installPath           string
	version               semver
	debug                 bool
	cmdHost               string
	lastSystemDataCheckin time.Time
	systemData            any
	commandChan           chan Command
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Agent starting...")

	cmdHost := "localhost"
	d := newDaemon()
	d.daemonCfg = getPlatformAgentConfig()
	d.cmdHost = fmt.Sprintf("http://%s:2213", cmdHost)
	d.commandChan = make(chan Command, 50)

	go d.commandProcessor()

	if service.Interactive() {
		d.debug = true
		d.runAgent()
	} else {
		d.runAgent()
	}

}

func (d *agentDaemon) runAgent() {
	log.Printf("Agent running")
	d.checkinProcessor()
}
