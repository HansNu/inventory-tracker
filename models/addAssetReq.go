package models

type AddAssetReq struct {
	Id           int     `json:"id"`
	AssetCode    string  `json:"assetCode"`
	AssetName    string  `json:"assetName"`
	Brand        *string `json:"brand"`
	SerialNumber *string `json:"serialNumber"`
	CategoryId   int     `json:"categoryId"`
	Status       string  `json:"status"`
	Location     string  `json:"location"`
	User         *string `json:"user"`
	PurchaseDate *string `json:"purchaseDate"`
	Description  *string `json:"description"`
}
