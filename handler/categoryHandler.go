package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"inventory-tracker/models"
)

// @Summary  List asset categories
// @Tags     categories
// @Produce  json
// @Security BearerAuth
// @Success  200 {array} models.AssetCategory
// @Router   /getAssetCategoryList [get]
func (h *Handler) GetAssetCategoryList(c *gin.Context) {
	categories, err := h.Categories.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, categories)
}

// @Summary  Create an asset category
// @Tags     categories
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    category body models.AssetCategory true "Category to create"
// @Success  201 {object} map[string]string
// @Failure  400 {object} map[string]string
// @Failure  409 {object} map[string]string
// @Router   /addAssetCategory [post]
func (h *Handler) AddAssetCategory(c *gin.Context) {
	var category models.AssetCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		badRequest(c, err)
		return
	}

	if err := h.Categories.Create(c.Request.Context(), category); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category added successfully"})
}

// @Summary  Delete an asset category
// @Tags     categories
// @Produce  json
// @Security BearerAuth
// @Param    id path int true "Category ID"
// @Success  200 {object} map[string]string
// @Failure  400 {object} map[string]string
// @Failure  404 {object} map[string]string
// @Router   /deleteAssetCategoryById/{id} [delete]
func (h *Handler) DeleteAssetCategoryById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}

	if err := h.Categories.DeleteByID(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
