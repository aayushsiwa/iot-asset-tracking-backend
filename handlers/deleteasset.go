package handlers

import (
	"errors"
	"net/http"

	"crud/domain"
	"crud/helpers"

	"github.com/google/uuid"
)

func DeleteAsset(w http.ResponseWriter, r *http.Request) {
	locationID := r.PathValue("locationID")
	assetID := r.PathValue("assetID")

	locationUUID, err := uuid.Parse(locationID)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	assetUUID, err := uuid.Parse(assetID)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := domain.DeleteAsset(r.Context(), locationUUID, assetUUID); err != nil {
		if errors.Is(err, helpers.ErrAssetDoesNotExist) {
			http.Error(w, `{"error":"`+helpers.ErrAssetDoesNotExist.Error()+`"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"failed to delete asset"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
