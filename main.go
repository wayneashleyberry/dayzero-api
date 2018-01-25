package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/wayneashleyberry/dayzero-app/pkg/coct"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kelseyhightower/envconfig"
)

type specification struct {
	Port int64 `default:"8080"`
}

func main() {
	var s specification
	envconfig.MustProcess("", &s)

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Get("/api/dashboard", func(w http.ResponseWriter, r *http.Request) {
		body, err := coct.Get()
		if err != nil {
			log.Fatal(err)
		}
		dash, err := coct.Parse(body)
		if err != nil {
			log.Fatal(err)
		}
		b, err := json.Marshal(dash)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	port := strconv.FormatInt(s.Port, 10)
	http.ListenAndServe(":"+port, r)
}
