package http

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"health-server/internal/kit"
	"health-server/internal/logger"
	"net/http"
	"time"
)

type Server struct {
	config ServerConfig
	engine *gin.Engine
	srv    *http.Server
}

type ServerConfig struct {
	Port         string
	Env          string
	CloseTimeout time.Duration
}

type RouteFunc func(*gin.Engine)

func NewHttpServer(config ServerConfig, routes RouteFunc) *Server {
	if config.CloseTimeout == 0 {
		config.CloseTimeout = 5 * time.Second
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(CustomLogger())

	if routes != nil {
		routes(engine)
	}

	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: engine,
	}

	return &Server{
		config: config,
		engine: engine,
		srv:    srv,
	}
}

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		latency := time.Since(startTime)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		logger.Logger.Sugar().Infof("| Http Request | %s | %d | %s | %s | %s", method, status, latency, path, c.ClientIP())
	}
}

func (s *Server) Start(ctx *kit.RunnerContext) {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			ctx.Error(err)
		}
	}()
	logger.Logger.Info("HTTP Server started", zap.String("port", s.config.Port))
}

func (s *Server) Stop(ctx *kit.RunnerContext) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), s.config.CloseTimeout)
	defer cancel()

	if err := s.srv.Shutdown(timeoutCtx); err != nil {
		logger.Logger.Error("HTTP Server close error:", zap.Error(err))
		return
	}

	logger.Logger.Info("HTTP Server closed")
}
