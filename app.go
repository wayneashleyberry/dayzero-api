package app

import (
	"encoding/json"
	"fmt"
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
	body, err := coct.Get(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dash, err := coct.Parse(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(dash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
