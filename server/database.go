package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

func InitializeMSSQLDatabase(host, name, user, pass, appname string) *sql.DB {

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(user, pass),
		Host:   fmt.Sprintf("%s:1433", host),
	}

	var db *sql.DB
	var err error

	if appname != "" {
		db, err = sql.Open("sqlserver", fmt.Sprintf("%s://%s@%s/?database=%s&app%%20name=%s", u.Scheme, u.User.String(), u.Host, name, appname))
	} else {
		db, err = sql.Open("sqlserver", fmt.Sprintf("%s://%s@%s/?database=%s", u.Scheme, u.User.String(), u.Host, name))
	}

	if err != nil {
		log.Printf("Failed to connect to db %s, returned error: %v", name, err)
		return nil
	}
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(100)
	return db
}

func InitializeMySQLDatabase(host, name, user, pass string) *sql.DB {

	if !strings.Contains(host, ":") {
		host = host + ":3306"
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&tls=false", user, pass, host, name))

	if err != nil {
		log.Printf("Failed to connect to db %s, returned error: %v", name, err)
		return nil
	}
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(100)
	return db
}
