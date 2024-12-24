package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var currentAgentVersion semver = semver{Major: 0, Minor: 0, Patch: 2}
var currentUpdaterVersion semver = semver{Major: 0, Minor: 0, Patch: 2}

type semver struct {
	Major int
	Minor int
	Patch int
}

func (v semver) isOlderThan(sv semver) bool {
	if v.Major < sv.Major {
		return true
	}
	if v.Minor < sv.Minor {
		return true
	}
	if v.Patch < sv.Patch {
		return true
	}
	return false
}

func (v semver) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (d *serverDaemon) versionCheckHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {

	version := semver{}
	var err error

	version.Major, err = strconv.Atoi(params.ByName("Major"))
	if checkError(err) {
		return
	}

	version.Minor, err = strconv.Atoi(params.ByName("Minor"))
	if checkError(err) {
		return
	}

	version.Patch, err = strconv.Atoi(params.ByName("Patch"))
	if checkError(err) {
		return
	}

	switch params.ByName("App") {
	case "updater":
		if version.isOlderThan(currentUpdaterVersion) {
			w.WriteHeader(201)
		}
	case "agent":
		if version.isOlderThan(currentAgentVersion) {
			w.WriteHeader(201)
		}
	}

}

func (d *serverDaemon) buildAppHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	app := params.ByName("App")

	var commandText string
	var cmd *exec.Cmd
	switch app {
	case "agent":
		commandText = "../agent/build.sh"
	case "updater":
		commandText = "../updater/build.sh"
	}
	cmd = exec.Command(commandText)
	out, err := cmd.CombinedOutput()
	if checkError(err) {
		return
	}

	log.Printf("done building: %s", string(out))

}