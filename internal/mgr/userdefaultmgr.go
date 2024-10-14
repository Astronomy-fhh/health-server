package mgr

import (
	"go.uber.org/zap"
	"health-server/internal/kit"
	"health-server/internal/logger"
	"health-server/internal/model"
	"sync"
	"time"
)

type UserDefaultMgr struct {
	sync.RWMutex
	cancelTick   chan struct{}
	defaultUsers []*model.UserDefault
}

var userDefaultMgrInstance *UserDefaultMgr
var userDefaultMgrOnce sync.Once

func GetUserDefaultMgr() *UserDefaultMgr {
	userDefaultMgrOnce.Do(func() {
		userDefaultMgrInstance = &UserDefaultMgr{
			cancelTick:   make(chan struct{}),
			defaultUsers: make([]*model.UserDefault, 0),
		}
	})
	return userDefaultMgrInstance
}

func (u *UserDefaultMgr) Start(ctx *kit.RunnerContext) {
	err := u.Load()
	if err != nil {
		logger.Logger.Error("start load user defaults failed", zap.Error(err))
		ctx.Error(err)
		return
	}
	go u.TickLoad()
}

func (u *UserDefaultMgr) TickLoad() {
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

func (u *UserDefaultMgr) Load() error {

	userDefaults, err := model.GetUserDefaults()
	if err != nil {
		return err
	}

	u.Lock()
	u.defaultUsers = userDefaults
	defer u.Unlock()

	logger.Logger.Info("load user defaults success")

	return nil
}

func (u *UserDefaultMgr) Stop(ctx *kit.RunnerContext) {
	close(u.cancelTick)
}

func (u *UserDefaultMgr) GetDefaultUser(uid uint64) model.UserDefault {
	u.RLock()
	defer u.RUnlock()

	userLen := len(u.defaultUsers)
	if userLen == 0 {
		return model.UserDefault{
			ID:   0,
			Name: "",
			Img:  "",
		}
	}
	userDefault := u.defaultUsers[uid%uint64(userLen)]
	return *userDefault
}
