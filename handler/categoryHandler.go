package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"fmt"
	"inventory-tracker/models"
)

func (h *Handler) GetAssetCategoryList(c *gin.Context) {
	//in go assetList here represents direct connection to database, it doesnt work like models in JS/.NET
	assetCategoryList, err := h.DB.Query(context.Background(),
		`Select id, category_name, category_group from category`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//that's why you have to reassign all data queried from the db to a struct in GO, this is the data we're returning
	var assetCategory []models.AssetCategory
	for assetCategoryList.Next() {
		var a models.AssetCategory
		err := assetCategoryList.Scan(&a.Id, &a.CategoryName, &a.CategoryGroup)
		//& is to inject data directly into variables otherwise, struct stays empty
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		assetCategory = append(assetCategory, a)
	}

	c.JSON(http.StatusOK, assetCategory)
}

func (h *Handler) AddAssetCategory(c *gin.Context) {
	var category models.AssetCategory

	if err := c.ShouldBindBodyWithJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.DB.Exec(context.Background(),
		`INSERT INTO category (category_name, category_group)
		VALUES ($1, $2)`,
		category.CategoryName, category.CategoryGroup)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category Added Successfully"})
}

func (h *Handler) DeleteAssetCategoryById(c *gin.Context) {
	var category models.AssetCategory

	if err := c.ShouldBindBodyWithJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.DB.Exec(context.Background(),
		`DELETE FROM category where id = $1`, category.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s Deleted Successfully", category.CategoryName)})
}
