package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"health-server/internal/kit"
	"net/http"
)

const TokenKey = "token"

const (
	CodeSuccess    = 1
	CodeError      = 2
	CodeAuthError  = 3
	CodeParamError = 4
	CodeSysError   = 5
)

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (r *Resp) success(data any) {
	r.Code = CodeSuccess
	r.Msg = "ok"
	r.Data = data
}

func (r *Resp) errorWithCode(err error, code int) {
	r.Code = code
	r.Msg = err.Error()
}

type Context struct {
	ginCtx *gin.Context
	token  kit.UserTokenPayload
}

func GetContext(ginCtx *gin.Context) *Context {
	return &Context{ginCtx: ginCtx}
}

func (c *Context) GetToken() *kit.UserTokenPayload {
	value, exists := c.ginCtx.Get(TokenKey)
	if !exists {
		return nil
	}
	token, ok := value.(*kit.UserTokenPayload)
	if !ok {
		return nil
	}
	return token
}

func (c *Context) GetReq(obj interface{}) error {
	return c.ginCtx.ShouldBindBodyWithJSON(obj)
}

func (c *Context) Success(reply any) {
	resp := &Resp{}
	resp.success(reply)
	c.ginCtx.JSON(http.StatusOK, resp)
}

func (c *Context) Error(err error) {
	resp := &Resp{}
	resp.errorWithCode(err, CodeError)
	c.Success(resp)
}

func (c *Context) DefaultError() {
	resp := &Resp{}
	resp.errorWithCode(errors.New("unnamed error"), CodeError)
	c.Success(resp)
}

func (c *Context) AuthError() {
	resp := &Resp{}
	resp.errorWithCode(errors.New("auth failed"), CodeAuthError)
	c.ginCtx.JSON(http.StatusOK, resp)
}

func (c *Context) ParamError(err error) {
	resp := &Resp{}
	resp.errorWithCode(err, CodeParamError)
	c.ginCtx.JSON(http.StatusOK, resp)
}

func (c *Context) SysError(err error) {
	resp := &Resp{}
	resp.errorWithCode(err, CodeSysError)
	c.ginCtx.JSON(http.StatusOK, resp)
}
