package handlers

import "net/http"

func ClinicsHandler(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("[]"))
}

func EmergencyClinicsHandler(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("[]"))
}
