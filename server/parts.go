package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type navItem struct {
	Name string
	ID   int
}

func (d *serverDaemon) leftNavHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

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

}
