package mgr

import (
	"go.uber.org/zap"
	"health-server/internal/kit"
	"health-server/internal/logger"
	"health-server/internal/model"
	"sync"
	"time"
)

type AdditiveMgr struct {
	sync.RWMutex
	cancelTick    chan struct{}
	additivesMap  map[uint64]*model.Additive
	categoriesMap map[uint64]*model.AdditiveCategory
}

var additiveMgrInstance *AdditiveMgr
var additiveMgrOnce sync.Once

func GetAdditiveMgr() *AdditiveMgr {
	additiveMgrOnce.Do(func() {
		additiveMgrInstance = &AdditiveMgr{
			cancelTick:    make(chan struct{}),
			additivesMap:  make(map[uint64]*model.Additive),
			categoriesMap: make(map[uint64]*model.AdditiveCategory),
		}
	})
	return additiveMgrInstance
}

func (a *AdditiveMgr) Start(ctx *kit.RunnerContext) {
	err := a.Load()
	if err != nil {
		logger.Logger.Error("start load additives failed", zap.Error(err))
		ctx.Error(err)
	}
	go a.TickLoad()
}

func (a *AdditiveMgr) TickLoad() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-a.cancelTick:
			return
		case <-ticker.C:
			err := a.Load()
			if err != nil {
				logger.Logger.Error("tick load additives failed", zap.Error(err))
			}
		}
	}
}

func (a *AdditiveMgr) Load() error {
	a.Lock()
	defer a.Unlock()
	additives, err := model.GetAllAdditives()
	if err != nil {
		return err
	}
	additivesMap := make(map[uint64]*model.Additive)
	for _, additive := range additives {
		additivesMap[additive.ID] = additive
	}
	a.additivesMap = additivesMap

	categories, err := model.GetAllAdditiveCategories()
	if err != nil {
		return err
	}

	categoriesMap := make(map[uint64]*model.AdditiveCategory)
	for _, category := range categories {
		categoriesMap[category.ID] = category
	}
	a.categoriesMap = categoriesMap

	logger.Logger.Info("load additives success")

	return nil
}

func (a *AdditiveMgr) Stop(ctx *kit.RunnerContext) {
	close(a.cancelTick)
}

func (a *AdditiveMgr) GetAdditives() map[uint64]*model.Additive {
	a.RLock()
	defer a.RUnlock()
	return a.additivesMap
}

func (a *AdditiveMgr) GetCategories() map[uint64]*model.AdditiveCategory {
	a.RLock()
	defer a.RUnlock()
	return a.categoriesMap
}
