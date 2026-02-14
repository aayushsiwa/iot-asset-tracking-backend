package handlers

import (
	"crud/domain"
	"crud/model"
	"encoding/json"
	"net/http"
)

func GetLocation(w http.ResponseWriter, r *http.Request) {
	locations, err := domain.GetLocations(r.Context())

	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := struct {
		Locations []model.Location `json:"locations"`
	}{
		Locations: locations,
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "json marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
