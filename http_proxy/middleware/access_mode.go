package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin_scaffiold/dao"
	"github.com/gin_scaffiold/middleware"
)

// HTTPAccessModeMiddleware 匹配接入方式,基于请求信息
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, err := dao.ServiceManagerHandler.HTTPAccessMode(ctx)
		if err != nil {
			middleware.ResponseError(ctx, 1001, err)
			ctx.Abort()

			return
		}

		ctx.Set("service", service)
		ctx.Next()
	}
}
