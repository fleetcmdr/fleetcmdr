package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type checkinData struct {
	ID      int
	Serial  string
	Version semver
}

type semver struct {
	Major int
	Minor int
	Patch int
}

func (v semver) string() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (d *serverDaemon) checkinHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	//var a agent

	var cd checkinData

	bodyBytes, err := io.ReadAll(r.Body)
	if checkError(err) {
		return
	}

	b := bytes.NewReader(bodyBytes)
	gd := gob.NewDecoder(b)

	err = gd.Decode(&cd)
	if checkError(err) {
		return
	}

	log.Printf("Agent (%d) with serial '%s' and version '%s'", cd.ID, cd.Serial, cd.Version.string())

	q := "INSERT INTO checkins (agent_id, v_major, v_minor, v_patch) VALUES (?,?,?,?)"
	_, err = d.db.ExecContext(context.Background(), q, cd.ID, cd.Version.Major, cd.Version.Minor, cd.Version.Patch)
	if checkError(err) {
		return
	}

}

func (d *serverDaemon) systemDataHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}
