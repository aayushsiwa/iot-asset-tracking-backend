package handlers

import (
	"crud/domain"
	"crud/helpers"
	"crud/model"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func GetAssetsByLocation(w http.ResponseWriter, r *http.Request) {
	locationID := r.PathValue("locationID")

	uid, err := uuid.Parse(locationID)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	assets, err := domain.GetAssetsByLocation(r.Context(), uid)

	if err != nil {
		if errors.Is(err, helpers.ErrLocationDoesNotExist) {
			http.Error(w, `{"error":"`+helpers.ErrLocationDoesNotExist.Error()+`"}`, http.StatusNotFound)
			return
		}

		http.Error(w, `{"error":"db error: `+err.Error()+`"}`, http.StatusInternalServerError)
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
