package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"inventory-tracker/models"
)

// *  means i get actual data not a copy of the data from the dbcontext
func (h *Handler) AddAsset(c *gin.Context) {
	var asset models.Asset

	// ShouldBindJSON reads the request body and maps it to the struct
	// using the json tags we defined (e.g. json:"asset_code")
	if err := c.ShouldBindBodyWithJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.DB.Exec(context.Background(),
		`INSERT INTO asset (asset_code, asset_name, brand, serial_number, asset_category, status, location, "user", purchase_date, description)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		asset.AssetCode, asset.AssetName, asset.Brand, asset.SerialNumber,
		asset.AssetCategory, asset.Status, asset.Location, asset.User,
		asset.PurchaseDate, asset.Description)

	if err != nil {
		// http.StatusInternalServerError = 500, means something went wrong on our end
		if strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Asset code already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// http.StatusCreated = 201, standard response for successful POST/create
	c.JSON(http.StatusCreated, gin.H{"message": "Asset added successfully"})
}

func (h *Handler) GetAssetList(c *gin.Context) {
	//in go assetList here represents direct connection to database, it doesnt work like models in JS/.NET
	var totalAsset int
	assetList, err := h.DB.Query(context.Background(),
		`Select id, asset_code, asset_name, asset_category, brand, serial_number, status, location, "user", purchase_date, description, COUNT(*) over()as total_count from asset`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer assetList.Close()

	//that's why you have to reassign all data queried from the db to a struct in GO, this is the data we're returning
	var assets []models.Asset
	for assetList.Next() {
		var a models.Asset
		err := assetList.Scan(&a.Id, &a.AssetCode, &a.AssetName, &a.AssetCategory, &a.Brand, &a.SerialNumber, &a.Status, &a.Location, &a.User, &a.PurchaseDate, &a.Description, &totalAsset)
		//& is to inject data directly into variables otherwise, struct stays empty
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		assets = append(assets, a)
	}

	c.JSON(http.StatusOK, models.AssetListResponse{
		Data:  assets,
		Total: totalAsset,
	})
}

func (h *Handler) UpdateAsset(c *gin.Context) {
	var asset models.Asset

	if err := c.ShouldBindBodyWithJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.DB.Exec(context.Background(),
		`UPDATE asset SET asset_code = $1, asset_name = $2, brand = $3, serial_number = $4, asset_category = $5, status = $6, location = $7, "user" = $8, purchase_date = $9, description = $10 WHERE id = $11`,
		asset.AssetCode, asset.AssetName, asset.Brand, asset.SerialNumber,
		asset.AssetCategory, asset.Status, asset.Location, asset.User,
		asset.PurchaseDate, asset.Description, asset.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Asset Id %d not found", asset.Id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset Updated Successfully"}) //sprintf returns a usable string like a variable
}

func (h *Handler) DeleteAsset(c *gin.Context) {
	var asset models.Asset

	if err := c.ShouldBindBodyWithJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.DB.Exec(context.Background(),
		`DELETE FROM asset where asset_code = $1`, asset.AssetCode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s Deleted Successfully", asset.AssetCode)})
}
