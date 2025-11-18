package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
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

	asset, err := domain.UpdateAsset(r.Context(), locationUUID, assetUUID, patch)
	if err != nil {
		slog.Error(`{"error":"` + err.Error() + `"}`)

		if errors.Is(err, helpers.ErrLocationDoesNotExist) {
			http.Error(w, `{"error":"`+helpers.ErrLocationDoesNotExist.Error()+`"}`, http.StatusBadRequest)
			return
		}

		if errors.Is(err, helpers.ErrAssetDoesNotExist) {
			http.Error(w, `{"error":"`+helpers.ErrAssetDoesNotExist.Error()+`"}`, http.StatusNotFound)
			return
		}

		if errors.Is(err, helpers.ErrAssetAlreadyExists) {
			http.Error(w, `{"error":"`+helpers.ErrAssetAlreadyExists.Error()+`"}`, http.StatusConflict)
			return
		}

		http.Error(w, `{"error":"failed to update asset"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(struct {
		ID *uuid.UUID `json:"ID"`
	}{
		ID: asset.ID,
	})
}
