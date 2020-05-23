package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
	"log"
	"net/http"
	"strconv"
)

//global var for testing
var a App

//ensure DB properly setup, clears DB at end
func TestMain(m *testing.M) {
	a.InitialiseDB(
        os.Getenv("APP_DB_USERNAME"),
        os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))
		
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS listItems
(
	id SERIAL,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	CONSTRAINT listItems_pkey PRIMARY KEY (id)
)`

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM listItems")
	a.DB.Exec("ALTER SEQUENCE listItems_id_seq RESTART WITH 1")
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/listItems", nil)
	response := executeRequest(req);

	checkResponseCode(t, http.StatusOK, response.Code)

	//check textual body of response is empty array
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected empty array. Got %s", body)
	}
}

/*
- try to access non-existent listitem in DB
- test HTTP response is 404
- test reponse error message = "Product not found"
*/
func TestGetNonExistentListItem(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/listItem/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	//transfer json into hashmap
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the rror key of response to be set to Product not found. Got %s", m["error"])
	}
}

/*
- manually add a listItem into DB
- test HTTP response is 201
- test JSON response contains correct information
*/
func TestCreateListItem(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"name": "test item", "description": "this is a test"}`)
	req, _ := http.NewRequest("POST", "/listItem", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	//check response json
	if m["name"] != "test item" {
		t.Errorf("Expected name to be test item. Got %v", m["name"])
	}
	if m["description"] != "this is a test" {
		t.Errorf("Expected description to be this is a test. Got %v", m["description"])
	}
	//id is comapred to 1.0 as JSON unmarshalling converts numbers to float when target is map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected ID to be 1. Got %v", m["id"])
	}
}

/*
- test fetching a legit listitem from DB
- test HTTP response code is 200
*/
func TestGetListItem(t *testing.T) {
	clearTable()
	addListItems(1)

	req, _ := http.NewRequest("GET", "/listItem/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

/*
- update an existing listItem in DB
- test HTTP response code is 200
- test response body contains correct updated details
*/
func TestUpdateListItem(t *testing.T) {
	clearTable()
	addListItems(1)

	//get item first
	req, _ := http.NewRequest("GET", "/listItem/1", nil)
	response := executeRequest(req)
	var oriListItem map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oriListItem)

	//update item
	var jsonStr = []byte(`{"name":"test listItem - updated name", "description": "test description1"}`)
	req, _ = http.NewRequest("PUT", "/listItem/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-type", "application/json")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	//check values
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["id"] != oriListItem["id"] {
		t.Errorf("Expected id to remain the same %v. Go %v", oriListItem["id"], m["id"])
	}
	if m["name"] != "test listItem - updated name" {
		t.Errorf("Expected name to change to test listItem - updated name. Got %v", m["name"])
	}
	if m["description"] != "test description1" {
		t.Errorf("Expected description to change to test description1. Got %v", m["description"])
	}
}

/*
- delete listItem from DB
- test status is 404 not found
*/
func TestDeleteListItem(t *testing.T) {
	clearTable()
	addListItems(1)

	req, _ := http.NewRequest("GET", "/listItem/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/listItem/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/listItem/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}


//-- helper functions --
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addListItems(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO listItems(name, description) VALUES($1, $2)", "ListItem "+strconv.Itoa(i), "test description")
	}
}
 

