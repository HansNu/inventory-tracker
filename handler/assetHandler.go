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
	var asset models.AddAssetReq
	if err := c.ShouldBindBodyWithJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Printf("[handler] parsed asset: %+v\n", asset)
		return
	}

	err := h.Service.AddAsset(context.Background(), asset)
	fmt.Println("[handler] service call returned, err =", err)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Asset added successfully"})
}

func (h *Handler) GetAssetList(c *gin.Context) {
	params := models.AssetListParams{
		Page:          c.DefaultQuery("page", "1"),
		PageSize:      c.DefaultQuery("pageSize", "10"),
		Search:        c.DefaultQuery("search", ""),
		Status:        c.DefaultQuery("status", ""),
		AssetCategory: c.DefaultQuery("type", ""),
		SortField:     c.DefaultQuery("sortField", "purchase_date"),
		SortOrder:     c.DefaultQuery("sortOrder", "descend"),
	}

	assets, total, err := h.Service.GetAssetList(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.AssetListResponse{
		Data:  assets,
		Total: total,
	})
}

func (h *Handler) GetAssetByAssetCode(c *gin.Context) {
	assetCode := c.Param("assetCode")
	if assetCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "asset_code is required"})
		return
	}

	result, err := h.DB.Query(context.Background(),
		`Select a.id, a.asset_code, a.asset_name, a.category_id, c.category_name, a.brand, a.serial_number, a.status, a.location, a."user", a.purchase_date, a.description 
     			from asset a
				join category c on c.id = a.category_id
				where asset_code = $1`, assetCode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer result.Close()

	var a models.AssetResponse
	if !result.Next() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}

	if err := result.Scan(&a.Id, &a.AssetCode, &a.AssetName, &a.CategoryId, &a.CategoryName, &a.Brand, &a.SerialNumber, &a.Status, &a.Location, &a.User, &a.PurchaseDate, &a.Description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, a)
}

func (h *Handler) UpdateAsset(c *gin.Context) {
	var asset models.AddAssetReq

	if err := c.ShouldBindBodyWithJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.DB.Exec(context.Background(),
		`UPDATE asset SET asset_code = $1, asset_name = $2, brand = $3, serial_number = $4, category_id = $5, status = $6, location = $7, "user" = $8, purchase_date = $9, description = $10 WHERE id = $11`,
		asset.AssetCode, asset.AssetName, asset.Brand, asset.SerialNumber,
		asset.CategoryId, asset.Status, asset.Location, asset.User,
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

func (h *Handler) DeleteAssetByAssetCode(c *gin.Context) {
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
