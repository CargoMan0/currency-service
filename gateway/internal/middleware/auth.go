package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type AuthorizationClient interface {
	ValidateToken(ctx context.Context, token string) error
}

type Authorization struct {
	authClient  AuthorizationClient
	skipperFunc func(c *gin.Context) bool
	logger      *zap.Logger
}

func NewAuthorization(
	authClient AuthorizationClient,
	logger *zap.Logger,
	skipperFunc func(c *gin.Context) bool,
) *Authorization {
	return &Authorization{
		authClient:  authClient,
		skipperFunc: skipperFunc,
		logger:      logger,
	}
}

func (a *Authorization) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		if a.skipperFunc(c) {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if err := a.authClient.ValidateToken(c.Request.Context(), authHeaderParts[1]); err != nil {
			a.logger.Error(
				"Invalid token",
				zap.String("token", authHeaderParts[1]),
				zap.String("client_ip", c.ClientIP()),
				zap.String("user_agent", c.GetHeader("User-Agent")),
				zap.Error(err),
			)

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}

func shouldSkipEndpoint(c *gin.Context) bool {
	if c.Request.URL.Path == "/login" {
		return true
	}
	return false
}
