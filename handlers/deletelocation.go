package handlers

import (
	"net/http"

	"crud/domain"

	"github.com/google/uuid"
)

func DeleteLocation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := domain.DeleteLocation(r.Context(), uid); err != nil {
		http.Error(w, `{"error":"failed to delete location"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
