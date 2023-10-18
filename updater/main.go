package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/kardianos/service"
)

type updaterDaemon struct {
	daemonCfg   *service.Config
	daemon      service.Service
	hc          http.Client
	programUrl  url.URL
	installPath string
	version     semver
}

type semver struct {
	Major int
	Minor int
	Patch int
}

func newDaemon() *updaterDaemon {
	d := &updaterDaemon{}
	d.hc.Timeout = time.Minute * 2
	d.hc.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	return d
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// largely just going to sit around and wait until a newer agent is available
	// checking every 24 hours for new agent
	// localhost listener allows agent to poke and perform on-demand agent updates
	d := newDaemon()
	d.programUrl.Scheme = "http://"
	d.programUrl.Host = "localhost:2213"
	d.programUrl.Path = fmt.Sprintf("/static/agent/%s/fc_agent", runtime.GOOS)
	d.daemonCfg = getPlatformAgentConfig()
	var err error
	d.daemon, err = service.New(d, d.daemonCfg)
	if checkError(err) {
		return
	}

	if service.Interactive() {
		err = d.daemon.Install()
		if checkError(err) {
			return
		}

		err = d.daemon.Start()
		if checkError(err) {
			return
		}
	} else {
		t := time.NewTicker(time.Hour * 25)
		for {
			select {
			case <-t.C:
				err = d.checkForUpdates()
				if checkError(err) {
					//return
				}
			}
		}
	}

}

func (d *updaterDaemon) Start(s service.Service) error {
	return nil
}

func (d *updaterDaemon) Stop(s service.Service) error {
	return nil
}

func (d *updaterDaemon) downloadAgent() (err error) {

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

func (d *updaterDaemon) installAgent() (err error) {
	err = d.daemon.Install()
	if checkError(err) {
		if errors.Is(err, service.ErrNotInstalled) {

		}
		return
	}

	err = d.daemon.Start()
	if checkError(err) {
		return
	}
	return
}

func (d *updaterDaemon) checkForUpdates() (err error) {

	maxRetryWait, err := time.ParseDuration("8h")
	if checkError(err) {
		return
	}
	retryWait, err := time.ParseDuration("10s")
	if checkError(err) {
		return
	}
	retryStage := int64(0)

	for {
		log.Printf("Retrying download of agent")
		err = d.downloadAgent()
		if err == nil {
			log.Printf("Agent download successful")
			err = d.uninstallAgent()
			if checkError(err) {
				return
			}

			err = d.installAgent()
			if checkError(err) {
				return
			}
			return
		}
		retryStage++

		// exponential backoff
		retryWait = time.Duration(retryWait.Nanoseconds() ^ retryStage)
		if retryWait.Nanoseconds() > maxRetryWait.Nanoseconds() {
			retryWait = maxRetryWait
		}
		log.Printf("Waiting %s to retry", retryWait.String())
		time.Sleep(retryWait)
	}

}

func (d *updaterDaemon) uninstallAgent() (err error) {

	err = d.daemon.Stop()
	if checkError(err) {
		//return
	}

	err = d.daemon.Uninstall()
	if checkError(err) {
		return
	}
	return
}
