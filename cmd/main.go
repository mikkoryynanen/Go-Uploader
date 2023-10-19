package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mikkoryynanen/uploader/api/handlers"
)

func main() {
	r := mux.NewRouter()
    r.HandleFunc("/", handlers.GetHandler)
	r.HandleFunc("/download", handlers.DownloadHandler)
    http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":3000", r))
}

