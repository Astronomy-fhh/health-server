package main

import (
	"flag"
	"go.uber.org/zap"
	"health-server/config"
	"health-server/internal/controller_app"
	"health-server/internal/controller_web"
	"health-server/internal/db"
	"health-server/internal/http"
	"health-server/internal/kit"
	"health-server/internal/logger"
	"health-server/internal/mgr"
	"health-server/internal/s3"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 参数
	var configFile *string

	configFile = flag.String("config", "", "config file path")
	flag.Parse()

	if *configFile == "" {
		log.Fatal("config file is required")
	}

	// 初始化配置
	err := config.InitConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = logger.InitLog()
	if err != nil {
		log.Fatal(err)
	}
	logger.Logger.Info("InitConfig", zap.Any("config", config.Get()))

	// 初始化DB
	err = db.InitDB(config.Get().Db)
	if err != nil {
		log.Fatal(err)
	}

	// 等待关闭信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	runnerCtx := kit.NewRunnerContext()
	runner := kit.NewRunnerSlice()

	// 其他运行项
	runner.WithRunner(mgr.GetAdditiveMgr())
	runner.WithRunner(mgr.GetUserDefaultMgr())
	runner.WithRunner(mgr.GetAppMessageMgr())
	runner.WithRunner(s3.InitInstance(config.Get().S3))
	runner.WithRunner(http.NewHttpServer(http.ServerConfig{Port: config.Get().AppApiPort, Env: config.Get().Env}, controller_app.Routes))
	runner.WithRunner(http.NewHttpServer(http.ServerConfig{Port: config.Get().WebApiPort, Env: config.Get().Env}, controller_web.Routes))

	runner.Start(runnerCtx)
	logger.Logger.Info("server start")

	select {
	case <-quit:
		logger.Logger.Info("receive signal")
	case err := <-runnerCtx.Errored():
		logger.Logger.Error("runner error", zap.Error(err))
	}

	runner.Stop(runnerCtx)

	logger.Logger.Info("server shutdown")
	_ = logger.Logger.Sync()
}
