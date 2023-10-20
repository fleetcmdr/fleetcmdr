package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

type serverDaemon struct {
	hc        http.Client
	hs        http.Server
	router    *httprouter.Router
	templates *template.Template
	db        *sql.DB
}

func parseTemplates() *template.Template {
	templ := template.New("")
	err := filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		_, err = templ.ParseFiles(path)
		if err != nil {
			log.Println(err)
		}

		return err
	})

	if err != nil {
		panic(err)
	}

	return templ
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	d := &serverDaemon{}

	d.db = InitializeMySQLDatabase("localhost", "fleetcmdr", os.Getenv("FLEETCMDR_MYSQL_USER"), os.Getenv("FLEETCMDR_MYSQL_PASS"))
	d.router = httprouter.New()
	d.templates = parseTemplates()

	d.hs = http.Server{
		Addr:    "localhost:2213",
		Handler: d.router,
	}

	log.Printf("Binding routes...")

	d.bindRoutes()

	log.Printf("Starting server...")
	err := d.hs.ListenAndServe()
	if checkError(err) {
		return
	}

	log.Printf("Shutting down...")

}
