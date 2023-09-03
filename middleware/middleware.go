package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mukulmantosh/ecommerce-gin/tokens"
	"net/http"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": fmt.Sprintf("No Authorization Header provided")})
			c.Abort()
			return
		}

		claims, err := tokens.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("uid", claims.UID)
		c.Next()

	}
}
