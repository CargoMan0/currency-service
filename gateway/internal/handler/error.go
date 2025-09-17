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

func (c *controller) handleError(ctx *gin.Context, err error) {
	var custom customerr.NotFoundError
	if errors.As(err, &custom) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": custom.Error(),
		})
	}

	log.Printf("internal error: %v", err)
	switch {
	case errors.Is(err, service.ErrNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	case errors.Is(err, service.ErrAlreadyExists):
		ctx.JSON(http.StatusConflict, gin.H{"error": "User already exist"})
	case errors.Is(err, service.ErrInvalidCredentials):
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	case errors.Is(err, auth.ErrUnexpectedStatusCode):
		log.Printf("unexpected status code error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
	case errors.Is(err, auth.ErrInvalidCredentials):
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	case errors.Is(err, auth.ErrTokenGeneration):
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate token"},
		)
	case errors.Is(err, auth.ErrTokenNotFound):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Token not found"})
	case errors.Is(err, auth.ErrInvalidOrExpiredToken):
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid or expired"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
