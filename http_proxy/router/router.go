package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin_scaffiold/middleware"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.New()
	router.Use(middlewares...)
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	oauth := router.Group("/oauth")
	oauth.Use(middleware.TranslationMiddleware())
	{
		controller.OAuthRegister(oauth)
	}

	router.Use(
		router.HTTPAccessModeMiddleware(),
		router.HTTPFlowCountMiddleware(),
		router.HTTPFlowLimitMiddleware(),
		router.HTTPJwtAuthTokenMiddleware(),
		router.HTTPJwtFlowCountMiddleware(),
		router.HTTPJwtFlowLimitMiddleware(),
		router.HTTPWhiteListMiddleware(),
		router.HTTPBlackListMiddleware(),
		router.HTTPHeaderTransferMiddleware(),
		router.HTTPStripUriMiddleware(),
		router.HTTPUrlRewriteMiddleware(),
		router.HTTPReverseProxyMiddleware())

	return router
}
