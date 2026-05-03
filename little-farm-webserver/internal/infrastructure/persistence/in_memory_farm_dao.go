package persistence

import (
	"sync"

	"little-farm-webserver/internal/domain/littlefarm"
	"little-farm-webserver/pkg/domain"
)

type InMemoryFarmDAO struct {
	mt    sync.Mutex
	farms map[string]littlefarm.Farm
}

func NewInMemoryFarmDAO() *InMemoryFarmDAO {
	return &InMemoryFarmDAO{
		farms: map[string]littlefarm.Farm{},
	}
}

func (dao *InMemoryFarmDAO) Save(farm littlefarm.Farm) error {
	dao.mt.Lock()
	defer dao.mt.Unlock()

	dao.farms[farm.ID().String()] = farm
	return nil
}

func (dao *InMemoryFarmDAO) Find(id domain.ID) (littlefarm.Farm, bool, error) {
	dao.mt.Lock()
	defer dao.mt.Unlock()

	farm, found := dao.farms[id.String()]
	return farm, found, nil
}
