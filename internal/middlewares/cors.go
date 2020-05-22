package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jz222/loggy/internal/keys"
)

// Cors creates a new middleware that enables cors.
func Cors(c *gin.Context) {
	origin := keys.GetKeys().CLIENT_URL
	path := c.Request.URL.Path

	if path == "/logging/error" || path == "/logging/analytics" {
		origin = "*"
	}

	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, GET, PUT, OPTIONS")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
