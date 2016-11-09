package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func ConnectSQLite(dbname string) error {
	var err error
	db, err = sqlx.Connect("sqlite3", dbname)
	return err
}

func Get() *sqlx.DB {
	return db
}
