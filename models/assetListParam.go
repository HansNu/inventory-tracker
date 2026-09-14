package models

type AssetListParams struct {
	Page          string
	PageSize      string
	Search        string
	Status        string
	AssetCategory string
	SortField     string
	SortOrder     string
}
