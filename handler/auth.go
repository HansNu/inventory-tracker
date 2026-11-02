package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"inventory-tracker/models"
)

// @Summary  Register a new user and receive a JWT
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    user body models.RegisterReq true "New user"
// @Success  201 {object} map[string]interface{}
// @Failure  400 {object} map[string]string
// @Failure  409 {object} map[string]string
// @Router   /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}

	user, token, err := h.Auth.Register(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
		"token":   token,
	})
}

// @Summary  Authenticate and receive a JWT
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    credentials body models.LoginReq true "Username and password"
// @Success  200 {object} map[string]interface{}
// @Failure  400 {object} map[string]string
// @Failure  401 {object} map[string]string
// @Router   /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}

	user, token, err := h.Auth.Login(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// @Summary  Return the authenticated user
// @Tags     auth
// @Produce  json
// @Security BearerAuth
// @Success  200 {object} models.Users
// @Failure  401 {object} map[string]string
// @Failure  404 {object} map[string]string
// @Router   /me [get]
func (h *Handler) GetCurrentUser(c *gin.Context) {
	// The id comes from the verified token via AuthMiddleware, never from the
	// request. A client may say what it wants done; it may not say who it is.
	value, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := value.(int)
	if !ok {
		// Two-value assertion: the single-value form panics on a type mismatch,
		// which turns a wiring bug into a crashed request.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	user, err := h.Auth.CurrentUser(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
