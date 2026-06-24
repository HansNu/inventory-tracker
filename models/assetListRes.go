package models

type AssetListResponse struct {
	Data  []Asset `json:"data"`
	Total int     `json:"total"`
}
