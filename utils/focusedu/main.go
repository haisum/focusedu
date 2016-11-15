package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/haisum/focusedu/db"
	"github.com/haisum/focusedu/states"
)

const (
	SQLiteDBName = "focusedu.db"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	db, err := db.ConnectSQLite(SQLiteDBName)
	if err != nil {
		log.WithError(err).Fatalf("Couldn't connect to database.")
	}
	defer db.Close()
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", states.MasterHandler)
	r.HandleFunc("/logout", states.LogoutHandler)
	port := os.Getenv("FOCUSEDU_PORT")
	if port == "" {
		port = "4444"
	} else {
		log.Info("Found port %s from env var FOCUSEDU_PORT", port)
	}
	log.Infof("Listening for requests on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handlers.CombinedLoggingHandler(os.Stdout, r)))
}
