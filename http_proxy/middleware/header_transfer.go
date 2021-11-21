package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin_scaffiold/dao"
	"github.com/gin_scaffiold/middleware"
	"strings"
)

// HTTPHeaderTransferMiddleware 基于请求信息,配置接入方式
func HTTPHeaderTransferMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		serviceInterface, ok := ctx.Get("service")
		if !ok {
			middleware.ResponseError(ctx, 2001, errors.New("service not found"))

			ctx.Abort()
			return
		}

		serviceDetail := serviceInterface.(dao.ServiceDetail)
		for _, item := range strings.Split(serviceDetail.HTTPRule.HeaderTransfor, ",") {
			items := strings.Split(item, ",")
			if len(items) != 3 {
				continue
			}

			if items[0] == "add" || items[0] == "edit" {
				ctx.Request.Header.Set(items[1], items[2])
			}

			if items[0] == "del" {
				ctx.Request.Header.Del(items[1])
			}
		}

		ctx.Next()
	}
}
