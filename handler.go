package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/pressly/chi"
	"github.com/rs/xlog"
)

func errorsFromContext(ctx context.Context) string {
	e, ok := ctx.Value("error").(string)
	if ok {
		return e
	}
	return ""
}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()
		t1 := time.Now()
		next.ServeHTTP(rec, r)

		l := xlog.FromRequest(r)
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		w.Write(rec.Body.Bytes())

		ctx := r.Context()
		t2 := time.Now()
		l.Info(xlog.F{
			"duration": t2.Sub(t1),
			"status":   rec.Code,
			"size":     rec.Body.Len(),
			"error":    errorsFromContext(ctx),
		})
	})
}

func alertCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		alertID, err := base64.URLEncoding.DecodeString(chi.URLParam(r, "alertID"))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "alertID", string(alertID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAlert(db store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		alertID, ok := ctx.Value("alertID").(string)
		l := xlog.FromRequest(r)
		if !ok {
			http.Error(w, http.StatusText(422), 422)
			return
		}
		l.SetField("message", "Successfully returned alert: "+alertID)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(db.getAlert(string(alertID)))
	}
}

func deleteAlert(db store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		alertID, ok := ctx.Value("alertID").(string)
		if !ok {
			http.Error(w, http.StatusText(422), 422)
			return
		}

		data := db.getAlert(string(alertID))
		err := db.deleteAlert(string(alertID))
		l := xlog.FromRequest(r)
		if err != nil {
			l.SetField("error", fmt.Sprintf("Wasn't able to delete %s: %s", alertID, err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			l.SetField("message", "Successfully deleted alert: "+alertID)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(data)
		}
		return
	}
}

func postAlert(db store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var alert alertData
		err := decoder.Decode(&alert)
		l := xlog.FromRequest(r)
		if err != nil {
			l.SetField("error", "Wasn't able to decode the message body: "+err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = db.putAlert(alert)
		if err != nil {
			l.SetField("error", "Wasn't able to save the alert: "+err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		l.SetField("message", "Successfully created alert: "+alert.ID)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Location", "/api/alert/"+base64.URLEncoding.EncodeToString([]byte(alert.ID)))
		w.WriteHeader(201) // Status 201 -- created
		w.Write(db.getAlert(alert.ID))
		return
	}
}

func listAlerts(db store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		prefix := chi.URLParam(r, "prefix")
		data, num, err := db.getAlertsByPrefix(prefix)
		l := xlog.FromRequest(r)
		if err != nil {
			l.SetField("error", "Wasn't able to get alerts: "+err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			l.SetField("message", fmt.Sprintf("Successfully returned %d alerts", num))
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(data)
		}
	}
}

func boltBackupHandler(db store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := xlog.FromRequest(r)
		db.backup(w)
		l.SetField("message", "Database backup started")
	}
}
