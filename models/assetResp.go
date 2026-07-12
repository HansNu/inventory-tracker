// Add this to your models or in a separate DTO file
package models

import (
	"time"
)

type AssetResponse struct {
	Id           int        `json:"id"`
	AssetCode    string     `json:"assetCode"`
	AssetName    string     `json:"assetName"`
	Brand        *string    `json:"brand"`
	SerialNumber *string    `json:"serialNumber"`
	CategoryId   int        `json:"categoryId"`
	CategoryName string     `json:"categoryName"`
	Status       string     `json:"status"`
	Location     string     `json:"location"`
	User         *string    `json:"user"`
	PurchaseDate *time.Time `json:"purchaseDate"`
	Description  *string    `json:"description"`
	CreateDt     time.Time  `json:"createDt"`
	UpdateDt     time.Time  `json:"updateDt"`
}
