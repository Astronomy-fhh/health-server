package additive

import (
	"github.com/gin-gonic/gin"
	"health-server/internal/controller"
	"health-server/internal/mgr"
)

func GetAdditive(c *gin.Context) {
	ctx := controller.GetContext(c)
	additives := mgr.GetAdditiveMgr().GetAdditives()
	ctx.Success(additives)
}

func GetAdditiveCategory(c *gin.Context) {
	ctx := controller.GetContext(c)
	categories := mgr.GetAdditiveMgr().GetCategories()
	ctx.Success(categories)
}
