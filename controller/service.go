package controller

import "github.com/gin-gonic/gin"

type ServiceController struct {
}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}

	group.GET("/service_list", service.List)
	group.GET("/service_delete", service.Delete)
	group.GET("/service_detail", service.Detail)
	group.GET("/service_stat", service.Stat)

	group.POST("/service_add_http", service.AddHTTP)
	group.POST("/service_update_http", service.UpdateHTTP)

	group.POST("/service_add_tcp", service.AddTcp)
	group.POST("/service_update_tcp", service.UpdateTcp)

	group.POST("/service_add_grpc", service.AddGrpc)
	group.POST("/service_update_grpc", service.UpdateGrpc)
}

func (s ServiceController) List(ctx *gin.Context) {

}

func (s ServiceController) Delete(ctx *gin.Context) {

}

func (s ServiceController) Detail(ctx *gin.Context) {

}

func (s ServiceController) Stat(ctx *gin.Context) {

}

func (s ServiceController) AddHTTP(ctx *gin.Context) {

}

func (s ServiceController) UpdateHTTP(ctx *gin.Context) {

}

func (s ServiceController) AddTcp(ctx *gin.Context) {

}

func (s ServiceController) UpdateTcp(ctx *gin.Context) {

}

func (s ServiceController) AddGrpc(ctx *gin.Context) {

}

func (s ServiceController) UpdateGrpc(ctx *gin.Context) {

}
