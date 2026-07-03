package models

type AssetCategory struct {
	Id            int    `json:"id"`
	CategoryName  string `json:"categoryName"`
	CategoryGroup string `json:"categoryGroup"`
}
