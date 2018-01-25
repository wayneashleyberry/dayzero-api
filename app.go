package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/wayneashleyberry/dayzero-app/pkg/coct"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/api/dashboard", dashboardHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
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
}
