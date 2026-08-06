package handler

import (
	"context"
	"net/http"

	"fmt"
	"inventory-tracker/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCategoryGroup(c *gin.Context) {
	catGroupList, err := h.DB.Query(context.Background(),
		`select id, group_name from category_group`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var catGroup []models.CategoryGroup
	for catGroupList.Next() {
		var a models.CategoryGroup
		err := catGroupList.Scan(&a.Id, &a.GroupName)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		catGroup = append(catGroup, a)
	}

	c.JSON(http.StatusOK, catGroup)
}

func (h *Handler) AddCategoryGroup(c *gin.Context) {
	var catGroup models.CategoryGroup

	if err := c.ShouldBindBodyWithJSON(&catGroup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err := h.DB.Exec(context.Background(),
		`INSERT into category_group (group_name)
		VALUES ($1)`, catGroup.GroupName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) DeleteCategoryGroupById(c *gin.Context) {
	var category models.AssetCategory

	if err := c.ShouldBindBodyWithJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.DB.Exec(context.Background(),
		`DELETE FROM category_group where id = $1`, category.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s Deleted Successfully", category.CategoryName)})
}
