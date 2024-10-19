package mgr

import (
	"go.uber.org/zap"
	"health-server/internal/kit"
	"health-server/internal/logger"
	"health-server/internal/model"
	"sync"
	"time"
)

type AppMessageMgr struct {
	sync.RWMutex
	cancelTick chan struct{}
	messages   map[string]string
}

var appMessageMgrInstance *AppMessageMgr
var appMessageMgrInstanceOnce sync.Once

func GetAppMessageMgr() *AppMessageMgr {
	appMessageMgrInstanceOnce.Do(func() {
		appMessageMgrInstance = &AppMessageMgr{
			cancelTick: make(chan struct{}),
			messages:   make(map[string]string),
		}
	})
	return appMessageMgrInstance
}

func (u *AppMessageMgr) Start(ctx *kit.RunnerContext) {
	err := u.Load()
	if err != nil {
		logger.Logger.Error("start load user defaults failed", zap.Error(err))
		ctx.Error(err)
		return
	}
	go u.TickLoad()
}

func (u *AppMessageMgr) TickLoad() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-u.cancelTick:
			return
		case <-ticker.C:
			err := u.Load()
			if err != nil {
				logger.Logger.Error("tick load user defaults failed", zap.Error(err))
			}
		}
	}
}

func (u *AppMessageMgr) Load() error {

	messages, err := model.GetAllAppMessages()
	if err != nil {
		return err
	}
	messagesMap := make(map[string]string)
	for _, message := range messages {
		messagesMap[message.Name] = message.Message
	}

	u.Lock()
	u.messages = messagesMap
	defer u.Unlock()

	logger.Logger.Info("load user defaults success")

	return nil
}

func (u *AppMessageMgr) Stop(ctx *kit.RunnerContext) {
	close(u.cancelTick)
}

func (u *AppMessageMgr) GetMessages() map[string]string {
	u.RLock()
	defer u.RUnlock()
	return u.messages
}
