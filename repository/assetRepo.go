package repository

import (
	"context"

	"inventory-tracker/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AssetRepository interface {
	CreateAsset(ctx context.Context, asset models.AddAssetReq) error
	GetAssetList(ctx context.Context, query string, args []any) ([]models.AssetResponse, int, error)
}

// testing
type DBTX interface {
	Close()
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}
type AssetRepo struct {
	DB DBTX
}

func NewAssetRepo(db DBTX) AssetRepository {
	return &AssetRepo{DB: db}
}

// prod
func (r *AssetRepo) CreateAsset(ctx context.Context, asset models.AddAssetReq) error {
	_, err := r.DB.Exec(ctx,
		`INSERT INTO asset (asset_code, asset_name, brand, serial_number, category_id, status, location, "user", purchase_date, description)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		asset.AssetCode, asset.AssetName, asset.Brand, asset.SerialNumber,
		asset.CategoryId, asset.Status, asset.Location, asset.User,
		asset.PurchaseDate, asset.Description,
	)
	return err
}

func (r *AssetRepo) GetAssetList(ctx context.Context, query string, args []any) ([]models.AssetResponse, int, error) {
	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var assets []models.AssetResponse
	var total int

	for rows.Next() {
		var a models.AssetResponse
		if err := rows.Scan(
			&a.Id, &a.AssetCode, &a.AssetName, &a.CategoryId, &a.CategoryName,
			&a.Brand, &a.SerialNumber, &a.Status, &a.Location, &a.User,
			&a.PurchaseDate, &a.Description, &total,
		); err != nil {
			return nil, 0, err
		}
		assets = append(assets, a)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}
