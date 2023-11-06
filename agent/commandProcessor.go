package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

func (d *agentDaemon) commandProcessor() {
	var err error

	for {
		select {
		case c := <-d.commandChan:
			c.Output, err = run(c.Input)
			if checkError(err) {
				// Do something smart?
			}

			d.returnCommandResult(c)
		}
	}
}

func (d *agentDaemon) returnCommandResult(c Command) {

	b := &bytes.Buffer{}

	gob.Register(c)

	ge := gob.NewEncoder(b)
	err := ge.Encode(c)
	if checkError(err) {
		return
	}

	log.Printf("Returning command result for UUID %s: '%s'", c.UUID, c.Output)

	_, err = d.hc.Post(fmt.Sprintf("%s/%s", d.cmdHost, commandResultURL), "application/octet-stream", b)
	if checkError(err) {
		return
	}
}
