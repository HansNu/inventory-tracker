package models

type AssetListResponse struct {
	Data  []AssetResponse `json:"data"`
	Total int             `json:"total"`
}
