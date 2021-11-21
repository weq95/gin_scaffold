package public

import (
	"golang.org/x/time/rate"
	"sync"
)

var FlowLimiterHandler *FlowLimiter

type FlowLimiterItem struct {
	ServiceName string
	Limter      *rate.Limiter
}

type FlowLimiter struct {
	FlowLmiterMap   map[string]*FlowLimiterItem
	FlowLmiterSlice []*FlowLimiterItem
	Locker          sync.RWMutex
}


func NewFlowLimiter() *FlowLimiter {
	return &FlowLimiter{
		FlowLmiterMap:   map[string]*FlowLimiterItem{},
		FlowLmiterSlice: []*FlowLimiterItem{},
		Locker:          sync.RWMutex{},
	}
}

func init() {
	FlowLimiterHandler = NewFlowLimiter()
}

func (counter *FlowLimiter) GetLimiter(servieName string, qps float64) (*rate.Limiter, error) {
	for _, item := range counter.FlowLmiterSlice {
		if item.ServiceName == servieName {
			return item.Limter, nil
		}
	}

	newLister := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	item := &FlowLimiterItem{
		ServiceName: servieName,
		Limter:      newLister,
	}

	counter.FlowLmiterSlice = append(counter.FlowLmiterSlice, item)
	counter.Locker.Unlock()
	defer counter.Locker.Unlock()

	counter.FlowLmiterMap[servieName] = item

	return newLister, nil
}
