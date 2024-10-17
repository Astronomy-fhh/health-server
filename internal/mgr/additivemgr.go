package mgr

import (
	"encoding/json"
	"go.uber.org/zap"
	"health-server/internal/kit"
	"health-server/internal/logger"
	"health-server/internal/model"
	"sync"
	"time"
)

type AdditiveMgr struct {
	sync.RWMutex
	cancelTick   chan struct{}
	additivesMap map[uint64]*Additive
	tagsMap      map[uint64]*model.AdditiveTag
}

var additiveMgrInstance *AdditiveMgr
var additiveMgrOnce sync.Once

func GetAdditiveMgr() *AdditiveMgr {
	additiveMgrOnce.Do(func() {
		additiveMgrInstance = &AdditiveMgr{
			cancelTick:   make(chan struct{}),
			additivesMap: make(map[uint64]*Additive),
			tagsMap:      make(map[uint64]*model.AdditiveTag),
		}
	})
	return additiveMgrInstance
}

type Additive struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	GB        string `json:"gb"`
	Status    string `json:"status"`
	Category  string `json:"category"`
	Tags      []int  `json:"tags"`
	ImageURL  string `json:"image_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (a *AdditiveMgr) Start(ctx *kit.RunnerContext) {
	err := a.Load()
	if err != nil {
		logger.Logger.Error("start load additives failed", zap.Error(err))
		ctx.Error(err)
		return
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
	additivesMap := make(map[uint64]*Additive)
	for _, additive := range additives {
		tags := make([]int, 0)
		if additive.Tags != nil {
			err := json.Unmarshal(additive.Tags, &tags)
			if err != nil {
				logger.Logger.Error("unmarshal tags failed", zap.Error(err), zap.String("tags", string(additive.Tags)))
				return err
			}
		}
		additivesMap[additive.ID] = &Additive{
			ID:        additive.ID,
			Name:      additive.Name,
			Desc:      additive.Desc,
			GB:        additive.GB,
			Category:  additive.Category,
			Tags:      tags,
			ImageURL:  additive.ImageURL,
			CreatedAt: additive.CreatedAt.String(),
			UpdatedAt: additive.UpdatedAt.String(),
		}
	}
	a.additivesMap = additivesMap

	tags, err := model.GetAllAdditiveTags()
	if err != nil {
		return err
	}

	tagsMap := make(map[uint64]*model.AdditiveTag)
	for _, category := range tags {
		tagsMap[category.ID] = category
	}
	a.tagsMap = tagsMap

	logger.Logger.Info("load additives success")

	return nil
}

func (a *AdditiveMgr) Stop(ctx *kit.RunnerContext) {
	close(a.cancelTick)
}

func (a *AdditiveMgr) GetAdditives() map[uint64]*Additive {
	a.RLock()
	defer a.RUnlock()
	return a.additivesMap
}

func (a *AdditiveMgr) GetTags() map[uint64]*model.AdditiveTag {
	a.RLock()
	defer a.RUnlock()
	return a.tagsMap
}
