package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5xx error:%s\n", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Marshaling Json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(dat)
}
