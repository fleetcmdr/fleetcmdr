package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (d *serverDaemon) bindRoutes() {

	d.router.GET("/", d.baseHandler)
	d.router.GET("/static/*path", d.staticHandler)

	d.router.GET("/api/v1/parts/leftNav", d.leftNavHandler)
	d.router.GET("/api/v1/parts/agent/:id", d.viewAgentHandler)
	d.router.GET("/api/v1/parts/commands/:id", d.commandHistoryForAgentHandler)

	d.router.POST("/api/v1/sendCommand/:id", d.sendCommandHandler)
	d.router.POST("/api/v1/checkin", d.checkinHandler)
	d.router.POST("/api/v1/sendSystemData", d.systemDataHandler)

	d.router.POST("/api/v1/sendCommandResult", d.commandResultHandler)

}

func (d *serverDaemon) baseHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	log.Printf("Returning index")

	var data struct{}

	err := d.templates.ExecuteTemplate(w, "index", data)
	if checkError(err) {
		return
	}
}

func (d *serverDaemon) staticHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	fName := fmt.Sprintf("./static%s", params.ByName("path"))

	log.Printf("Handling static request '%s'", fName)

	fileBytes, err := os.ReadFile(fName)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.HasSuffix(fName, "js") {
		w.Header().Set("Content-Type", "application/javascript")
	}
	if strings.HasSuffix(fName, "css") {
		w.Header().Set("Content-Type", "text/css")
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))

	w.Write(fileBytes)
}
