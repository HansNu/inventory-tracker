package repository

import (
	"context"
	"fmt"
	"strings"

	"inventory-tracker/apperr"
	"inventory-tracker/models"
)

type AssetRepository interface {
	CreateAsset(ctx context.Context, asset models.AddAssetReq) error
	GetAssetList(ctx context.Context, q models.AssetListQuery) ([]models.AssetResponse, int, error)
	GetAssetByCode(ctx context.Context, assetCode string) (models.AssetResponse, error)
	UpdateAsset(ctx context.Context, id int, asset models.AddAssetReq) error
	DeleteAssetByCode(ctx context.Context, assetCode string) error
}

type AssetRepo struct {
	DB DBTX
}

func NewAssetRepo(db DBTX) AssetRepository {
	return &AssetRepo{DB: db}
}

// sortableColumns is an allowlist, not a convenience. A sort column cannot be a
// bound parameter, so it has to be interpolated into the SQL string — which makes
// any unvalidated value an injection point. Mapping the caller's name to a column
// we control closes that.
var sortableColumns = map[string]string{
	"id":            "a.id",
	"asset_code":    "a.asset_code",
	"asset_name":    "a.asset_name",
	"brand":         "a.brand",
	"serial_number": "a.serial_number",
	"status":        "a.status",
	"location":      "a.location",
	"user":          `a."user"`,
	"purchase_date": "a.purchase_date",
	"category_name": "c.category_name",
}

const assetSelect = `SELECT a.id, a.asset_code, a.asset_name, c.id, c.category_name, a.brand,
	a.serial_number, a.status, a.location, a."user", a.purchase_date, a.description`

func (r *AssetRepo) CreateAsset(ctx context.Context, asset models.AddAssetReq) error {
	_, err := r.DB.Exec(ctx,
		`INSERT INTO asset (asset_code, asset_name, brand, serial_number, category_id, status, location, "user", purchase_date, description)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		asset.AssetCode, asset.AssetName, asset.Brand, asset.SerialNumber,
		asset.CategoryId, asset.Status, asset.Location, asset.User,
		asset.PurchaseDate, asset.Description,
	)
	return translate(err)
}

func (r *AssetRepo) GetAssetList(ctx context.Context, q models.AssetListQuery) ([]models.AssetResponse, int, error) {
	column, ok := sortableColumns[q.SortField]
	if !ok {
		return nil, 0, apperr.ErrInvalidSortField
	}
	order := "DESC"
	if q.SortOrder == "ASC" {
		order = "ASC"
	}

	sql := assetSelect + `, COUNT(*) OVER() AS total_count
		FROM asset a
		JOIN category c ON c.id = a.category_id
		WHERE 1=1`

	args := []any{}
	idx := 1

	if q.Search != "" {
		cols := []string{
			"a.asset_code", "a.asset_name", "c.category_name", "a.brand",
			"a.serial_number", "a.status", "a.location", `a."user"`,
		}
		conds := make([]string, 0, len(cols))
		for _, col := range cols {
			conds = append(conds, fmt.Sprintf("%s ILIKE $%d", col, idx))
			args = append(args, "%"+q.Search+"%")
			idx++
		}
		sql += " AND (" + strings.Join(conds, " OR ") + ")"
	}
	if q.Status != "" {
		sql += fmt.Sprintf(" AND a.status = $%d", idx)
		args = append(args, q.Status)
		idx++
	}
	if q.AssetCategory != "" {
		sql += fmt.Sprintf(" AND c.category_name = $%d", idx)
		args = append(args, q.AssetCategory)
		idx++
	}

	sql += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", column, order, idx, idx+1)
	args = append(args, q.Limit, q.Offset)

	rows, err := r.DB.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, translate(err)
	}
	defer rows.Close()

	// Initialised, not declared nil: a nil slice marshals to JSON null, which
	// breaks any client that calls .map on the result of an empty table.
	assets := []models.AssetResponse{}
	total := 0

	for rows.Next() {
		var a models.AssetResponse
		if err := rows.Scan(
			&a.Id, &a.AssetCode, &a.AssetName, &a.CategoryId, &a.CategoryName,
			&a.Brand, &a.SerialNumber, &a.Status, &a.Location, &a.User,
			&a.PurchaseDate, &a.Description, &total,
		); err != nil {
			return nil, 0, translate(err)
		}
		assets = append(assets, a)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, translate(err)
	}
	return assets, total, nil
}

func (r *AssetRepo) GetAssetByCode(ctx context.Context, assetCode string) (models.AssetResponse, error) {
	var a models.AssetResponse
	err := r.DB.QueryRow(ctx,
		assetSelect+`
		 FROM asset a
		 JOIN category c ON c.id = a.category_id
		 WHERE a.asset_code = $1`, assetCode).
		Scan(&a.Id, &a.AssetCode, &a.AssetName, &a.CategoryId, &a.CategoryName,
			&a.Brand, &a.SerialNumber, &a.Status, &a.Location, &a.User,
			&a.PurchaseDate, &a.Description)
	if err != nil {
		return models.AssetResponse{}, translate(err)
	}
	return a, nil
}

func (r *AssetRepo) UpdateAsset(ctx context.Context, id int, asset models.AddAssetReq) error {
	tag, err := r.DB.Exec(ctx,
		`UPDATE asset SET asset_code=$1, asset_name=$2, brand=$3, serial_number=$4,
		        category_id=$5, status=$6, location=$7, "user"=$8, purchase_date=$9, description=$10
		 WHERE id=$11`,
		asset.AssetCode, asset.AssetName, asset.Brand, asset.SerialNumber,
		asset.CategoryId, asset.Status, asset.Location, asset.User,
		asset.PurchaseDate, asset.Description, id)
	if err != nil {
		return translate(err)
	}
	// The write's own response is the existence check. A SELECT first would be a
	// second round trip and a race.
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *AssetRepo) DeleteAssetByCode(ctx context.Context, assetCode string) error {
	tag, err := r.DB.Exec(ctx, `DELETE FROM asset WHERE asset_code = $1`, assetCode)
	if err != nil {
		return translate(err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}
