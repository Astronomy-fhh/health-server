package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"health-server/internal/controller"
	"health-server/internal/controller/api/additive"
	"health-server/internal/controller/api/user"
	"health-server/internal/kit"
	"health-server/internal/logger"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := controller.GetContext(c)
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			logger.Logger.Error("token is required")
			ctx.AuthError()
			c.Abort()
			return
		}

		payload, expired, err := kit.ParseUserToken(tokenString)
		if err != nil {
			logger.Logger.Error("parse token failed", zap.Error(err), zap.String("token", tokenString))
			ctx.AuthError()
			c.Abort()
			return
		}
		if expired {
			logger.Logger.Warn("token expired", zap.String("token", tokenString))
		}
		c.Set(controller.TokenKey, payload)
		c.Next()
	}
}

func Routes(engine *gin.Engine) {
	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello!")
	})

	userGroup := engine.Group("/app/user")
	{
		userGroup.POST("/login", user.Login)
	}

	itemGroup := engine.Group("/app/item")
	itemGroup.Use(AuthMiddleware())
	{
		itemGroup.GET("/additive/get", additive.GetAdditive)
		itemGroup.GET("/additive_category/get", additive.GetAdditiveCategory)
	}
}
