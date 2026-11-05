package service

import (
	"context"
	"strconv"

	"inventory-tracker/models"
	"inventory-tracker/repository"
)

type AssetService interface {
	AddAsset(ctx context.Context, asset models.AddAssetReq) error
	GetAssetList(ctx context.Context, params models.AssetListParams) ([]models.AssetResponse, int, error)
	GetAssetByCode(ctx context.Context, assetCode string) (models.AssetResponse, error)
	UpdateAsset(ctx context.Context, id int, asset models.AddAssetReq) error
	DeleteAssetByCode(ctx context.Context, assetCode string) error
}

type assetService struct {
	Repo repository.AssetRepository
}

func NewAssetService(repo repository.AssetRepository) AssetService {
	return &assetService{Repo: repo}
}

const (
	defaultPageSize = 10
	maxPageSize     = 100
)

func (s *assetService) AddAsset(ctx context.Context, asset models.AddAssetReq) error {
	return s.Repo.CreateAsset(ctx, asset)
}

// GetAssetList turns the raw query-string values into a query the repository can
// execute. This is the service's actual job here: deciding what a missing page
// number, a nonsense page size or an unrecognised sort order mean. The SQL that
// results from those decisions belongs to the repository.
func (s *assetService) GetAssetList(ctx context.Context, params models.AssetListParams) ([]models.AssetResponse, int, error) {
	page, err := strconv.Atoi(params.Page)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(params.PageSize)
	if err != nil || pageSize < 1 {
		pageSize = defaultPageSize
	}
	// Capped so a caller cannot ask for the whole table in one request.
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	sortOrder := "DESC"
	if params.SortOrder == "ascend" || params.SortOrder == "ASC" || params.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	sortField := params.SortField
	if sortField == "" {
		sortField = "purchase_date"
	}

	return s.Repo.GetAssetList(ctx, models.AssetListQuery{
		Search:        params.Search,
		Status:        params.Status,
		AssetCategory: params.AssetCategory,
		SortField:     sortField,
		SortOrder:     sortOrder,
		Limit:         pageSize,
		Offset:        (page - 1) * pageSize,
	})
}

func (s *assetService) GetAssetByCode(ctx context.Context, assetCode string) (models.AssetResponse, error) {
	return s.Repo.GetAssetByCode(ctx, assetCode)
}

func (s *assetService) UpdateAsset(ctx context.Context, id int, asset models.AddAssetReq) error {
	return s.Repo.UpdateAsset(ctx, id, asset)
}

func (s *assetService) DeleteAssetByCode(ctx context.Context, assetCode string) error {
	return s.Repo.DeleteAssetByCode(ctx, assetCode)
}
