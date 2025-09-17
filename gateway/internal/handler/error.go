package handler

import (
	"errors"
	"github.com/BernsteinMondy/currency-service/gateway/internal/clients/auth"
	customerr "github.com/BernsteinMondy/currency-service/gateway/internal/errors"
	"github.com/BernsteinMondy/currency-service/gateway/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (s *controller) handleError(c *gin.Context, err error) {
	var nferr customerr.NotFoundError
	if errors.As(err, &nferr) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": nferr.Error(),
		})
	}

	log.Printf("internal error: %v", err)
	switch {
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	case errors.Is(err, service.ErrAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": "User already exist"})
	case errors.Is(err, service.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	case errors.Is(err, auth.ErrUnexpectedStatusCode):
		log.Printf("unexpected status code error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
	case errors.Is(err, auth.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	case errors.Is(err, auth.ErrTokenGeneration):
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate token"},
		)
	case errors.Is(err, auth.ErrTokenNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found"})
	case errors.Is(err, auth.ErrInvalidOrExpiredToken):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid or expired"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
