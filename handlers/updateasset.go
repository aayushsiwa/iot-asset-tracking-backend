package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"crud/domain"
	"crud/helpers"
	"crud/model"

	"github.com/google/uuid"
)

func UpdateAsset(w http.ResponseWriter, r *http.Request) {
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
	patch := model.AssetPatch{}
	if err := helpers.ValidateRequest(w, r, &patch); err != nil {
		return
	}

	if patch.Name == nil && patch.Status == nil {
		http.Error(w, `{"error":"`+helpers.ErrNoValidFieldsToUpdate.Error()+`"}`, http.StatusBadRequest)
		return
	}

	asset, err := domain.UpdateAsset(r.Context(), locationUUID, assetUUID, patch)
	if err != nil {
		if errors.Is(err, helpers.ErrAssetDoesNotExist) {
			w.WriteHeader(http.StatusOK)
			return
		}

		if errors.Is(err, helpers.ErrAssetAlreadyExists) {
			http.Error(w, `{"error":"`+helpers.ErrAssetAlreadyExists.Error()+`"}`, http.StatusConflict)
			return
		}

		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		ID *uuid.UUID `json:"ID"`
	}{
		ID: asset.ID,
	})
}
