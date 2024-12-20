package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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

		// parse out cpu count
		procString := sData.SystemData.SPHardwareDataType[0].NumberProcessors

		procString = strings.ReplaceAll(procString, "proc ", "")
		procCountStrings := strings.Split(procString, ":")
		sData.AgentData.CPUCountEfficiency, err = strconv.Atoi(procCountStrings[2])
		if checkError(err) {
			return
		}

		sData.AgentData.CPUCountPerformance, err = strconv.Atoi(procCountStrings[1])
		if checkError(err) {
			return
		}

		cs, err := d.getAgentCommands(id)
		if checkError(err) {
			return
		}

		sData.Commands = cs

		scripts, err := d.getScriptsForAgent(a.ID)
		if checkError(err) {
			return
		}

		sData.Scripts = scripts

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

func (d *serverDaemon) sendCommandHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {

	req.ParseForm()

	id, err := strconv.Atoi(params.ByName("id"))
	if checkError(err) {
		return
	}

	input := req.Form.Get("input")

	cUUID, err := uuid.NewUUID()
	if checkError(err) {
		return
	}

	log.Printf("Recieved command '%s' from agent %d", input, id)

	var agentID int

	q := "INSERT INTO commands (agent_id, input, c_uuid, scheduled_ts) VALUES ($1, $2, $3, NOW()) RETURNING id"
	err = d.db.QueryRowContext(context.Background(), q, id, input, cUUID.String()).Scan(&agentID)
	if checkError(err) {
		return
	}

	time.Sleep(time.Millisecond * 200)

	params = append(params, httprouter.Param{Key: "agentID", Value: strconv.Itoa(id)})

	d.commandHistoryForAgentHandler(w, req, params)
}

func (d *serverDaemon) commandOutputRefreshHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {

	id, err := strconv.Atoi(params.ByName("commandID"))
	if checkError(err) {
		return
	}

	var co Command
	q := "SELECT id, COALESCE(output,''), executed_ts FROM commands WHERE id = $1"
	err = d.db.QueryRowContext(context.Background(), q, id).Scan(&co.ID, &co.Output, &co.ExecutedTS)
	if checkError(err) {
		return
	}

	if co.ExecutedTS.Valid {
		co.Executed = true
	}

	b := bytes.NewBuffer(nil)
	err = d.templates.ExecuteTemplate(b, "command-output", co)
	if checkError(err) {
		return
	}

	responseBytes := b.Bytes()

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseBytes)))
	w.Write(responseBytes)

}

type Script struct {
	ID          int
	Name        string
	Created     time.Time
	Modified    time.Time
	Notes       string
	Creator     int
	Description string
}

type ScriptParameters struct {
	ScriptID  int
	ID        int
	Name      string
	ValueType string
}

func (d *serverDaemon) getScriptsForAgent(id int) ([]Script, error) {

	q := "SELECT id, name, created_ts, modified_ts, notes, creator_id, description FROM script_library WHERE os = (SELECT os FROM agents WHERE id = $1)"
	rows, err := d.db.QueryContext(context.Background(), q, id)
	if checkError(err) {
		return nil, err
	}

	var scripts []Script

	for rows.Next() {
		s := Script{}
		err = rows.Scan(&s.ID, &s.Name, &s.Created, &s.Modified, &s.Notes, &s.Creator, &s.Description)
	}

	return scripts, nil

}

func (d *serverDaemon) getAgentCommands(id int) (commands []Command, err error) {
	q := "SELECT id, ts, input, COALESCE(output,''), scheduled_ts, delivered_ts, executed_ts FROM commands WHERE agent_id = $1 ORDER BY scheduled_ts DESC LIMIT 20"
	rows, err := d.db.QueryContext(context.Background(), q, id)
	if checkError(err) {
		return
	}

	cmds := []Command{}

	for rows.Next() {
		c := Command{}
		err = rows.Scan(&c.ID, &c.TS, &c.Input, &c.Output, &c.ScheduledTS, &c.DeliveredTS, &c.ExecutedTS)
		if checkError(err) {
			return
		}

		if c.ExecutedTS.Valid {
			c.Executed = true
		}

		cmds = append([]Command{c}, cmds...)
	}

	return cmds, nil
}

type Command struct {
	ID          int
	UUID        string
	TS          time.Time
	Input       string
	Output      string
	ScheduledTS sql.NullTime
	DeliveredTS sql.NullTime
	ExecutedTS  sql.NullTime
	Executed    bool
}

type darwinSystemData struct {
	AgentData  *agent
	SystemData AppleSystemProfilerOutput
	Commands   []Command
	Scripts    []Script
}
