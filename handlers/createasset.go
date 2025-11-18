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

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	locationID := r.PathValue("locationID")

	req := struct {
		Name   string `json:"name" validate:"required,min=5,max=50"`
		Status string `json:"status" validate:"required,oneof=online offline"`
	}{}

	if err := helpers.ValidateRequest(w, r, &req); err != nil {
		return
	}

	asset := &model.CreateAssetRequest{
		Name:       req.Name,
		Status:     model.Status(req.Status),
		LocationID: locationID,
	}

	if err := domain.CreateAsset(r.Context(), asset); err != nil {
		if errors.Is(err, helpers.ErrLocationDoesNotExist) {
			http.Error(w, `{"error":"location does not exist"}`, http.StatusBadRequest)
			return
		}

		if errors.Is(err, helpers.ErrAssetAlreadyExists) {
			http.Error(w, `{"error":"`+helpers.ErrAssetAlreadyExists.Error()+`"}`, http.StatusConflict)
			return
		}

		http.Error(w, `{"error":"failed to create asset"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID *uuid.UUID `json:"ID"`
	}{
		ID: asset.ID,
	})
}
