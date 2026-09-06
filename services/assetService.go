package service

import (
	"context"

	"inventory-tracker/models"
	"inventory-tracker/repository"
)

type AssetService interface {
	AddAsset(ctx context.Context, asset models.AddAssetReq) error
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
