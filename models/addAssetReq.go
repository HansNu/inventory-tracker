package models

// binding tags make the shape a 400 before any query runs. Without them a
// missing asset_code travels all the way to Postgres and comes back as a
// constraint violation, which is a 500-shaped answer to a client mistake.
type AddAssetReq struct {
	Id           int     `json:"id"`
	AssetCode    string  `json:"assetCode" binding:"required"`
	AssetName    string  `json:"assetName" binding:"required"`
	Brand        *string `json:"brand"`
	SerialNumber *string `json:"serialNumber"`
	CategoryId   int     `json:"categoryId" binding:"required"`
	Status       string  `json:"status" binding:"required"`
	Location     string  `json:"location" binding:"required"`
	User         *string `json:"user"`
	PurchaseDate *string `json:"purchaseDate"`
	Description  *string `json:"description"`
}
