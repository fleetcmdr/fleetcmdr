package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type navItem struct {
	Name string
	ID   int
}

func (svc *service) leftNavHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	var navItems []navItem

	agents := svc.getAgents(50, 0)

	for _, a := range agents {
		var ni navItem
		ni.ID = a.ID
		ni.Name = a.Name
		navItems = append(navItems, ni)
	}

	err := svc.templates.ExecuteTemplate(w, "leftNav", navItems)
	if checkError(err) {
		return
	}
}

func (svc *service) viewAgentHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}
