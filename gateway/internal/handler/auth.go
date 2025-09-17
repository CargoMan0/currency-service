package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) Login(ctx *gin.Context) {
	var req loginRequest

	err := ctx.BindJSON(&req)
	if err != nil {
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

func (c *controller) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Authorization token is required"})
		return
	}

	err := c.authService.Logout(ctx.Request.Context(), token)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
