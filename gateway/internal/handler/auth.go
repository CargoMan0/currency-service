package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (c *controller) Login(ctx *gin.Context) {
	var req loginRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		c.logger.Error("Error binding request parameters", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.authService.Login(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *controller) Register(ctx *gin.Context) {
	var req registerRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		c.logger.Error("Error binding request parameters", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.authService.Register(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}
