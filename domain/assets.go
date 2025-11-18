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

func GetAllAssets(ctx context.Context) ([]model.Asset, error) {
	query := `
		SELECT a."ID", a."name", a."status", l."name" AS location, a."lastUpdatedAtUTC", a."createdAtUTC"
		FROM assets a
		JOIN locations l ON a."locationID" = l."ID";
	`

	rows, err := db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assets := []model.Asset{}

	for rows.Next() {
		var a model.Asset

		if err := rows.Scan(&a.ID, &a.Name, &a.Status, &a.Location, &a.LastUpdatedAtUTC, &a.CreatedAtUTC); err != nil {
			fmt.Println(err)
			// return nil, ErrGetAllAssetsFailed
			return nil, err
		}

		assets = append(assets, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}

func GetAssetsByLocation(ctx context.Context, locationID uuid.UUID) ([]model.Asset, error) {
	var exists bool
	checkQuery := `SELECT EXISTS (SELECT 1 FROM locations WHERE "ID" = $1);`

	if err := db.DB.QueryRowContext(ctx, checkQuery, locationID).Scan(&exists); err != nil {
		return nil, err
	}

	if !exists {
		return nil, helpers.ErrLocationDoesNotExist
	}

	query := `
		SELECT a."ID", a."name", a."status", l."name" AS location, a."lastUpdatedAtUTC", a."createdAtUTC"
		FROM assets a
		JOIN locations l ON a."locationID" = l."ID"
    WHERE a."locationID" = $1;
	`

	rows, err := db.DB.QueryContext(ctx, query, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assets := []model.Asset{}

	for rows.Next() {
		var a model.Asset

		if err := rows.Scan(&a.ID, &a.Name, &a.Status, &a.Location, &a.LastUpdatedAtUTC, &a.CreatedAtUTC); err != nil {
			return nil, err
		}

		assets = append(assets, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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

		return err
	}

	return nil
}

func UpdateAsset(ctx context.Context, locationID uuid.UUID, assetID uuid.UUID, patch model.AssetPatch) (*model.Asset, error) {
	var setParts []string
	var args []any
	argIdx := 1

	if patch.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, patch.Name)
		argIdx++
	}

	if patch.Status != "" {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, patch.Status)
		argIdx++
	}

	if len(setParts) == 0 {
		return nil, errors.New("no valid fields to update")
	}

	setParts = append(setParts, `"lastUpdatedAtUTC" = NOW()`)

	query := fmt.Sprintf(`
        UPDATE assets
        SET %s
        WHERE "ID" = $%d AND "locationID" = $%d
        RETURNING "ID";
    `,
		strings.Join(setParts, ", "),
		argIdx,
		argIdx+1,
	)

	args = append(args, assetID, locationID)

	asset := &model.Asset{}
	if err := db.DB.QueryRowContext(ctx, query, args...).Scan(&asset.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helpers.ErrAssetDoesNotExist
		}

		if err := helpers.HandlePostgresError(err); err != nil {
			return nil, err
		}

		return nil, err
	}

	return asset, nil
}

func DeleteAsset(ctx context.Context, locationID uuid.UUID, assetID uuid.UUID) error {
	query := `
        DELETE FROM assets 
        WHERE "locationID" = $1 AND "ID" = $2;
    `

	result, err := db.DB.ExecContext(ctx, query, locationID, assetID)
	if err != nil {
		if err := helpers.HandlePostgresError(err); err != nil {
			return err
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return helpers.ErrAssetDoesNotExist
	}

	return nil
}
