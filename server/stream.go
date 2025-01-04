package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Activity struct {
	CPUConsumedPercent           float64
	MemoryPressurePercent        int64
	DiskIOOperationsPerSecond    int     // `ioutil -d`` unknown baseline
	DiskLatencyMilliseconds      float64 //  < 1 is good?
	DiskSizeBytes                int
	DiskUsedBytes                int
	NetworkUploadBytesPerSecond  int
	NetworkDownloadBytesPerSeond int
}

func (d *serverDaemon) agentStreamActivityReaderHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {

	agendIDStr := params.ByName("id")
	agentID, err := strconv.Atoi(agendIDStr)
	if checkError(err) {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	activityName := params.ByName("ActivityName")

	var a Activity

	d.agentsLocker.RLock()
	v := d.agents[agentID]
	d.agentsLocker.RUnlock()
	if v.LatestActivityLocker == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	v.LatestActivityLocker.RLock()
	a = v.LatestActivity
	v.LatestActivityLocker.RUnlock()

	switch activityName {
	case "cpu":
		w.Write([]byte(fmt.Sprintf("{\"cpu\": %.1f}", a.CPUConsumedPercent)))
		return
	}

	jsonBytes, err := json.Marshal(a)
	if checkError(err) {
		return
	}

	w.Write(jsonBytes)
}

func (d *serverDaemon) agentStreamActivityMomentReaderHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {

	agendIDStr := params.ByName("id")
	agentID, err := strconv.Atoi(agendIDStr)
	if checkError(err) {
		return
	}

	var a Activity

	d.agentsLocker.RLock()
	v := d.agents[agentID]
	d.agentsLocker.RUnlock()
	v.LatestActivityLocker.RLock()
	a = v.LatestActivity
	v.LatestActivityLocker.RUnlock()

	jsonBytes, err := json.Marshal(a)
	if checkError(err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func (d *serverDaemon) agentStreamActivityMomentHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {

	agendIDStr := params.ByName("id")
	agentID, err := strconv.Atoi(agendIDStr)
	if checkError(err) {
		return
	}

	var act Activity

	gob.Register(act)

	ge := gob.NewDecoder(req.Body)
	err = ge.Decode(&act)
	if checkError(err) {
		return
	}

	log.Printf("received activity: %#v", act)

	d.agentsLocker.Lock()
	a, ok := d.agents[agentID]
	d.agentsLocker.Unlock()
	if !ok {
		a, err = d.getAgentByID(agentID)
		if checkError(err) {
			return
		}
	}
	a.LatestActivityLocker.Lock()
	a.LatestActivity = act
	a.LatestActivityLocker.Unlock()
	if !a.StreamingActivity {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (d *serverDaemon) agentStartStreamActivityHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// we want to get agent streaming updates as soon as possible

	agendIDStr := params.ByName("id")
	agentID, err := strconv.Atoi(agendIDStr)
	if checkError(err) {
		return
	}

	q := "UPDATE agents SET streaming_activity=true WHERE id = $1"
	_, err = d.db.ExecContext(context.Background(), q, agentID)
	if checkError(err) {
		return
	}

	d.agentsLocker.Lock()
	v := d.agents[agentID]
	d.agentsLocker.Unlock()
	v.StreamingActivity = true
	d.agents[agentID] = v

}

func (d *serverDaemon) agentEndStreamActivityHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// shut down the stream

	agendIDStr := params.ByName("id")
	agentID, err := strconv.Atoi(agendIDStr)
	if checkError(err) {
		return
	}

	q := "UPDATE agents SET streaming_activity=false WHERE id = $1"
	_, err = d.db.ExecContext(context.Background(), q, agentID)
	if checkError(err) {
		return
	}

	d.agentsLocker.Lock()
	v := d.agents[agentID]
	d.agentsLocker.Unlock()
	v.StreamingActivity = false
	d.agents[agentID] = v

}
