package main

import (
	"fmt"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"encoding/json"

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
	a.initializeRoutes()
}

//Run - run app
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}



//-- Route Handlers --
/*
- handler retrieves id to be fected from requested URL
- uses GetListItem to fetch details of that item
- if item not found -> return 404
*/
func (a *App) getListItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if (err != nil) {
		respondWithError(w, http.StatusBadRequest, "Invalid listItem ID")
		return
	}

	l := listItem{ID: id}
	if err := l.getListItem(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, l)
}

/*
- handler fetches all list items
- use count & start from query string fetch count of items starting from position start in DB
*/
func (a *App) getAllListItems(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	listItems, err := getAllListItems(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, listItems)
}

/*
- handler creates product
*/
func (a *App) createListItem(w http.ResponseWriter, r *http.Request) {
	var l listItem
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&l); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := l.createListItem(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, l)
}

/*
- handler to update listitem
*/
func (a *App) updateListItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid listItem ID")
        return
	}
	
	var l listItem
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&l); err != nil {
		respondWithError(w, http.StatusBadGateway, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	l.ID = id

	if err := l.updateListItem(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, l)
}

/*
- handler to delete a product
*/
func (a *App) deleteListItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
        return
	}
	
	l := listItem{ID : id}
	if err := l.deleteListItem(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "Success"})
}

func (a *App) initializeRoutes() {
    a.Router.HandleFunc("/listItems", a.getAllListItems).Methods("GET")
    a.Router.HandleFunc("/listItem", a.createListItem).Methods("POST")
    a.Router.HandleFunc("/listItem/{id:[0-9]+}", a.getListItem).Methods("GET")
    a.Router.HandleFunc("/listItem/{id:[0-9]+}", a.updateListItem).Methods("PUT")
    a.Router.HandleFunc("/listItem/{id:[0-9]+}", a.deleteListItem).Methods("DELETE")
}



// -- Helper Functions --
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}