package domain

import (
	"context"
	"crud/db"
	"crud/helpers"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"crud/model"

	"github.com/google/uuid"
)

var (
	ErrGetAllAssetsFailed       = errors.New("failed to get all assets")
	ErrGetAssetByLocationFailed = errors.New("failed to get assets")
	ErrCreateAssetFailed        = errors.New("failed to create asset")
	ErrUpdateAssetFailed        = errors.New("failed to update asset")
	ErrDeleteAssetFailed        = errors.New("failed to delete asset")
)

func GetAllAssets(ctx context.Context) ([]model.Asset, error) {
	query := `
		SELECT a."ID", a."name", a."status", l."name" AS location, a."lastUpdatedAtUTC", a."createdAtUTC"
		FROM assets a
		JOIN locations l ON a."locationID" = l."ID";
	`

	rows, err := db.DB.QueryContext(ctx, query)
	if err != nil {
		slog.Error(`{"error":"` + err.Error() + `"}`)

		return nil, ErrGetAllAssetsFailed
	}
	defer rows.Close()

	assets := []model.Asset{}

	for rows.Next() {
		var a model.Asset

		if err := rows.Scan(&a.ID, &a.Name, &a.Status, &a.Location, &a.LastUpdatedAtUTC, &a.CreatedAtUTC); err != nil {
			slog.Error(`{"error":"` + err.Error() + `"}`)

			return nil, ErrGetAllAssetsFailed
		}

		assets = append(assets, a)
	}

	if err := rows.Err(); err != nil {
		slog.Error(`{"error":"` + err.Error() + `"}`)

		return nil, ErrGetAllAssetsFailed
	}

	return assets, nil
}

func GetAssetsByLocation(ctx context.Context, locationID uuid.UUID) ([]model.Asset, error) {
	query := `
	SELECT a."ID", a."name", a."status", l."name" AS location, a."lastUpdatedAtUTC", a."createdAtUTC"
	FROM assets a
	JOIN locations l ON a."locationID" = l."ID"
    WHERE a."locationID" = $1;
	`

	rows, err := db.DB.QueryContext(ctx, query, locationID)
	if err != nil {
		slog.Error(`{"error":"` + err.Error() + `"}`)

		return nil, ErrGetAssetByLocationFailed
	}
	defer rows.Close()

	assets := []model.Asset{}

	for rows.Next() {
		var a model.Asset

		if err := rows.Scan(&a.ID, &a.Name, &a.Status, &a.Location, &a.LastUpdatedAtUTC, &a.CreatedAtUTC); err != nil {
			slog.Error(`{"error":"` + err.Error() + `"}`)

			return nil, ErrGetAssetByLocationFailed
		}

		assets = append(assets, a)
	}

	if err := rows.Err(); err != nil {
		slog.Error(`{"error":"` + err.Error() + `"}`)

		return nil, ErrGetAssetByLocationFailed
	}

	return assets, nil
}

func CreateAsset(ctx context.Context, a *model.CreateAssetRequest) error {
	query := `
		INSERT INTO assets ("name", "status", "locationID")
		VALUES ($1,$2,$3)
		RETURNING "ID";
	`

	if err := db.DB.QueryRowContext(ctx, query,
		a.Name,
		a.Status,
		a.LocationID,
	).Scan(&a.ID); err != nil {
		slog.Error(`{"error":"` + err.Error() + `"}`)

		if err := helpers.HandlePostgresError(err); err != nil {
			return err
		}

		return ErrCreateAssetFailed
	}

	return nil
}

func UpdateAsset(ctx context.Context, locationID uuid.UUID, assetID uuid.UUID, patch model.AssetPatch) (*model.Asset, error) {
	var b strings.Builder
	var args []any
	argIdx := 1

	b.WriteString(`UPDATE assets SET "lastUpdatedAtUTC" = NOW(), `)

	first := true

	if patch.Name != nil {
		if !first {
			b.WriteString(", ")
		}

		fmt.Fprintf(&b, `name = $%d`, argIdx)
		args = append(args, patch.Name)
		argIdx++
		first = false
	}

	if patch.Status != nil {
		if !first {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, `status = $%d`, argIdx)
		args = append(args, patch.Status)
		argIdx++
		first = false
	}

	// TODO: move this to handler
	if first {
		return nil, errors.New("no valid fields to update")
	}

	fmt.Fprintf(&b, ` WHERE "ID" = $%d AND "locationID" = $%d RETURNING "ID"`, argIdx, argIdx+1)
	args = append(args, assetID, locationID)

	query := b.String()

	asset := &model.Asset{}
	if err := db.DB.QueryRowContext(ctx, query, args...).Scan(&asset.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helpers.ErrAssetDoesNotExist
		}

		slog.Error(`{"error":"` + err.Error() + `"}`)

		if err := helpers.HandlePostgresError(err); err != nil {
			return nil, err
		}

		return nil, ErrUpdateAssetFailed
	}

	return asset, nil
}

func DeleteAsset(ctx context.Context, locationID uuid.UUID, assetID uuid.UUID) error {
	query := `
        DELETE FROM assets 
        WHERE "locationID" = $1 AND "ID" = $2;
    `

	_, err := db.DB.ExecContext(ctx, query, locationID, assetID)
	if err != nil {
		slog.Error(`{"error":"` + err.Error() + `"}`)

		return ErrDeleteAssetFailed
	}

	return nil
}
