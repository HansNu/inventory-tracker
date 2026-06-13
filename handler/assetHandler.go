package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"inventory-tracker/models"
)

var DB *pgx.Conn

// *  means i get actual data not a copy of the data from the dbcontext
func AddAsset(c *gin.Context) {
	var asset models.Asset

	// ShouldBindJSON reads the request body and maps it to the struct
	// using the json tags we defined (e.g. json:"asset_code")
	if err := c.ShouldBindBodyWithJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := DB.Exec(context.Background(),
		`INSERT INTO asset (asset_code, asset_name, brand, serial_number, asset_category, status, location, "user", purchase_date, description)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		asset.AssetCode, asset.AssetName, asset.Brand, asset.SerialNumber,
		asset.AssetCategory, asset.Status, asset.Location, asset.User,
		asset.PurchaseDate, asset.Description,
	)
	if err != nil {
		// http.StatusInternalServerError = 500, means something went wrong on our end
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// http.StatusCreated = 201, standard response for successful POST/create
	c.JSON(http.StatusCreated, gin.H{"message": "Asset added successfully"})
}
