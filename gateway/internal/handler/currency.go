package handler

import (
	"github.com/BernsteinMondy/currency-service/gateway/internal/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (c *controller) GetCurrencyRates(ctx *gin.Context) {
	var req currencyRequest

	err := ctx.BindQuery(&req)
	if err != nil {
		c.logger.Error("Error binding request parameters", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dateFrom, err := time.Parse("2006-01-02", req.DateFrom)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for date_from, expected YYYY-MM-DD"})
		return
	}

	dateTo, err := time.Parse("2006-01-02", req.DateTo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for date_to, expected YYYY-MM-DD"})
		return
	}

	parsedCurrencyRequest := dto.ParsedCurrencyRequest{
		Currency: req.Currency,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	data, err := c.currencyService.GetCurrencyRates(ctx.Request.Context(), parsedCurrencyRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}
