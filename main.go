package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		az := NewAzureService()
		files, err := az.GetFiles()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		WriteJSON(w, http.StatusOK, files)
	} 
	
	if r.Method == "POST" {
		log.Println("posting")
	}
}

func WriteJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(value)
}

func main() {
	r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":3000", r))
}

