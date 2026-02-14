package domain

import (
	"context"
	"crud/db"
	"crud/helpers"
	"crud/model"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrCreateLocationFailed = errors.New("failed to create location")
	ErrGetLocationsFailed   = errors.New("failed to get locations")
	ErrDeleteLocationFailed = errors.New("failed to delete location")
	ErrUpdateLocationFailed = errors.New("failed to update location")
)

func CreateLocation(ctx context.Context, location *model.Location) error {
	query := `
		INSERT INTO locations ("name", "code")
		VALUES ($1, $2)
		RETURNING "ID";
	`

	if err := db.DB.QueryRowContext(ctx, query, location.Name, location.Code).Scan(&location.ID); err != nil {
		slog.Error(`{"error":"` + err.Error() + `"}`)
		if err := helpers.HandlePostgresError(err); err != nil {
			return err
		}
		return ErrCreateLocationFailed
	}

	return nil
}

func GetLocations(ctx context.Context) ([]model.Location, error) {
	query := `
		SELECT "ID", "name", "code", "createdAtUTC", "lastUpdatedAtUTC"
		FROM locations;
	`

	rows, err := db.DB.QueryContext(ctx, query)
	if err != nil {
		slog.Error(`{"error":"` + err.Error() + `"}`)

		return nil, ErrGetLocationsFailed
	}
	defer rows.Close()

	locations := []model.Location{}

	for rows.Next() {
		var loc model.Location

		if err := rows.Scan(&loc.ID, &loc.Name, &loc.Code, &loc.CreatedAtUTC, &loc.LastUpdatedAtUTC); err != nil {
			slog.Error(`{"error":"` + err.Error() + `"}`)

			return nil, ErrGetLocationsFailed
		}

		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		slog.Error(`{"error":"` + err.Error() + `"}`)

		return nil, ErrGetLocationsFailed
	}

	return locations, nil
}

func DeleteLocation(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM locations
		WHERE "ID" = $1;
	`

	_, err := db.DB.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error(`{"error":"` + err.Error() + `"}`)

		return ErrDeleteLocationFailed
	}

	return nil
}

func UpdateLocation(ctx context.Context, p model.LocationPatch) (*model.Location, error) {
	var b strings.Builder
	args := []any{}
	argIdx := 1

	b.WriteString(`UPDATE locations SET "lastUpdatedAtUTC" = NOW(), `)

	first := true

	if p.Name != nil {
		if !first {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, `"name" = $%d`, argIdx)
		args = append(args, *p.Name)
		argIdx++
		first = false
	}

	if p.Code != nil {
		if !first {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, `"code" = $%d`, argIdx)
		args = append(args, *p.Code)
		argIdx++
		first = false
	}

	fmt.Fprintf(&b, ` WHERE "ID" = $%d`, argIdx)
	args = append(args, p.ID)

	b.WriteString(` RETURNING "ID"`)

	query := b.String()

	loc := &model.Location{}

	if err := db.DB.QueryRowContext(ctx, query, args...).Scan(&loc.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helpers.ErrLocationDoesNotExist
		}

		slog.Error(`{"error":"` + err.Error() + `"}`)

		if err := helpers.HandlePostgresError(err); err != nil {
			return nil, err
		}

		return nil, ErrUpdateLocationFailed
	}

	return loc, nil
}
