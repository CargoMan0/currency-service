package handler

import (
	"errors"
	serviceErrors "github.com/BernsteinMondy/currency-service/gateway/internal/service/errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (c *controller) handleError(ctx *gin.Context, err error) {
	log.Printf("internal error: %v", err)
	switch {
	case errors.Is(err, serviceErrors.ErrNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	case errors.Is(err, serviceErrors.ErrAlreadyExists):
		ctx.JSON(http.StatusConflict, gin.H{"error": "User already exist"})
	case errors.Is(err, serviceErrors.ErrInvalidCredentials):
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
