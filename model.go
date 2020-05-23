package main

import (
	"database/sql"
	"errors"
)

//listItem entity
type listItem struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}

//get
func (l *listItem) getListItem(db *sql.DB) error {
	return errors.New("Not Implemented")
}

//update
func (l *listItem) updateListItem(db *sql.DB) error {
	return errors.New("Not Implemented")
}

//delete
func (l *listItem) deleteListItem(db *sql.DB) error {
	return errors.New("Not Implemented")
}

//create
func (l *listItem) createListItem(db *sql.DB) error {
	return errors.New("Not Implemented")
}

//get all
func getAllListItems(db *sql.DB, start, count int) ([]listItem, error) {
	return nil, errors.New("Not implemented")
}
