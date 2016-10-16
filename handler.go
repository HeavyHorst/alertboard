package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// some middleware handlers
func timeoutHandler(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, 10*time.Second, "timed out")
}

func loggingHandler(h http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, h)
}

// the real handlers
func alertHandler(db store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		alertID, err := base64.URLEncoding.DecodeString(vars["alertID"])
		if err != nil {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(db.getAlert(string(alertID)))
			return

		case "POST":
			decoder := json.NewDecoder(r.Body)
			var alert alertData
			err := decoder.Decode(&alert)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = db.putAlert(alert)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Location", "/api/alert/"+base64.URLEncoding.EncodeToString([]byte(alert.ID)))
			w.WriteHeader(201) // Status 201 -- created
			w.Write(db.getAlert(alert.ID))
			return

		case "DELETE":
			data := db.getAlert(string(alertID))
			err := db.deleteAlert(string(alertID))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.Write(data)
			}
			return

		default:
			http.Error(w, "Method not allowed", 405)
			return
		}
	}
}

func alertListHandler(db store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		prefix := vars["prefix"]
		data, err := db.getAlertsByPrefix(prefix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(data)
		}
	}
}

func boltBackupHandler(db store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db.backup(w)
	}
}
