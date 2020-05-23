package main

import (
	"fmt"
	"database/sql"
	"log"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//App - main app
type App struct {
	Router *mux.Router
	DB 	   *sql.DB
}

//InitialiseDB - method to initialise Postgres DB
func (a *App) InitialiseDB(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
}

//Run - run app
func (a *App) Run(addr string) {}
