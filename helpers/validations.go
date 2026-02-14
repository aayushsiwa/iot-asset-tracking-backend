package helpers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

var validate = validator.New()

const (
	postgresUniqueConstraintViolationCode  = "23505"
	postgresForeignConstraintViolationCode = "23503"
)

var (
	ErrLocationDoesNotExist  = errors.New("location does not exist")
	ErrLocationAlreadyExists = errors.New("location already exists")
	ErrCodeAlreadyExists     = errors.New("code already exists")
	ErrAssetAlreadyExists    = errors.New("asset already exists")
	ErrAssetDoesNotExist     = errors.New("asset does not exist")
	ErrNoValidFieldsToUpdate = errors.New("no valid fields to update")
)

var (
	pqErrorMap = map[string]error{
		"locations_name_key":     ErrLocationAlreadyExists,
		"locations_code_key":     ErrCodeAlreadyExists,
		"assets_name_key":        ErrAssetAlreadyExists,
		"assets_locationID_fkey": ErrLocationDoesNotExist,
	}
)

func ValidateRequest[T any](w http.ResponseWriter, r *http.Request, req *T) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return err
	}

	if err := validate.Struct(req); err != nil {
		slog.Error(err.Error())
		if ve, ok := err.(validator.ValidationErrors); ok {
			writeValidationErrors(w, ve)
			return err
		}

		writeError(w, "validation_error", "Invalid request")
		return err
	}

	return nil
}

func writeValidationErrors(w http.ResponseWriter, errs validator.ValidationErrors) {
	out := make(map[string]string)

	for _, fe := range errs {
		field := strings.ToLower(fe.Field())

		switch fe.Tag() {
		case "required":
			out[field] = field + " is required"

		case "len":
			out[field] = field + " must be " + fe.Param() + " letters long"

		case "min":
			out[field] = field + " should have minimum " + fe.Param() + " letters"

		case "max":
			out[field] = field + " should have maximum " + fe.Param() + " letters"

		case "uppercase":
			out[field] = field + " must be uppercase"

		case "oneof":
			out[field] = field + " must be one of: " + fe.Param()

		default:
			out[field] = field + " is invalid"
		}
	}

	writeJSON(w, http.StatusBadRequest, map[string]interface{}{
		"error": out,
	})
}

func writeError(w http.ResponseWriter, code string, message string) {
	writeJSON(w, http.StatusBadRequest, map[string]interface{}{
		"error": map[string]string{
			code: message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func HandlePostgresError(err error) error {
	var pqErr *pq.Error
	// print the error for debugging
	if errors.As(err, &pqErr) {
		// print the error for debugging

		if errValue, ok := pqErrorMap[pqErr.Constraint]; ok && (pqErr.Code == postgresUniqueConstraintViolationCode || pqErr.Code == postgresForeignConstraintViolationCode) {
			return errValue
		}
	}
	return nil
}
