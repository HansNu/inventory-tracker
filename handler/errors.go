package handler

import (
	"errors"
	"log"
	"net/http"

	"inventory-tracker/apperr"

	"github.com/gin-gonic/gin"
)

// respondError is the only place in this package that turns an error into a
// status code. Services and repositories return domain errors; the mapping from
// those to HTTP lives here and nowhere else.
//
// The default branch logs the real error and returns a sentence. Returning
// err.Error() to a caller leaks internals; logging nothing leaves you debugging
// a 500 with no information, which costs hours. Both, in the right directions.
func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperr.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.Is(err, apperr.ErrDuplicate):
		c.JSON(http.StatusConflict, gin.H{"error": "already exists"})
	case errors.Is(err, apperr.ErrReferenced):
		c.JSON(http.StatusConflict, gin.H{"error": "still referenced by other records"})
	case errors.Is(err, apperr.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	case errors.Is(err, apperr.ErrInvalidSortField):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sortField"})
	default:
		log.Printf("[%s %s] unhandled error: %v", c.Request.Method, c.Request.URL.Path, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

// badRequest is for shape problems found before anything downstream is called.
// Binding errors are safe to echo: they describe the caller's own payload.
func badRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
