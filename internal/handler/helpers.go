package handler

import (
	"net/http"

	apierrors "github.com/azizjon-top/techno-re-ecosystem/internal/errors"
	"github.com/gin-gonic/gin"
)

// respondError writes a structured error response based on *apierrors.APIError or a generic 500.
func respondError(c *gin.Context, err error) {
	if apiErr, ok := err.(*apierrors.APIError); ok {
		c.JSON(apiErr.StatusCode, apiErr)
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    apierrors.ErrCodeInternalError,
		"message": "internal server error",
	})
}
