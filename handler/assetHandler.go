package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"inventory-tracker/models"
)

// @Summary  Create an asset
// @Tags     assets
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    asset body models.AddAssetReq true "Asset to create"
// @Success  201 {object} map[string]string
// @Failure  400 {object} map[string]string
// @Failure  409 {object} map[string]string
// @Router   /addAsset [post]
func (h *Handler) AddAsset(c *gin.Context) {
	var asset models.AddAssetReq
	if err := c.ShouldBindJSON(&asset); err != nil {
		badRequest(c, err)
		return
	}

	if err := h.Service.AddAsset(c.Request.Context(), asset); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Asset added successfully"})
}

// @Summary  List assets, paged and filtered
// @Tags     assets
// @Produce  json
// @Security BearerAuth
// @Param    page      query int    false "Page number"                        default(1)
// @Param    pageSize  query int    false "Items per page (max 100)"           default(10)
// @Param    search    query string false "Free-text search across asset fields"
// @Param    status    query string false "Filter by status"
// @Param    type      query string false "Filter by category name"
// @Param    sortField query string false "Column to sort by"                  default(purchase_date)
// @Param    sortOrder query string false "ascend or descend"                  default(descend)
// @Success  200 {object} models.AssetListResponse
// @Failure  400 {object} map[string]string
// @Failure  401 {object} map[string]string
// @Router   /getAssetList [get]
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
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.AssetListResponse{Data: assets, Total: total})
}

// @Summary  Get one asset by its asset code
// @Tags     assets
// @Produce  json
// @Security BearerAuth
// @Param    assetCode path string true "Asset code"
// @Success  200 {object} models.AssetResponse
// @Failure  404 {object} map[string]string
// @Router   /getAssetByAssetCode/{assetCode} [get]
func (h *Handler) GetAssetByAssetCode(c *gin.Context) {
	assetCode := c.Param("assetCode")
	if assetCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "assetCode is required"})
		return
	}

	asset, err := h.Service.GetAssetByCode(c.Request.Context(), assetCode)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, asset)
}

// @Summary  Update an asset
// @Tags     assets
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id    path int                true "Asset ID"
// @Param    asset body models.AddAssetReq true "Updated fields"
// @Success  200 {object} map[string]string
// @Failure  400 {object} map[string]string
// @Failure  404 {object} map[string]string
// @Router   /updateAsset/{id} [put]
func (h *Handler) UpdateAsset(c *gin.Context) {
	// c.Param returns a string. A non-numeric id is the caller's mistake, so it
	// is a 400 — discarding this error turns it into a confusing 404 instead.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}

	var asset models.AddAssetReq
	if err := c.ShouldBindJSON(&asset); err != nil {
		badRequest(c, err)
		return
	}

	if err := h.Service.UpdateAsset(c.Request.Context(), id, asset); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset updated successfully"})
}

// @Summary  Delete an asset by its asset code
// @Tags     assets
// @Produce  json
// @Security BearerAuth
// @Param    assetCode path string true "Asset code"
// @Success  200 {object} map[string]string
// @Failure  404 {object} map[string]string
// @Router   /deleteAssetByAssetCode/{assetCode} [delete]
func (h *Handler) DeleteAssetByAssetCode(c *gin.Context) {
	assetCode := c.Param("assetCode")
	if assetCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "assetCode is required"})
		return
	}

	if err := h.Service.DeleteAssetByCode(c.Request.Context(), assetCode); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s deleted successfully", assetCode)})
}
