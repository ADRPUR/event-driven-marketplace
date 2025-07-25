package utils

import (
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
	"github.com/gin-gonic/gin"
)

// ExtractPayload returns *token.Payload from Gin context if present, else nil.
func ExtractPayload(c *gin.Context) *token.Payload {
	v, ok := c.Get("payload")
	if !ok || v == nil {
		return nil
	}
	pl, ok := v.(*token.Payload)
	if !ok {
		return nil
	}
	return pl
}
