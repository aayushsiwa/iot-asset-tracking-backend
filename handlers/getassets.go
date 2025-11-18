package handlers

import (
	"crud/domain"
	"crud/model"
	"encoding/json"
	"net/http"
)

func GetAssets(w http.ResponseWriter, r *http.Request) {
	assets, err := domain.GetAllAssets(r.Context())

	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := struct {
		Assets []model.Asset `json:"assets"`
	}{
		Assets: assets,
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
