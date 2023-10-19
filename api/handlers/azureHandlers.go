package handlers

import (
	"fmt"
	"io"
	"net/http"
	"unicode/utf8"

	"github.com/mikkoryynanen/uploader/internal/azure"
	"github.com/mikkoryynanen/uploader/internal/utils"
)

var az = azure.NewAzureService()

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		files, err := az.GetFiles()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, files)
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
			utils.WriteJSON(w, http.StatusBadRequest, fmt.Sprint("Unsupported type. Currently only files containing strings are allowed"))
			return
		}
		utils.WriteJSON(w, http.StatusOK, string(bytes))
	} else {
		utils.WriteJSON(w, http.StatusBadRequest, fmt.Sprintf("Unsupported method %v", r.Method))
	}
}
