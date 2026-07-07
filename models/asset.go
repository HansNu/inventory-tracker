package models

import (
	"time"
)

type Asset struct {
	Id            int        `json:"id"`
	AssetCode     string     `json:"assetCode"`
	AssetName     string     `json:"assetName"`
	Brand         *string    `json:"brand"`
	SerialNumber  *string    `json:"serialNumber"`
	CategoryId    string     `json:"CategoryId"`
	AssetCategory string     `json:"AssetCategory"`
	Status        string     `json:"status"`
	Location      string     `json:"location"`
	User          *string    `json:"user"`
	PurchaseDate  *time.Time `json:"purchaseDate"`
	Description   *string    `json:"description"`
	CreateDt      time.Time  `json:"createDt"`
	UpdateDt      time.Time  `json:"updateDt"`
}
