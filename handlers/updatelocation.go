package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"crud/domain"
	"crud/helpers"
	"crud/model"

	"github.com/google/uuid"
)

func UpdateLocation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	patchReq := model.UpdateLocationRequest{}
	if err := helpers.ValidateRequest(w, r, &patchReq); err != nil {
		return
	}

	if patchReq.Code != nil {
		code := strings.ToUpper(*patchReq.Code)
		patchReq.Code = &code
	}

	patch := model.LocationPatch{
		ID:   uid,
		Name: patchReq.Name,
		Code: patchReq.Code,
	}

	loc, err := domain.UpdateLocation(r.Context(), patch)
	if err != nil {
		if errors.Is(err, helpers.ErrLocationAlreadyExists) {
			http.Error(w, `{"error":"`+helpers.ErrLocationAlreadyExists.Error()+`"}`, http.StatusConflict)
			return
		}

		if errors.Is(err, helpers.ErrCodeAlreadyExists) {
			http.Error(w, `{"error":"`+helpers.ErrCodeAlreadyExists.Error()+`"}`, http.StatusConflict)
			return
		}

		if errors.Is(err, helpers.ErrLocationDoesNotExist) {
			http.Error(w, `{"error":"`+helpers.ErrLocationDoesNotExist.Error()+`"}`, http.StatusNotFound)
			return
		}

		http.Error(w, `{"error":"update location failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		ID *uuid.UUID `json:"ID"`
	}{
		ID: loc.ID,
	})
}
