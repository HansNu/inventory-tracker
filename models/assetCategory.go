package models

type AssetCategory struct {
	Id            int    `json:"id"`
	CategoryName  string `json:"categoryName" binding:"required"`
	CategoryGroup string `json:"categoryGroup" binding:"required"`
}
