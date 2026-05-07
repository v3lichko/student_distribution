package main

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}
	writeJson(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func main() {
	serve := http.NewServeMux()
	serve.HandleFunc("/health", healthHandler)
	err := http.ListenAndServe(":8080", serve)
	if err != nil {
		panic(err)
	}

}
