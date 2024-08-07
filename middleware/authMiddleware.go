package middleware

import (
	"backend/service"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	jwtService service.JwtService
}

type UserIdContextKey struct{}

type RoleContextKey struct{}

func (m Middleware) RequireAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "require access token",
		})
		return
	}
	header := strings.Fields(authHeader)
	if len(header) != 2 && header[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong access token format",
		})
		return
	}

	accessToken := header[1]
	claims, err := m.jwtService.VerifyToken(accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx := context.WithValue(c.Request.Context(), UserIdContextKey{}, claims["userId"])
	ctx = context.WithValue(ctx, RoleContextKey{}, claims["role"])
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

func (Middleware) RequireAdminRole(c *gin.Context) {
	role := c.Request.Context().Value(RoleContextKey{})
	if role == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "role not found",
		})
		return
	}
	if role.(string) != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	c.Next()
}

func NewMiddleware(jwtService service.JwtService) Middleware {
	return Middleware{
		jwtService: jwtService,
	}
}
