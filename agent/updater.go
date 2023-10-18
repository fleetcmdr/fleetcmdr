package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kardianos/service"
)

type semver struct {
	Major int
	Minor int
	Patch int
}

func newDaemon() *agentDaemon {
	d := &agentDaemon{}
	d.hc.Timeout = time.Minute * 2
	d.hc.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	return d
}

func (d *agentDaemon) Start(s service.Service) error {
	err := s.Start()
	if checkError(err) {
		return err
	}

	return nil
}

func (d *agentDaemon) Stop(s service.Service) error {
	err := s.Stop()
	if checkError(err) {
		return err
	}

	return nil
}

func deployInstaller() {
	d := newDaemon()
	d.daemonCfg = getPlatformUpdaterConfig()
	var err error

	d.daemon, err = service.New(d, d.daemonCfg)
	if checkError(err) {
		return
	}

	err = d.daemon.Install()
	if checkError(err) {
		return
	}

	err = d.daemon.Start()
	if checkError(err) {
		return
	}

}

func (d *agentDaemon) downloaInstaller() (err error) {

	log.Printf("Attempting to download agent from '%s'", d.programUrl.String())
	resp, err := d.hc.Get(d.programUrl.String())
	if checkError(err) {
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if checkError(err) {
		return
	}
	defer resp.Body.Close()

	log.Printf("Received %s file", BytesToHuman(int64(len(bodyBytes))))

	f, err := os.Create(d.installPath)
	if checkError(err) {
		return
	}

	n, err := f.Write(bodyBytes)
	if checkError(err) {
		return
	}

	log.Printf("Wrote %s to file at '%s'", BytesToHuman(int64(n)), d.installPath)

	err = f.Sync()
	if checkError(err) {
		return
	}

	err = f.Close()
	if checkError(err) {
		return
	}

	return nil
}
