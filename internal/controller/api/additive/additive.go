package additive

import (
	"github.com/gin-gonic/gin"
	"health-server/internal/controller"
	"health-server/internal/mgr"
)

func GetAdditive(c *gin.Context) {
	ctx := controller.GetContext(c)
	token := ctx.GetToken()
	if token == nil {
		ctx.AuthError()
		return
	}
	additives := mgr.GetAdditiveMgr().GetAdditives()
	ctx.Success(additives)
}

func GetAdditiveCategory(c *gin.Context) {
	ctx := controller.GetContext(c)
	token := ctx.GetToken()
	if token == nil {
		ctx.AuthError()
		return
	}
	categories := mgr.GetAdditiveMgr().GetCategories()
	ctx.Success(categories)
}
