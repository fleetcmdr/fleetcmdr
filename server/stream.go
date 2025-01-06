package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Activity struct {
	PowerMetrics                 darwinPowerMetrics
	CPUConsumedPercent           float64
	MemoryPressurePercent        int64
	DiskIOOperationsPerSecond    int     // `ioutil -d`` unknown baseline
	DiskLatencyMilliseconds      float64 //  < 1 is good?
	DiskSizeBytes                int
	DiskUsedBytes                int
	NetworkUploadBytesPerSecond  int
	NetworkDownloadBytesPerSeond int
}

type criticality string

const (
	criticalityNominal  criticality = "nominal"
	criticalityWarning              = "warning"
	criticalityCritical             = "critical"
)

func (d *serverDaemon) agentStreamActivityReaderHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {

	agendIDStr := params.ByName("id")
	agentID, err := strconv.Atoi(agendIDStr)
	if checkError(err) {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	activityName := params.ByName("ActivityName")

	var act Activity

	a, err := d.getAgentByID(agentID)
	a.LatestActivityLocker.RLock()
	act = a.LatestActivity
	a.LatestActivityLocker.RUnlock()

	// log.Printf("sending activity: %#v", act)

	var rd struct {
		Value       any
		Criticality string
		Extra       any
		Text        string
	}

	switch activityName {
	case "battery":

		rd.Value = act.PowerMetrics.Battery.PercentCharge
		if act.PowerMetrics.Battery.PercentCharge < 20 {
			rd.Criticality = criticalityWarning
		}
		if act.PowerMetrics.Battery.PercentCharge < 5 {
			rd.Criticality = criticalityCritical
		}
		rd.Text = fmt.Sprintf("%d%%", act.PowerMetrics.Battery.PercentCharge)
	case "cpu":
		rd.Value = act.CPUConsumedPercent / (float64(a.CPUCountEfficiency + a.CPUCountPerformance))
		if act.CPUConsumedPercent/(float64(a.CPUCountEfficiency+a.CPUCountPerformance)) > 90 {
			rd.Criticality = criticalityWarning
		}
		if act.CPUConsumedPercent/(float64(a.CPUCountEfficiency+a.CPUCountPerformance)) > 95 {
			rd.Criticality = criticalityCritical
		}
		rd.Extra = act.PowerMetrics.Processor.Clusters
		rd.Text = fmt.Sprintf("%.1f", act.CPUConsumedPercent/(float64(a.CPUCountEfficiency+a.CPUCountPerformance)))
	case "ram":
		rd.Value = 100 - act.MemoryPressurePercent
		if 100-act.MemoryPressurePercent > 80 {
			rd.Criticality = criticalityWarning
		}
		if 100-act.MemoryPressurePercent > 90 {
			rd.Criticality = criticalityCritical
		}
		rd.Text = fmt.Sprintf("%d", 100-act.MemoryPressurePercent)
	case "disk":
		rd.Value = float64(act.DiskUsedBytes) / float64(act.DiskSizeBytes) * 100
		if act.DiskSizeBytes-act.DiskUsedBytes < 20000 { // 20GB remains
			rd.Criticality = criticalityWarning
		}
		if act.DiskSizeBytes-act.DiskUsedBytes < 10000 { // 10GB remains
			rd.Criticality = criticalityCritical
		}
		rd.Text = fmt.Sprintf("%.1fGB Remaining", float64(act.DiskSizeBytes-act.DiskUsedBytes)/1024)
	}

	jsonBytes, err := json.Marshal(rd)
	if checkError(err) {
		return
	}
	w.Write(jsonBytes)
	return
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

	// log.Printf("received activity: %#v", act)

	a, err := d.getAgentByID(agentID)
	if checkError(err) {
		return
	}
	a.LatestActivityLocker.Lock()
	a.LatestActivity = act
	a.LatestActivityLocker.Unlock()
	d.agentsLocker.Lock()
	d.agents[agentID] = a
	d.agentsLocker.Unlock()
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

	a, err := d.getAgentByID(agentID)
	if checkError(err) {
		return
	}
	d.agentsLocker.Lock()
	a.StreamingActivity = true
	d.agentsLocker.Unlock()
	d.agents[agentID] = a

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

	a, err := d.getAgentByID(agentID)
	if checkError(err) {
		return
	}
	d.agentsLocker.Lock()
	a.StreamingActivity = false
	d.agentsLocker.Unlock()

}
