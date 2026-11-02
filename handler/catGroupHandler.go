package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"inventory-tracker/models"
)

// @Summary  List category groups
// @Tags     category-groups
// @Produce  json
// @Success  200 {array} models.CategoryGroup
// @Router   /getCategoryGroup [get]
func (h *Handler) GetCategoryGroup(c *gin.Context) {
	groups, err := h.CategoryGroups.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, groups)
}

// @Summary  Create a category group
// @Tags     category-groups
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    group body models.CategoryGroup true "Group to create"
// @Success  201 {object} map[string]string
// @Failure  400 {object} map[string]string
// @Failure  409 {object} map[string]string
// @Router   /addCategoryGroup [post]
func (h *Handler) AddCategoryGroup(c *gin.Context) {
	var group models.CategoryGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		badRequest(c, err)
		return
	}

	if err := h.CategoryGroups.Create(c.Request.Context(), group); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category group added successfully"})
}

// @Summary     Delete a category group
// @Description Requires the IT role. A group still referenced by categories returns 409.
// @Tags        category-groups
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Category group ID"
// @Success     200 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Failure     409 {object} map[string]string
// @Router      /deleteCategoryGroup/{id} [delete]
func (h *Handler) DeleteCategoryGroupById(c *gin.Context) {
	// The identifier moved from the request body into the path. DELETE with a
	// body has no defined meaning in HTTP and proxies are free to strip it.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}

	if err := h.CategoryGroups.DeleteByID(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category group deleted successfully"})
}
