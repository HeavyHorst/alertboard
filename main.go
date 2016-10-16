package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func main() {
	var err error
	db, err := newBoltStore()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	router := mux.NewRouter().StrictSlash(true)
	chain := alice.New(loggingHandler, timeoutHandler)

	router.HandleFunc("/api/backup", boltBackupHandler(db))
	router.HandleFunc("/api/alerts", alertListHandler(db))
	router.HandleFunc("/api/alerts/{prefix}", alertListHandler(db))
	router.HandleFunc("/api/alert", alertHandler(db))
	router.HandleFunc("/api/alert/{alertID}", alertHandler(db))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./dashboard/")))

	log.Fatal(http.ListenAndServe(":8080", chain.Then(router)))
}
