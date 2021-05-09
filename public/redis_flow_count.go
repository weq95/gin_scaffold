package public

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin_scaffiold/common/lib"
	"sync/atomic"
	"time"
)

type RedisFlowCountService struct {
	AppID       string
	Interval    time.Duration
	QPS         int64
	Unix        int64
	TickerCount int64
	TotalCount  int64
}

func (r *RedisFlowCountService) GetDayKey(t time.Time) string {
	hourStr := t.In(lib.TimeLocaltion).Format("20060102")

	return fmt.Sprintf("%s_%s_%s", RedisFlowHourKey, hourStr, r.AppID)
}

func (r *RedisFlowCountService) GetHourKey(t time.Time) string {
	hourStr := t.In(lib.TimeLocaltion).Format("2006010215")

	return fmt.Sprintf("%s_%s_%s", RedisFlowHourKey, hourStr, r.AppID)
}

func (r *RedisFlowCountService) GetHourData(t time.Time) (int64, error) {
	return redis.Int64(RedisConfDo("GET", r.GetHourKey(t)))
}

func (r *RedisFlowCountService) GetDayData(t time.Time) (int64, error) {
	return redis.Int64(RedisConfDo("GET", r.GetDayKey(t)))
}

//原子增加
func (r *RedisFlowCountService) Increase() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		atomic.AddInt64(&r.TickerCount, 1)
	}()
}

func NewRedisFlowCountService(appId string, interval time.Duration) *RedisFlowCountService {
	reqCounter := &RedisFlowCountService{
		AppID:    appId,
		Interval: interval,
		QPS:      0,
		Unix:     0,
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		ticker := time.NewTicker(interval)
		for true {
			<-ticker.C
			tickerCount := atomic.LoadInt64(&reqCounter.TickerCount) //获取数据
			atomic.StoreInt64(&reqCounter.TickerCount, 0)            //重置数据

			currentTime := time.Now()
			dayKey := reqCounter.GetDayKey(currentTime)
			hourKey := reqCounter.GetHourKey(currentTime)
			if err := RedisConfPipline(func(c redis.Conn) {
				_ = c.Send("INCRBY", dayKey, tickerCount)
				_ = c.Send("EXPIRE", dayKey, 86400*2)
				_ = c.Send("INCRBY", hourKey, tickerCount)
				_ = c.Send("EXPIRE", hourKey, 86400*2)
			}); err != nil {
				fmt.Println("RedisConfPipline err", err)
				continue
			}

			totalCount, err := reqCounter.GetDayData(currentTime)
			if err != nil {
				fmt.Println("reqCounter.GetDayData err", err)
				continue
			}

			nowUnix := time.Now().Unix()
			if reqCounter.Unix == 0 {
				reqCounter.Unix = time.Now().Unix()
				continue
			}

			tickerCount = totalCount - reqCounter.TotalCount
			if nowUnix > reqCounter.Unix {
				reqCounter.TotalCount = totalCount
				reqCounter.QPS = tickerCount / (nowUnix - reqCounter.Unix)
				reqCounter.Unix = time.Now().Unix()
			}
		}
	}()

	return reqCounter
}
