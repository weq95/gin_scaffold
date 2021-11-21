package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin_scaffiold/dao"
	"github.com/gin_scaffiold/middleware"
	"github.com/gin_scaffiold/public"
	"strings"
)

// HTTPJwtAuthTokenMiddleware jwt auth token
func HTTPJwtAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))

			c.Abort()
			return
		}

		serviceDetail := serviceInterface.(*dao.ServiceDetail)

		token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		var appMatched bool

		if token != "" {
			claims, err := public.JwtDecode(token)
			if err != nil {
				middleware.ResponseError(c, 2002, err)
				c.Abort()
				return
			}

			appList := dao.AppManagerHandler.GetAppList()
			for _, appInfo := range appList {
				if appInfo.AppID == claims.Issuer {
					c.Set("app", appInfo)
					appMatched = true

					break
				}
			}
		}

		if serviceDetail.AccessControl.OpenAuth == 1 && !appMatched {
			middleware.ResponseError(c, 2003, errors.New("not match valid app"))

			c.Abort()
			return
		}

		c.Next()
	}
}
