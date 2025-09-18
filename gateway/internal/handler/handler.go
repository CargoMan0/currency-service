package handler

import (
	"context"
	"github.com/BernsteinMondy/currency-service/gateway/internal/models/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type controller struct {
	authService     AuthService
	currencyService CurrencyService
	router          *gin.Engine
	logger          *zap.Logger
}

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	Register(ctx context.Context, login, password string) error
}

type CurrencyService interface {
	GetCurrencyRates(ctx context.Context, request dto.ParsedCurrencyRequest) (*dto.CurrencyResponse, error)
}

func RegisterRoutes(
	authService AuthService,
	currencyService CurrencyService,
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

	return ctrl
}
