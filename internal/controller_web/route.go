package controller_web

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"health-server/internal/controller_app"
	"health-server/internal/logger"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := controller_app.GetContext(c)
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			logger.Logger.Error("token is required")
			ctx.AuthError()
			c.Abort()
			return
		}

		payload, expired, err := ParseAdminToken(tokenString)
		if err != nil {
			logger.Logger.Error("parse token failed", zap.Error(err), zap.String("token", tokenString))
			ctx.AuthError()
			c.Abort()
			return
		}
		if expired {
			logger.Logger.Warn("token expired", zap.String("token", tokenString))
			ctx.AuthError()
			c.Abort()
			return
		}
		c.Set(controller_app.TokenKey, payload)
		c.Next()
	}
}

func Routes(engine *gin.Engine) {
	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello!")
	})

	userGroup := engine.Group("/api/admins")
	{
		userGroup.POST("/login", Login)
	}

	systemGroup := engine.Group("/api/system")
	systemGroup.Use(AuthMiddleware())
	{
		systemGroup.GET("/info", Info)
	}

	productGroup := engine.Group("/api/products")
	productGroup.Use(AuthMiddleware())
	{
		productGroup.POST("/reviews/get", GetReviewList)
		productGroup.POST("/get", GetProduct)
		productGroup.POST("/add", AddProduct)
	}
}
