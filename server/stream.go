package main

import (
	"context"
	"encoding/gob"
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

func (d *serverDaemon) agentStreamActivityMomentHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {

	var a Activity

	gob.Register(a)

	var err error
	ge := gob.NewDecoder(req.Body)
	err = ge.Decode(&a)
	if checkError(err) {
		return
	}

	log.Printf("received activity: %#v", a)
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
	v.StreamingActivity = true
	d.agents[agentID] = v
	d.agentsLocker.Unlock()

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
	v.StreamingActivity = false
	d.agents[agentID] = v
	d.agentsLocker.Unlock()

}
