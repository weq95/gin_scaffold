package public

import (
	"sync"
	"time"
)

type FlowCounter struct {
	RedisFlowCountMap   map[string]*RedisFlowCountService
	RedisFlowCountSlice []*RedisFlowCountService
	Locker              sync.RWMutex
}

var FlowCounterHandler *FlowCounter

func NewFlowCounter() *FlowCounter {
	return &FlowCounter{
		RedisFlowCountMap:   map[string]*RedisFlowCountService{},
		RedisFlowCountSlice: []*RedisFlowCountService{},
		Locker:              sync.RWMutex{},
	}
}

func init() {
	FlowCounterHandler = NewFlowCounter()
}

func (f *FlowCounter) GetCounter(serviceName string) (*RedisFlowCountService, error) {
	for _, service := range f.RedisFlowCountSlice {
		if service.AppID == serviceName {
			return service, nil
		}
	}

	newF := NewRedisFlowCountService(serviceName, 1*time.Second)
	f.RedisFlowCountSlice = append(f.RedisFlowCountSlice, newF)

	f.Locker.Lock()
	defer f.Locker.Unlock()

	f.RedisFlowCountMap[serviceName] = newF

	return newF, nil
}
