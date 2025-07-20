package middleware

import (
	"net/http"
	"strings"

	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(maker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		fields := strings.Fields(auth)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or malformed auth header"})
			return
		}
		payload, err := maker.VerifyToken(fields[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("payload", payload) // available in handlers
		c.Next()
	}
}
