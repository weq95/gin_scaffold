package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin_scaffiold/dao"
	"github.com/gin_scaffiold/dto"
	"github.com/gin_scaffiold/middleware"
	"github.com/gin_scaffiold/public"
)

type DashboardController struct {
}

func DashboardRegister(group *gin.RouterGroup) {
	dial := &DashboardController{}

	group.GET("/panel_group_data", dial.PanelGroupData)
	group.GET("/flow_stat", dial.FlowStat)
	group.GET("/service_stat", dial.ServiceStat)
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 首页大盘
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [get]
func (c DashboardController) PanelGroupData(ctx *gin.Context) {
	serviceInfo := &dao.ServiceInfo{}
	_, total := serviceInfo.PageList(&dto.ServiceListInput{
		PageNo:   1,
		PageSize: 1,
	})

	if total <= 0 {
		middleware.ResponseError(ctx, 2001, errors.New("未获取服务数据"))
		return
	}

	app := dao.App{}
	_, appNum, err := app.AppList(&dto.APPListInput{
		PageNo:   1,
		PageSize: 1,
	})
	if err != nil {
		middleware.ResponseError(ctx, 2002, err)
		return
	}

	counter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
	if err != nil {
		middleware.ResponseError(ctx, 2003, err)
		return
	}

	middleware.ResponseSuccess(ctx, &dto.PanelGroupDataOutput{
		ServiceNum:      total,
		AppNum:          appNum,
		CurrentQPS:      counter.TotalCount,
		TodayRequestNum: counter.QPS,
	})
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /dashboard/service_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.DashServiceStatOutput} "success"
// @Router /dashboard/service_stat [get]
func (c DashboardController) FlowStat(ctx *gin.Context) {
	serviceInfo := &dao.ServiceInfo{}
	list, err := serviceInfo.GroupByLoadType()
	if err != nil {
		middleware.ResponseError(ctx, 2002, err)
		return
	}

	legend := []string{}
	for i, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			middleware.ResponseError(ctx, 2003, errors.New("load_type not found"))
			return
		}

		list[i].Name = name
		legend = append(legend, name)
	}

	middleware.ResponseSuccess(ctx, &dto.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	})
}

// FlowStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /dashboard/flow_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /dashboard/flow_stat [get]
func (c DashboardController) ServiceStat(ctx *gin.Context) {

}
