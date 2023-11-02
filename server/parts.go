package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type navItem struct {
	Name string
	ID   int
}

func (d *serverDaemon) leftNavHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	log.Printf("leftNav requested")

	var navItems []navItem

	agents := d.getAgents(50, 0)

	for _, a := range agents {
		var ni navItem
		ni.ID = a.ID
		ni.Name = a.Name
		navItems = append(navItems, ni)
	}

	err := d.templates.ExecuteTemplate(w, "leftNav", navItems)
	if checkError(err) {
		return
	}
}

func (d *serverDaemon) viewAgentHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	var err error
	id, err := strconv.Atoi(params.ByName("id"))
	if checkError(err) {
		return
	}

	log.Printf("agent %d requested", id)

	a, err := d.getAgentByID(id)
	if checkError(err) {
		return
	}

	switch a.OS {
	case "darwin":
		sData := darwinSystemData{}
		sData.AgentData = a

		err = json.Unmarshal([]byte(sData.AgentData.SystemData), &sData.SystemData)
		if checkError(err) {
			return
		}

		log.Printf(sData.SystemData.SPHardwareDataType[0].SerialNumber)

		// sData.systemData = a.SystemData.(AppleSystemProfilerOutput)
		b := bytes.NewBuffer(nil)
		err = d.templates.ExecuteTemplate(b, "agent-darwin", sData)
		if checkError(err) {
			return
		}

		responseBytes := b.Bytes()

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(responseBytes)))
		w.Write(responseBytes)

	case "windows":
	}

}

type darwinSystemData struct {
	AgentData  *agent
	SystemData AppleSystemProfilerOutput
}
