package domain

import (
	"context"
	"crud/db"
	"crud/helpers"
	"crud/model"
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

var ErrNoLocations = errors.New("no locations found")

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
		return err
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
		return nil, err
	}
	defer rows.Close()

	locations := []model.Location{}

	for rows.Next() {
		var loc model.Location
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.Code, &loc.CreatedAtUTC, &loc.LastUpdatedAtUTC); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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
		if err := helpers.HandlePostgresError(err); err != nil {
			return err
		}
		return err
	}

	return nil
}

func UpdateLocation(ctx context.Context, p model.LocationPatch) (*model.Location, error) {
	var b strings.Builder

	args := []any{}
	argIdx := 1

	b.WriteString(`UPDATE locations SET `)

	if p.Name != nil {
		b.WriteString(`"name" = $%d`)
		args = append(args, *p.Name)
		argIdx++
	}

	if p.Code != nil {
		b.WriteString(`"code" = $%d`)
		args = append(args, *p.Code)
		argIdx++
	}

	b.WriteString(`, "lastUpdatedAtUTC" = NOW() WHERE "ID" = $%d RETURNING "ID"`)

	args = append(args, p.ID, argIdx)

	query := b.String()

	loc := &model.Location{}
	if err := db.DB.QueryRowContext(ctx, query, args...).Scan(&loc.ID); err != nil {
		slog.Error(`{"error":"` + err.Error() + `"}`)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, helpers.ErrLocationDoesNotExist
		}

		if err := helpers.HandlePostgresError(err); err != nil {
			return nil, err
		}

		return nil, err
	}

	return loc, nil
}
