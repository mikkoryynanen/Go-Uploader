package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

// TODO Maybe should not be global?
var az = NewAzureService()

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		files, err := az.GetFiles()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		WriteJSON(w, http.StatusOK, files)
	} 
	
	if r.Method == "POST" {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// TODO Giving placeholder name for blob
		// NOTE This blob is not any specific type, only an array of bytes, types are deferred after download
		az.Upload(bytes, "uploaded_blob")
	}
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		bytes, err := az.Download("uploaded_blob")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !utf8.Valid(bytes) {
			WriteJSON(w, http.StatusBadRequest, fmt.Sprint("Unsupported type. Currently only files containing strings are allowed"))
			return
		}
		WriteJSON(w, http.StatusOK, string(bytes))
	} else {
		WriteJSON(w, http.StatusBadRequest, fmt.Sprintf("Unsupported method %v", r.Method))
	}
}

func WriteJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(value)
}

func main() {
	r := mux.NewRouter()
    r.HandleFunc("/", GetHandler)
	r.HandleFunc("/download", DownloadHandler)
    http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":3000", r))
}

