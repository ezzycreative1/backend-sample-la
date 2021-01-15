package middleware

import (
	"os"

	BaseHandler "backend-sample-la/app/base/handler"
	auth "backend-sample-la/app/user/usecase"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware for Authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := c.Request.Header.Get("secret-key")
		env := os.Getenv("ENV")
		if env != "local" && govalidator.IsNull(secretKey) {
			if secretKey != os.Getenv("SECRET_KEY") {
				c.AbortWithStatus(401)
			}
		}
		c.Next()
	}
}

// TokenAuthMiddleware ..
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c)
		if err != nil {
			BaseHandler.RespondUnauthorized(c, "")
			return
		}
		c.Next()
	}
}

// MaxAllowed specify max allowed connections
func MaxAllowed(n int) gin.HandlerFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(c *gin.Context) {
		acquire()       // before request
		defer release() // after request
		c.Next()
	}
}
