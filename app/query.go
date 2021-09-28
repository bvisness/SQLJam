package app

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// TODO: Surely this is pretty temporary. I just need to display boring query output.
type queryResult struct {
	Columns []string
	Rows    [][]interface{}
}

func openDB() func() {
	var err error
	db, err = sql.Open("sqlite3", "./sakila.db")
	if err != nil {
		panic(err)
	}

	return func() {
		db.Close()
	}
}

func doQuery(q string) *queryResult {
	rows, err := db.Query(q)
	if err != nil {
		log.Print(err)
		return &queryResult{}
	}
	defer rows.Close()

	var res queryResult

	res.Columns, err = rows.Columns()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		row := make([]interface{}, len(res.Columns))
		rowPointers := make([]interface{}, len(row))
		for i := range row {
			rowPointers[i] = &row[i]
		}

		err = rows.Scan(rowPointers...)
		if err != nil {
			panic(err)
		}
		res.Rows = append(res.Rows, row)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return &res
}
