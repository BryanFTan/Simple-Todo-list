package main

import (
	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//App - main app
type App struct {
	Router *mux.Router
	DB 	   *sql.DB
}


