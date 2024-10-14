package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"health-server/internal/controller"
	"health-server/internal/controller/api/product"
	"health-server/internal/controller/api/system"
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

	userGroup := engine.Group("/api/users")
	{
		userGroup.POST("/login", user.Login) // 用户登录接口，不需要权限
	}

	authUserGroup := engine.Group("/api/users")
	authUserGroup.Use(AuthMiddleware()) // 应用权限中间件
	{
		authUserGroup.GET("/info", user.GetInfo) // 获取用户信息
	}

	systemGroup := engine.Group("/api/system")
	systemGroup.Use(AuthMiddleware())
	{
		systemGroup.GET("/info", system.Info) // 获取系统信息 包含添加剂信息和配置信息
	}

	// 产品相关路由
	productGroup := engine.Group("/api/products")
	productGroup.Use(AuthMiddleware())
	{
		productGroup.GET("/imgUrl", product.GetImgUrl) // 获取产品图片上传地址
		productGroup.POST("/upload", product.Upload)   // 上传商品信息
		productGroup.GET("/:id", product.Get)          // 获取指定 ID 的产品
	}
}
