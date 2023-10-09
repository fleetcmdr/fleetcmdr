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

func checkError(err error) bool {
	if err != nil {
		PrintErrorN(err, 1)
		//log.Print(err)
		return true
	}
	return false
}

type service struct {
	hc        http.Client
	router    *httprouter.Router
	templates *template.Template
	db        *sql.DB
}

var svc *service

func parseTemplates() *template.Template {
	templ := template.New("")
	err := filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
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

	svc.db = InitializeMySQLDatabase("localhost", "fleetcmdr", "root", os.Getenv("FLEETCMDR_MYSQL_ROOT_PASS"))
	svc.router = httprouter.New()
	svc.templates = parseTemplates()

}
