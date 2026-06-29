package models

type AddAssetReq struct {
	Id            int     `json:"id"`
	AssetCode     string  `json:"asset_code"`
	AssetName     string  `json:"asset_name"`
	Brand         *string `json:"brand"`
	SerialNumber  *string `json:"serial_number"`
	AssetCategory string  `json:"asset_category"`
	Status        string  `json:"status"`
	Location      string  `json:"location"`
	User          *string `json:"user"`
	PurchaseDate  *string `json:"purchase_date"`
	Description   *string `json:"description"`
}
