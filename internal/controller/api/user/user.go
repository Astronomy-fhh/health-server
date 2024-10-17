package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"health-server/internal/controller"
	"health-server/internal/kit"
	"health-server/internal/logger"
	"health-server/internal/mgr"
	"health-server/internal/model"
	"time"
)

type LoginReq struct {
	Token string `json:"token"`
}

type LoginResp struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
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
		logger.Logger.Info("user login by token", zap.String("token", tokenString))
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

		// 设置默认的昵称和头像
		defaultUser := mgr.GetUserDefaultMgr().GetDefaultUser(user.ID)
		user.Name = defaultUser.Name
		user.SystemAvatar = defaultUser.Img
		err = model.UpdateUser(user)
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
		Token:  token,
		Expire: tokenExp,
	}
	logger.Logger.Info("user login success", zap.Any("user", user))
	ctx.Success(resp)
}

func GetInfo(c *gin.Context) {
	ctx := controller.GetContext(c)
	token := ctx.MustGetToken()

	user, err := model.GetUserByUID(token.Uid)
	if err != nil {
		logger.Logger.Error("get user by uid failed", zap.Error(err))
		ctx.DefaultError()
		return
	}
	if user == nil {
		ctx.ParamError(errors.New("user not found"))
		return
	}

	resp := LoginUserInfo{
		Uid:          user.UID,
		Name:         user.Name,
		Did:          user.DID,
		SystemAvatar: user.SystemAvatar,
		CustomAvatar: user.CustomAvatar,
		BindID:       user.BindID,
	}

	logger.Logger.Info("get user info success", zap.Any("resp", resp))

	ctx.Success(resp)
}

func ChangeAvatar(c *gin.Context) {
	ctx := controller.GetContext(c)
	token := ctx.MustGetToken()

	var req struct {
		Index int `json:"index"`
	}
	err := ctx.GetReq(&req)
	if err != nil {
		ctx.ParamError(err)
		return
	}

	user, err := model.GetUserByUID(token.Uid)
	if err != nil {
		logger.Logger.Error("get user by uid failed", zap.Error(err))
		ctx.DefaultError()
		return
	}
	if user == nil {
		ctx.ParamError(errors.New("user not found"))
		return
	}

	defaultUser := mgr.GetUserDefaultMgr().GetDefaultUser(uint64(req.Index))
	user.SystemAvatar = defaultUser.Img
	user.Name = defaultUser.Name
	err = model.UpdateUser(user)
	if err != nil {
		logger.Logger.Error("update user failed", zap.Error(err))
		ctx.DefaultError()
		return
	}

	logger.Logger.Info("change user avatar", zap.String("uid", user.UID), zap.Any("defaultUser", defaultUser))

	resp := LoginUserInfo{
		Uid:          user.UID,
		Name:         user.Name,
		Did:          user.DID,
		SystemAvatar: user.SystemAvatar,
		CustomAvatar: user.CustomAvatar,
		BindID:       user.BindID,
	}

	ctx.Success(resp)
}
