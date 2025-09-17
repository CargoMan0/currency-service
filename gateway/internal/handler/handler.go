package handler

import (
	"github.com/BernsteinMondy/currency-service/gateway/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type controller struct {
	authService     *service.AuthService
	currencyService *service.CurrencyService
	router          *gin.Engine
	logger          *zap.Logger
}

func RegisterRoutes(
	authService *service.AuthService,
	currencyService *service.CurrencyService,
	router *gin.Engine,
	logger *zap.Logger) controller {

	ctrl := controller{
		authService:     authService,
		currencyService: currencyService,
		router:          router,
		logger:          logger,
	}

	ctrl.router.GET(
		"/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		},
	)

	ctrl.router.GET("/api/v1/rate", ctrl.GetCurrencyRates)
	ctrl.router.POST("/api/v1/login", ctrl.Login)
	ctrl.router.POST("/api/v1/register", ctrl.Register)
	ctrl.router.POST("/api/v1/logout", ctrl.Logout)

	return ctrl
}
