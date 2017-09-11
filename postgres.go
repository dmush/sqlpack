package main

import (
	"database/sql"
	"strconv"

	"github.com/lib/pq"
)

const failTextOffset int = 20

type postgres struct {
	sql *sql.DB
}

func openPostgres(source string) (db *postgres, err error) {
	var sqldb *sql.DB
	sqldb, err = sql.Open("postgres", source)
	if err != nil {
		return
	}
	db = &postgres{sqldb}
	return
}

func (db *postgres) execLog(text string, l logger) {
	err := db.exec(text)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			i, err := strconv.Atoi(pqErr.Position)
			if err != nil {
				l.fail(pqErr.Message)
			} else {
				l.failSQL(pqErr.Message, text, i-1)
			}
		} else {
			l.fail(err.Error())
		}
		return
	}
	l.ok()
}

func (db *postgres) exec(text string) (err error) {
	_, err = db.sql.Exec(text)
	return
}

func (db *postgres) close() error {
	return db.sql.Close()
}
