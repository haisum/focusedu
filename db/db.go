package db

import (
	"bytes"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

const (
	FreshSQL = "sqls/fresh.sql"
)

func ConnectSQLite(dbname string) (*sqlx.DB, error) {
	var err error
	if _, err = os.Stat(dbname); os.IsNotExist(err) {
		err = CreateFresh(dbname)
		if err != nil {
			return db, err
		}
	} else {
		db, err = sqlx.Connect("sqlite3", dbname)
	}
	db.MapperFunc(func(s string) string {
		return s
	})
	return db, err
}

func Get() *sqlx.DB {
	return db
}

func CreateFresh(dbname string) error {
	b, err := ioutil.ReadFile(FreshSQL)
	if err != nil {
		log.WithError(err).Errorf("Error in reading fresh sql file from path %s", FreshSQL)
		return err
	}
	db, err = sqlx.Connect("sqlite3", dbname)
	if err != nil {
		log.WithError(err).Error("Couldn't connect to db")
	}
	var tx *sqlx.Tx
	tx, err = db.Beginx()
	n := bytes.IndexByte(b, 0)
	var query string
	//null char not found
	if n == -1 {
		query = string(b[:])
	} else {
		query = string(b[:n])
	}
	_, err = tx.Exec(query)
	if err != nil {
		log.WithError(err).Error("Error in executing fresh sql file.")
		tx.Rollback()
		db.Close()
		os.Remove(dbname)
		return err
	} else {
		err = tx.Commit()
		if err != nil {
			log.WithError(err).Error("Couldn't commit")
			tx.Rollback()
			db.Close()
			os.Remove(dbname)
			return err
		}
	}
	log.Info("Successfully executed fresh sql.")
	return nil
}
