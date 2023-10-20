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
	"path/filepath"
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

type agentDaemon struct {
	ID                    int
	daemonCfg             *service.Config
	daemon                service.Service
	hc                    http.Client
	programUrl            url.URL
	installPath           string
	version               semver
	debug                 bool
	cmdr                  string
	lastSystemDataCheckin time.Time
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
	ud := newDaemon()
	ud.programUrl.Scheme = "http"
	ud.programUrl.Host = "localhost:2213"
	ud.programUrl.Path = fmt.Sprintf("/static/downloads/updater/%s/fc_updater", runtime.GOOS)
	ud.daemonCfg = getPlatformUpdaterConfig()
	var err error
	ud.daemon, err = service.New(ud, ud.daemonCfg)
	if checkError(err) {
		return
	}

	ad := &agentDaemon{}
	ad.programUrl.Scheme = "http"
	ad.programUrl.Host = "localhost:2213"
	ad.programUrl.Path = fmt.Sprintf("/static/downloads/agent/%s/fc_agent", runtime.GOOS)
	ad.daemonCfg = getPlatformAgentConfig()
	ad.daemon, err = service.New(ad, ad.daemonCfg)
	if checkError(err) {
		return
	}

	if service.Interactive() {
		err = ad.daemon.Stop()
		if checkError(err) {
			//return
		}

		err = ad.daemon.Uninstall()
		if checkError(err) {
			//return
		}

		err = download(ad.programUrl.String(), ad.daemonCfg.Executable)
		if checkError(err) {
			//return
		}

		err = ad.daemon.Install()
		if checkError(err) {
			//return
		}

		err = ad.daemon.Start()
		if checkError(err) {
			return
		}
		log.Printf("Service started")
		return
	}

	t := time.NewTicker(time.Hour * 25)
	for {
		select {
		case <-t.C:
			err = ad.checkForUpdates()
			if checkError(err) {
				//return
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

func (d *agentDaemon) Start(s service.Service) error {
	return nil
}

func (d *agentDaemon) Stop(s service.Service) error {
	return nil
}

func download(endpoint, toPath string) (err error) {

	hc := &http.Client{Timeout: time.Minute * 2}

	log.Printf("Attempting to download from '%s'", endpoint)
	resp, err := hc.Get(endpoint)
	if checkError(err) {
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if checkError(err) {
		return
	}
	defer resp.Body.Close()

	log.Printf("Received %s file", BytesToHuman(int64(len(bodyBytes))))

	err = os.MkdirAll(filepath.Dir(toPath), 0700)
	if checkError(err) {
		return
	}

	f, err := os.OpenFile(toPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0700)
	if checkError(err) {
		return
	}

	n, err := f.Write(bodyBytes)
	if checkError(err) {
		return
	}

	log.Printf("Wrote %s to file at '%s'", BytesToHuman(int64(n)), toPath)

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

func (d *agentDaemon) checkForUpdates() (err error) {

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
		err = download(d.programUrl.String(), d.daemonCfg.Executable)
		if err == nil {
			log.Printf("Agent download successful")

			err = d.daemon.Stop()
			if checkError(err) {
				//return
			}

			err = d.daemon.Uninstall()
			if checkError(err) {
				//return
			}

			err = d.daemon.Install()
			if checkError(err) {
				//return
			}

			err = d.daemon.Start()
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
