package logger

import "go.uber.org/zap"

var Logger *zap.Logger

func InitLog() error {
	log, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	Logger = log
	return nil
}
