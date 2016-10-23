package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/rs/xlog"
)

func main() {
	var err error
	db, err := newBoltStore()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	host, _ := os.Hostname()
	r := chi.NewRouter()

	r.Use(xlog.NewHandler(xlog.Config{
		Output: xlog.NewOutputChannel(xlog.NewConsoleOutput()),
		Fields: xlog.F{"hostname": host},
	}))
	r.Use(xlog.MethodHandler("method"))
	r.Use(xlog.URLHandler("url"))
	r.Use(xlog.RemoteAddrHandler("ip"))
	r.Use(xlog.RequestIDHandler("req_id", "Request-Id"))
	r.Use(logHandler)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Route("/alert", func(r chi.Router) {
			r.Post("/", postAlert(db))
			r.Route("/:alertID", func(r chi.Router) {
				r.Use(alertCtx)
				r.Get("/", getAlert(db))
				r.Delete("/", deleteAlert(db))
			})
		})
		r.Route("/alerts", func(r chi.Router) {
			r.Get("/", listAlerts(db))
			r.Get("/:prefix", listAlerts(db))
		})
		r.Route("/backup", func(r chi.Router) {
			r.Get("/", boltBackupHandler(db))
		})
	})

	r.FileServer("/", http.Dir("./dashboard/"))

	log.Fatal(http.ListenAndServe(":8080", r))
}
