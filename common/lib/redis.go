package lib

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"time"
)

func RedisConnFactory(name string) (redis.Conn, error) {
	if ConfRedisMap == nil || ConfRedisMap.List == nil {
		return nil, errors.New("create redis conn fail")
	}

	for confName, cfg := range ConfRedisMap.List {
		if name != confName {
			continue
		}

		randHost := cfg.ProxyList[rand.Intn(len(cfg.ProxyList))]
		if cfg.ConnTimeout == 0 {
			cfg.ConnTimeout = 50
		}

		if cfg.ReadTimeout == 0 {
			cfg.ReadTimeout = 100
		}

		if cfg.WriteTimeout == 0 {
			cfg.WriteTimeout = 100
		}

		c, err := redis.Dial(
			"tcp",
			randHost,
			redis.DialConnectTimeout(time.Duration(cfg.ConnTimeout)*time.Millisecond),
			redis.DialReadTimeout(time.Duration(cfg.ReadTimeout)*time.Millisecond),
			redis.DialWriteTimeout(time.Duration(cfg.WriteTimeout)*time.Millisecond))

		if err != nil {
			return nil, err
		}

		if cfg.Password != "" {
			if _, err = c.Do("AUTH", cfg.Password); err != nil {
				_ = c.Close()

				return nil, err
			}
		}

		if cfg.Db != 0 {
			if _, err = c.Do("SELECT", cfg.Db); err != nil {
				return nil, err
			}
		}

		return c, nil
	}

	return nil, errors.New("create redis conn fail")
}

func RedisLogDo(ctx *TraceContext, conn redis.Conn, command string, args ...interface{}) (interface{}, error) {
	startExecTime := time.Now()
	reply, err := conn.Do(command, args)
	endExecTime := time.Now()

	if err != nil {
		Log.TagError(ctx, "_com_redis_failure", map[string]interface{}{
			"method":    command,
			"err":       err,
			"bind":      args,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})

		return reply, err
	}

	replyStr, _ := redis.String(reply, nil)
	Log.TagInfo(ctx, "_com_redis_success", map[string]interface{}{
		"method":    command,
		"bind":      args,
		"reply":     replyStr,
		"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
	})

	return reply, err
}

//通过配置 执行redis
func RedisConfDo(ctx *TraceContext, name, command string, args ...interface{}) (interface{}, error) {
	c, err := RedisConnFactory(name)
	if err != nil {
		Log.TagError(ctx, "_com_redis_failure", map[string]interface{}{
			"method": command,
			"err":    errors.New("RedisConnFactory_error:" + name),
			"bind":   args,
		})

		return nil, err
	}
	defer c.Close()

	startExecTime := time.Now()
	reply, err := c.Do(command, args...)
	endExecTime := time.Now()

	if err != nil {
		Log.TagError(ctx, "_com_redis_failure", map[string]interface{}{
			"method":    command,
			"err":       err,
			"bind":      args,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})

		return reply, err
	}

	replyStr, _ := redis.String(reply, nil)
	Log.TagInfo(ctx, "_com_redis_success", map[string]interface{}{
		"method":    command,
		"bind":      args,
		"reply":     replyStr,
		"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
	})

	return reply, err
}
