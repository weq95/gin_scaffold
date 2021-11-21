package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin_scaffiold/dao"
	"github.com/gin_scaffiold/middleware"
	"github.com/gin_scaffiold/public"
)

// HTTPFlowLimitMiddleware 客户端限流
func HTTPFlowLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		serviceInterface, ok := ctx.Get("service")
		if !ok {
			middleware.ResponseError(ctx, 2001, errors.New("service not found"))
			ctx.Abort()
			return
		}

		serviceDetail := serviceInterface.(*dao.ServiceDetail)
		if serviceDetail.AccessControl.ServiceFlowLimit != 0 {
			serviceLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				middleware.ResponseError(ctx, 5001, err)
				ctx.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				middleware.ResponseError(ctx, 5002, errors.New(fmt.Sprintf("service flow limit %v", serviceDetail.AccessControl.ServiceFlowLimit)))
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
