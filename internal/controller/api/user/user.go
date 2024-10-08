package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"health-server/internal/controller"
	"health-server/internal/kit"
	"health-server/internal/logger"
	"health-server/internal/model"
	"time"
)

type LoginReq struct {
	Token string `json:"token"`
}

type LoginResp struct {
	UserInfo *LoginUserInfo `json:"userInfo"`
	Token    string         `json:"token"`
	Expire   int64          `json:"expire"`
}

type LoginUserInfo struct {
	Uid          string `json:"uid"`
	Name         string `json:"name"`
	Did          string `json:"did"`
	SystemAvatar string `json:"system_avatar"`
	CustomAvatar string `json:"custom_avatar"`
	BindID       string `json:"bind_id"`
}

func Login(c *gin.Context) {
	ctx := controller.GetContext(c)

	tokenString := c.GetHeader("Authorization")

	var user *model.User
	// 已经注册的用户 找到该账号信息 返回
	if tokenString != "" {
		tokenInfo, _, err := kit.ParseUserToken(tokenString)
		if err != nil {
			logger.Logger.Error("parse user token failed", zap.Error(err))
			ctx.ParamError(errors.New("invalid token"))
			return
		}
		user, err = model.GetUserByUID(tokenInfo.Uid)
		if err != nil {
			logger.Logger.Error("get user by uid failed", zap.Error(err))
			ctx.DefaultError()
			return
		}
		if user == nil {
			ctx.ParamError(errors.New("user not found"))
			return
		}
		logger.Logger.Info("user login by uid", zap.String("uid", user.UID))
	} else {
		uid := model.GenUID()
		if uid == "" {
			logger.Logger.Error("gen uid failed")
			ctx.DefaultError()
			return
		}

		nowTime := time.Now().UTC()
		user = &model.User{
			UID:          uid,
			DID:          "",
			RegisteredAt: nowTime,
			LastLoginAt:  nowTime,
		}
		err := model.CreateUser(user)
		if err != nil {
			ctx.Error(err)
			return
		}
		logger.Logger.Info("new user created", zap.String("uid", user.UID), zap.String("did", user.DID))
	}

	tokenExp := time.Now().Add(time.Hour * 24 * 30).Unix()
	token, err := kit.GenerateUserToken(&kit.UserTokenPayload{Uid: user.UID}, tokenExp)
	if err != nil {
		logger.Logger.Error("generate user token failed", zap.Error(err))
		ctx.DefaultError()
		return
	}

	resp := LoginResp{
		UserInfo: &LoginUserInfo{
			Uid:          user.UID,
			Name:         user.Name,
			Did:          user.DID,
			SystemAvatar: user.SystemAvatar,
			CustomAvatar: user.CustomAvatar,
			BindID:       user.BindID,
		},
		Token:  token,
		Expire: tokenExp,
	}

	logger.Logger.Info("user login success", zap.Any("resp.userInfo", resp.UserInfo))

	ctx.Success(resp)
}
