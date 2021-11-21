package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin_scaffiold/dao"
	"github.com/gin_scaffiold/middleware"
	"github.com/gin_scaffiold/public"
)

// HTTPJwtFlowCountMiddleware 租户日清流流量统计
func HTTPJwtFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}

		appInfo := appInterface.(*dao.App)
		appCounter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + appInfo.AppID)

		if err != nil {
			middleware.ResponseError(c, 2002, err)
			c.Abort()
			return
		}

		//原子增加
		appCounter.Increase()

		if appInfo.Qpd > 0 && appCounter.TotalCount > appInfo.Qpd {
			middleware.ResponseError(c, 2003, errors.New(fmt.Sprintf("租户日请求量限流 limit:%v current:%v", appInfo.Qpd, appCounter.TotalCount)))

			c.Abort()
			return
		}

		c.Next()
	}
}
