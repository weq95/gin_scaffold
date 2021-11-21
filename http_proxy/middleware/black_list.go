package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin_scaffiold/dao"
	"github.com/gin_scaffiold/middleware"
	"github.com/gin_scaffiold/public"
	"strings"
)

// HTTPBlackListMiddleware 黑名单
func HTTPBlackListMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		serverInterface, ok := ctx.Get("service")
		if ok {
			middleware.ResponseError(ctx, 2001, errors.New("service not found"))

			ctx.Abort()
			return
		}

		serviceDetail := serverInterface.(*dao.ServiceDetail)
		whileIpList := make([]string, 0)
		if serviceDetail.AccessControl.WhiteList != "" {
			whileIpList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		blackIpList := make([]string, 0)
		if serviceDetail.AccessControl.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}

		if serviceDetail.AccessControl.OpenAuth == 1 && len(whileIpList) == 0 && len(blackIpList) > 0 {
			if public.InStringSlice(blackIpList, ctx.ClientIP()) {
				middleware.ResponseError(ctx, 3001, errors.New(fmt.Sprintf("%s in black ip list", ctx.ClientIP())))

				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
