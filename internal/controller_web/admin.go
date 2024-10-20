package controller_web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"health-server/internal/model"
	"time"
)

type LoginReq struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type LoginResp struct {
	Token  string     `json:"token"`
	Expire string     `json:"expire"`
	Admin  LoginAdmin `json:"admin"`
}

type LoginAdmin struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func Login(c *gin.Context) {
	ctx := GetContext(c)
	var req LoginReq
	if err := ctx.GetReq(&req); err != nil {
		ctx.ParamError(err)
		return
	}
	admin, err := model.GetAdminByName(req.Name)
	if err != nil {
		ctx.DefaultError()
		return
	}
	if admin == nil {
		ctx.Error(errors.New("admin not found"))
		return
	}
	if admin.Pass != req.Pass {
		ctx.Error(errors.New("password error"))
		return
	}

	tokenPayload := TokenPayload{
		AdminId: admin.ID,
	}
	token, err := GenerateAdminToken(tokenPayload, int64(Exp))
	if err != nil {
		ctx.DefaultError()
		return
	}

	resp := LoginResp{
		Token:  token,
		Expire: time.Now().Add(Exp).Format(time.RFC3339),
		Admin: LoginAdmin{
			Name:   admin.Name,
			Avatar: admin.Avatar,
		},
	}
	ctx.Success(resp)
}
