package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
)

type serverDaemon struct {
	hc                        http.Client
	hs                        http.Server
	router                    *httprouter.Router
	templates                 *template.Template
	db                        *sql.DB
	currentAgentVersion       semver
	currentAgentVersionLocker sync.RWMutex

	agentsLocker sync.RWMutex
	agents       map[int]agent
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

	// d.db = InitializeMySQLDatabase("localhost", "fleetcmdr", os.Getenv("FLEETCMDR_MYSQL_USER"), os.Getenv("FLEETCMDR_MYSQL_PASS"))
	d.db = InitializePgSQLDatabase("localhost", "fleetcmdr", os.Getenv("FLEETCMDR_PGSQL_USER"), os.Getenv("FLEETCMDR_PGSQL_PASS"))
	d.router = httprouter.New()
	d.templates = parseTemplates()
	d.agents = make(map[int]agent)

	d.hs = http.Server{
		Addr:    "localhost:2213",
		Handler: d.router,
	}
	d.getLatestAgentVersion()

	log.Printf("Binding routes...")

	d.bindRoutes()

	log.Printf("Starting server...")
	err := d.hs.ListenAndServe()
	if checkError(err) {
		return
	}

	log.Printf("Shutting down...")

}
