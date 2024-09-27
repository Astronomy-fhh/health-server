package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetHeader("X-User-Role") // 假设用户角色通过请求头传递
		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func Routes(engine *gin.Engine) {
	// 无权限控制的路由
	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello!")
	})

	// 定义用户路由组
	userGroup := engine.Group("/user")
	userGroup.Use(AuthMiddleware("user"))
	{
		userGroup.GET("/profile", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "User Profile"})
		})

		userGroup.GET("/settings", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "User Settings"})
		})
	}
}
