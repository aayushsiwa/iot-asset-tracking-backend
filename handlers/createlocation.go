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

func CreateLocation(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Name string `json:"name" validate:"required,min=5,max=50"`
		Code string `json:"code" validate:"required,len=4"`
	}{}

	if err := helpers.ValidateRequest(w, r, &req); err != nil {
		return
	}

	req.Code = strings.ToUpper(req.Code)

	location := &model.Location{
		Name: req.Name,
		Code: req.Code,
	}

	err := domain.CreateLocation(r.Context(), location)
	if err != nil {
		if errors.Is(err, helpers.ErrLocationAlreadyExists) {
			http.Error(w, `{"error":"`+helpers.ErrLocationAlreadyExists.Error()+`"}`, http.StatusConflict)
			return
		}
		if errors.Is(err, helpers.ErrCodeAlreadyExists) {
			http.Error(w, `{"error":"`+helpers.ErrCodeAlreadyExists.Error()+`"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error":"create location failed"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID *uuid.UUID `json:"ID"`
	}{
		ID: location.ID,
	})
}
