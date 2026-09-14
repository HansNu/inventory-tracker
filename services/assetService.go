package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"inventory-tracker/models"
	"inventory-tracker/repository"
)

type AssetService interface {
	AddAsset(ctx context.Context, asset models.AddAssetReq) error
	GetAssetList(ctx context.Context, params models.AssetListParams) ([]models.AssetResponse, int, error)
}
type assetService struct {
	Repo repository.AssetRepository
}

func NewAssetService(repo repository.AssetRepository) AssetService {
	return &assetService{Repo: repo}
}

func (s *assetService) AddAsset(ctx context.Context, asset models.AddAssetReq) error {
	// validation/business rules land here later — for now, straight passthrough
	return s.Repo.CreateAsset(ctx, asset)
}

func (s *assetService) GetAssetList(ctx context.Context, params models.AssetListParams) ([]models.AssetResponse, int, error) {
	pageInt, _ := strconv.Atoi(params.Page)
	pageSizeInt, _ := strconv.Atoi(params.PageSize)
	if pageInt < 1 {
		pageInt = 1
	}
	if pageSizeInt < 1 {
		pageSizeInt = 10
	}
	offset := (pageInt - 1) * pageSizeInt

	sortOrder := "DESC"
	if params.SortOrder == "ascend" {
		sortOrder = "ASC"
	}

	query := `SELECT a.id, a.asset_code, a.asset_name, c.id, c.category_name, a.brand, a.serial_number,
		a.status, a.location, a."user", a.purchase_date, a.description, COUNT(*) OVER() as total_count
		FROM asset a
		JOIN category c on c.id = a.category_id
		WHERE 1=1`

	args := []any{}
	argIdx := 1
	searchColumns := []string{"asset_code", "asset_name", "category_name", "brand", "serial_number", "status", "location", `"user"`}

	if params.Search != "" {
		conditions := []string{}
		for _, col := range searchColumns {
			conditions = append(conditions, fmt.Sprintf("%s ILIKE $%d", col, argIdx))
			args = append(args, "%"+params.Search+"%")
			argIdx++
		}
		query += " AND (" + strings.Join(conditions, " OR ") + ")"
	}

	if params.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, params.Status)
		argIdx++
	}

	if params.AssetCategory != "" {
		query += fmt.Sprintf(" AND c.category_name = $%d", argIdx)
		args = append(args, params.AssetCategory)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", params.SortField, sortOrder, argIdx, argIdx+1)
	args = append(args, pageSizeInt, offset)

	return s.Repo.GetAssetList(ctx, query, args)
}
