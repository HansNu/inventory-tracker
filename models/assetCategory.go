package models

type AssetCategory struct {
	Id            int    `json:"id"`
	CategoryName  string `json:"category_name"`
	CategoryGroup string `json:"category_group"`
}
