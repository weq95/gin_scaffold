package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin_scaffiold/dao"
	"github.com/gin_scaffiold/middleware"
	"github.com/gin_scaffiold/public"
)

// HTTPJwtFlowLimitMiddleware qps统计
func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}

		appInfo := appInterface.(*dao.App)
		if appInfo.Qps > 0 {
			clientLimiter, err := public.FlowLimiterHandler.GetLimiter(public.FlowAppPrefix+appInfo.AppID+"_"+c.ClientIP(), float64(appInfo.Qps))

			if err != nil {
				middleware.ResponseError(c, 5001, err)
				c.Abort()
				return
			}

			if clientLimiter.Allow() {
				middleware.ResponseError(c, 5002, errors.New(fmt.Sprintf("%v flow limit %v", c.ClientIP(), appInfo.Qps)))

				c.Abort()
				return
			}
		}

		c.Next()
	}
}
