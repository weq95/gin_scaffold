package public

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin_scaffiold/common/lib"
)

func RedisConfPipline(pip ...func(c redis.Conn)) error {
	c, err := lib.RedisConnFactory("default")
	if err != nil {
		return err
	}

	defer func() {
		_ = c.Close()
	}()

	for _, f := range pip {
		f(c)
	}

	_ = c.Flush()

	return nil
}

func RedisConfDo(commandName string, args ...interface{}) (interface{}, error) {
	c, err := lib.RedisConnFactory("default")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = c.Close()
	}()

	return c.Do(commandName, args...)
}
