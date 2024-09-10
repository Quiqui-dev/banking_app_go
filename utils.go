package main

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, v any) error {

	w.Header().Set("Content-Type", "application/json")

	// headers should be set before we call WriteHeader
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func ReadJson(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
