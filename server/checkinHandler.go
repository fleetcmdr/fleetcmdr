package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (d *serverDaemon) checkinHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	var a agent

	bodyBytes, err := io.ReadAll(r.Body)
	if checkError(err) {
		return
	}

	b := bytes.NewReader(bodyBytes)
	gd := gob.NewDecoder(b)

	err = gd.Decode(a)
	if checkError(err) {
		return
	}

	q := "INSERT INTO checkins (id) VALUES (?)"
	_, err = d.db.ExecContext(context.Background(), q, a.ID)
	if checkError(err) {
		return
	}

}

func (d *serverDaemon) systemDataHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}
