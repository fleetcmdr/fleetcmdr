package main

import (
	"log"
	"net/http"
	"os"
)

func (d *serverDaemon) bindRoutes() {

	d.router.GET("/api/v1/parts/leftNav", d.leftNavHandler)
	d.router.GET("/api/v1/parts/agents/:id", d.viewAgentHandler)

	d.router.POST("/api/v1/checkin", d.checkinHandler)
	d.router.POST("/api/v1/systemData", d.systemDataHandler)

	d.router.NotFound = http.HandlerFunc(d.staticHandler)

}

func (d *serverDaemon) staticHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Handling static request '%s'", r.URL.Path)

	if r.URL.Path == "/" {
		r.URL.Path = "/static/index.html"

		var data struct{}

		err := d.templates.ExecuteTemplate(w, "index", data)
		if checkError(err) {
			return
		}

		return
	}

	f, err := os.Open(r.URL.Path[1:])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if checkError(err) {
		return
	}

	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)

}
