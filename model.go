package main

import (
	"database/sql"
)

//listItem entity
type listItem struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}

//get
func (l *listItem) getListItem(db *sql.DB) error {
	return db.QueryRow("SELECT name, description FROM listItems WHERE id=$1",
		l.ID).Scan(&l.Name, &l.Description)
}

//update
func (l *listItem) updateListItem(db *sql.DB) error {
	_, err := db.Exec("UPDATE listItems SET name=$1, description=$2 WHERE id=$3", l.Name, l.Description, l.ID)
	return err
}

//delete
func (l *listItem) deleteListItem(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM listItems WHERE id=$1", l.ID)
	return err
}

//create
func (l *listItem) createListItem(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO listItems(name, description) VALUES($1, $2) RETURNING ID", l.Name, l.Description).Scan(&l.ID)

	if err != nil {
		return err
	}
	return nil
}

//get all
func getAllListItems(db *sql.DB, start, count int) ([]listItem, error) {
	rows, err := db.Query("SELECT id, name, description FROM listItems LIMIT $1 OFFSET $2", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	listItems := []listItem{}

	for rows.Next() {
		var l listItem
		if err := rows.Scan(&l.ID, &l.Name, &l.Description); err != nil {
			return nil, err
		}
		listItems = append(listItems, l)
	}
	return listItems, nil
}
