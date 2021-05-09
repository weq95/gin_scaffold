package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin_scaffiold/common/lib"
	"github.com/gin_scaffiold/dao"
	"github.com/gin_scaffiold/dto"
	"github.com/gin_scaffiold/middleware"
	"github.com/gin_scaffiold/public"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type ServiceController struct {
}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}

	group.GET("/service_list", service.List)
	group.DELETE("/service_delete", service.Delete)
	group.GET("/service_detail", service.Detail)
	group.GET("/service_stat", service.Stat)

	group.POST("/service_add_http", service.AddHTTP)
	group.POST("/service_update_http", service.UpdateHTTP)

	group.POST("/service_add_tcp", service.AddTcp)
	group.POST("/service_update_tcp", service.UpdateTcp)

	group.POST("/service_add_grpc", service.AddGrpc)
	group.POST("/service_update_grpc", service.UpdateGrpc)
}

// List godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query int true "每页条数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (s *ServiceController) List(ctx *gin.Context) {
	params := &dto.ServiceListInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	outList := make([]*dto.ServiceListItemOutput, 0)
	list, total := new(dao.ServiceInfo).PageList(params)

	for _, info := range list {

		detail := info.Detail(info)

		clusterIp := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")

		serviceAddr := "unknow"
		if detail.Info.LoadType == public.LoadTypeHTTP {
			if detail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
				//http
				serviceAddr = clusterIp + ":" + clusterPort + detail.HTTPRule.Rule

				//支持https
				if detail.HTTPRule.NeedHttps == 1 {
					serviceAddr = clusterIp + ":" + clusterSSLPort + detail.HTTPRule.Rule
				}
			}

			if detail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
				serviceAddr = detail.HTTPRule.Rule
			}
		}

		if detail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = clusterIp + ":" + strconv.Itoa(detail.TCPRule.Port)
		}

		if detail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = clusterIp + ":" + strconv.Itoa(detail.GRPCRule.Port)
		}

		outList = append(outList, &dto.ServiceListItemOutput{
			ID:          info.ID,
			ServiceName: info.ServiceName,
			ServiceDesc: info.ServiceDesc,
			LoadType:    info.LoadType,
			ServiceAddr: serviceAddr,
			Qps:         0,
			Qpd:         0,
			TotalNode:   len(detail.LoadBalance.GetIPListModel()), //节点数
		})
	}

	middleware.ResponseSuccess(ctx, &dto.ServiceListOutput{
		Total: total,
		List:  outList,
	})
}

// Delete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/service_delete
// @Accept  json
// @Produce  json
// @Param id query int true "服务id"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_delete [delete]
func (s *ServiceController) Delete(ctx *gin.Context) {
	params := &dto.ServiceDeleteInput{}

	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{
		ID:       params.ID,
		IsDelete: 1,
	}

	_ = serviceInfo.Update()

	middleware.ResponseSuccess(ctx, "操作成功")
}

// ServiceDetail godoc
// @Summary 服务详情
// @Description 服务详情
// @Tags 服务管理
// @ID /service/service_detail
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=dao.ServiceDetail} "success"
// @Router /service/service_detail [get]
func (s *ServiceController) Detail(ctx *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	//基本信息
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo = serviceInfo.Find(serviceInfo)
	if serviceInfo.ID <= 0 {
		middleware.ResponseError(ctx, 2003, errors.New("未获取到详情信息"))
		return
	}

	detail := serviceInfo.Detail(serviceInfo)
	middleware.ResponseSuccess(ctx, detail)

}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 服务管理
// @ID /service/service_stat
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /service/service_stat [get]
func (s *ServiceController) Stat(ctx *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	//基本信息
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	detail := serviceInfo.Detail(serviceInfo)
	if serviceInfo.ID <= 0 {
		middleware.ResponseError(ctx, 2004, errors.New("未获取到详情"))
		return
	}

	counter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + detail.Info.ServiceName)
	if err != nil {
		middleware.ResponseError(ctx, 2004, err)
		return
	}

	//今天的数据详情
	todayList := []int64{}
	nowTime := time.Now()
	for i := 0; i < nowTime.Hour(); i++ {
		dateTime := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), i, 0, 0, 0, lib.TimeLocaltion)
		hourDate, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, hourDate)
	}

	//昨天的数据详情
	yesterdayList := []int64{}
	yesterTime := nowTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i < 23; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, lib.TimeLocaltion)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, hourData)
	}

	middleware.ResponseSuccess(ctx, &dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}

// ServiceAddHTTP godoc
// @Summary 添加HTTP服务
// @Description 添加HTTP服务
// @Tags 服务管理
// @ID /service/service_add_http
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_http [post]
func (s *ServiceController) AddHTTP(ctx *gin.Context) {
	params := &dto.ServiceAddHTTPInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	//检测服务名称
	serviceInfo := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
	}
	serviceInfo.FindByWhere()
	if serviceInfo.ID > 0 {
		middleware.ResponseError(ctx, 2002, errors.New("服务名称已存在"))
		return
	}

	//检测接入域名或前缀
	httpUrl := &dao.HttpRule{
		RuleType: params.RuleType,
		Rule:     params.Rule,
	}
	httpUrl.FindByWhere()
	if httpUrl.ID > 0 {
		middleware.ResponseError(ctx, 2003, errors.New("服务接入前缀或域名已存在"))
		return
	}

	//开启事务,更新信息
	_ = lib.DBMySQL.Transaction(func(tx *gorm.DB) error {
		serviceModel := &dao.ServiceInfo{
			ServiceName: params.ServiceName,
			ServiceDesc: params.ServiceDesc,
		}
		err := serviceModel.Save()
		if err != nil {
			return err
		}

		httpRule := &dao.HttpRule{
			ServiceID:      serviceModel.ID,
			RuleType:       params.RuleType,
			Rule:           params.Rule,
			NeedHttps:      params.NeedHttps,
			NeedStripUri:   params.NeedStripUri,
			NeedWebsocket:  params.NeedWebsocket,
			UrlRewrite:     params.UrlRewrite,
			HeaderTransfor: params.HeaderTransfor,
		}
		err = httpRule.Save()
		if err != nil {
			return err
		}

		accessControl := &dao.AccessControl{
			ServiceID:         serviceModel.ID,
			OpenAuth:          params.OpenAuth,
			BlackList:         params.BlackList,
			WhiteList:         params.WhiteList,
			ClientIPFlowLimit: params.ClientipFlowLimit,
			ServiceFlowLimit:  params.ServiceFlowLimit,
		}
		err = accessControl.Save()
		if err != nil {
			return err
		}

		loadbalance := &dao.LoadBalance{
			ServiceID:              serviceModel.ID,
			RoundType:              params.RoundType,
			IpList:                 params.IpList,
			WeightList:             params.WeightList,
			UpstreamConnectTimeout: params.UpstreamConnectTimeout,
			UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
			UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
			UpstreamMaxIdle:        params.UpstreamMaxIdle,
		}
		err = loadbalance.Save()
		if err != nil {
			return err
		}

		return nil
	})

	middleware.ResponseSuccess(ctx, "")
}

// UpdateHTTP godoc
// @Summary 修改HTTP服务
// @Description 修改HTTP服务
// @Tags 服务管理
// @ID /service/service_update_http
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_http [post]
func (s *ServiceController) UpdateHTTP(ctx *gin.Context) {
	params := &dto.ServiceUpdateHTTPInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(ctx, 2001, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	serviceInfo.Find(serviceInfo)
	if serviceInfo.ID <= 0 {
		middleware.ResponseError(ctx, 2003, errors.New("服务不存在"))
		return
	}

	detail := serviceInfo.Detail(serviceInfo)
	if serviceInfo.ID <= 0 {
		middleware.ResponseError(ctx, 2004, errors.New("服务详情异常"))
		return
	}

	_ = lib.DBMySQL.Transaction(func(tx *gorm.DB) error {
		//更新服务信息
		info := detail.Info
		info.ServiceDesc = params.ServiceDesc
		err := info.Update()
		if err != nil {
			middleware.ResponseError(ctx, 2005, err)
			return err
		}

		//更新http服务规则
		httpRule := detail.HTTPRule
		httpRule.NeedHttps = params.NeedHttps
		httpRule.NeedStripUri = params.NeedStripUri
		httpRule.NeedWebsocket = params.NeedWebsocket
		httpRule.UrlRewrite = params.UrlRewrite
		httpRule.HeaderTransfor = params.HeaderTransfor
		err = httpRule.Update()
		if err != nil {
			middleware.ResponseError(ctx, 2006, err)
			return err
		}

		//更新权限
		accessControl := detail.AccessControl
		accessControl.OpenAuth = params.OpenAuth
		accessControl.BlackList = params.BlackList
		accessControl.WhiteList = params.WhiteList
		accessControl.ClientIPFlowLimit = params.ClientipFlowLimit
		accessControl.ServiceFlowLimit = params.ServiceFlowLimit
		err = accessControl.Update()
		if err != nil {
			middleware.ResponseError(ctx, 2007, err)
			return err
		}

		//更新负载均衡规则
		loadbalance := detail.LoadBalance
		loadbalance.RoundType = params.RoundType
		loadbalance.IpList = params.IpList
		loadbalance.WeightList = params.WeightList
		loadbalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
		loadbalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
		loadbalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
		loadbalance.UpstreamMaxIdle = params.UpstreamMaxIdle
		err = loadbalance.Update()
		if err != nil {
			middleware.ResponseError(ctx, 2008, err)
			return err
		}

		return nil
	})

	middleware.ResponseSuccess(ctx, "")
}

// AddTcp godoc
// @Summary tcp服务添加
// @Description tcp服务添加
// @Tags 服务管理
// @ID /service/service_add_tcp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddTcpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_tcp [post]
func (s *ServiceController) AddTcp(ctx *gin.Context) {
	params := &dto.ServiceAddTcpInput{}
	if err := params.GetValidParams(ctx); err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}

	//验证ip与权限数量是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(ctx, 2005, errors.New("ip列表与权重设置不匹配"))
		return
	}

	//验证服务名称是否已存在
	infoSearch := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	infoSearch.Find(infoSearch)
	if infoSearch.ID <= 0 {
		middleware.ResponseError(ctx, 2002, errors.New("服务名称已存在"))
		return
	}

	//验证端口是否被占用
	tcpRuleSearch := &dao.TcpRule{
		Port: params.Port,
	}
	tcpRuleSearch.Find(tcpRuleSearch)
	if tcpRuleSearch.ID <= 0 {
		middleware.ResponseError(ctx, 2004, errors.New("服务端口已被占用"))
		return
	}

	err := lib.DBMySQL.Transaction(func(tx *gorm.DB) error {
		//添加TCP服务
		info := &dao.ServiceInfo{
			LoadType:    public.LoadTypeTCP,
			ServiceName: params.ServiceName,
			ServiceDesc: params.ServiceDesc,
		}
		err := info.Save()
		if err != nil {
			return err
		}

		//提娜佳负载均衡规则
		loadBalance := &dao.LoadBalance{
			ServiceID:  info.ID,
			RoundType:  params.RoundType,
			IpList:     params.IpList,
			WeightList: params.WeightList,
			ForbidList: params.ForbidList,
		}
		err = loadBalance.Save()
		if err != nil {
			return err
		}

		//添加TCP规则
		httpRule := &dao.TcpRule{
			ServiceID: info.ID,
			Port:      params.Port,
		}
		err = httpRule.Save()
		if err != nil {
			return err
		}

		//添加权限
		accessControl := &dao.AccessControl{
			ServiceID:         info.ID,
			OpenAuth:          params.OpenAuth,
			BlackList:         params.BlackList,
			WhiteList:         params.WhiteList,
			WhiteHostName:     params.WhiteHostName,
			ClientIPFlowLimit: params.ClientIPFlowLimit,
			ServiceFlowLimit:  params.ServiceFlowLimit,
		}
		err = accessControl.Save()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	middleware.ResponseSuccess(ctx, "")
}

// UpdateTcp godoc
// @Summary tcp服务更新
// @Description tcp服务更新
// @Tags 服务管理
// @ID /service/service_update_tcp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateTcpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_tcp [post]
func (s *ServiceController) UpdateTcp(ctx *gin.Context) {
	params := &dto.ServiceUpdateTcpInput{}
	if err := params.GetValidParams(ctx); err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(ctx, 2002, errors.New("ip列表与权重列表不一致"))
		return
	}

	service := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail := service.Detail(service)

	err := lib.DBMySQL.Transaction(func(tx *gorm.DB) error {
		info := detail.Info
		info.ServiceDesc = params.ServiceDesc
		err := info.Update()
		if err != nil {
			return err
		}

		loadBalance := &dao.LoadBalance{}
		if detail.LoadBalance != nil {
			loadBalance = detail.LoadBalance
		}
		loadBalance.ServiceID = info.ID
		loadBalance.RoundType = params.RoundType
		loadBalance.IpList = params.IpList
		loadBalance.WeightList = params.WeightList
		loadBalance.ForbidList = params.ForbidList
		err = loadBalance.Save()
		if err != nil {
			return err
		}

		tcpRule := &dao.TcpRule{}
		if detail.TCPRule != nil {
			tcpRule = detail.TCPRule
		}
		tcpRule.ServiceID = info.ID
		tcpRule.Port = params.Port
		err = tcpRule.Save()
		if err != nil {
			return err
		}

		accessControl := &dao.AccessControl{}
		if detail.AccessControl != nil {
			accessControl = detail.AccessControl
		}
		accessControl.ServiceID = info.ID
		accessControl.OpenAuth = params.OpenAuth
		accessControl.BlackList = params.BlackList
		accessControl.WhiteList = params.WhiteList
		accessControl.WhiteHostName = params.WhiteHostName
		accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
		accessControl.ServiceFlowLimit = params.ServiceFlowLimit
		err = accessControl.Save()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	middleware.ResponseSuccess(ctx, "")
}

// AddGrpc godoc
// @Summary grpc服务添加
// @Description grpc服务添加
// @Tags 服务管理
// @ID /service/service_add_grpc
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddGrpcInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_grpc [post]
func (s *ServiceController) AddGrpc(ctx *gin.Context) {
	params := &dto.ServiceAddGrpcInput{}
	if err := params.GetValidParams(ctx); err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(ctx, 2005, errors.New("ip列表与权重设置不匹配"))
		return
	}

	//验证服务名称是否已被占用
	infoSearch := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	infoSearch.Find(infoSearch)
	if infoSearch.ID > 0 {
		middleware.ResponseError(ctx, 2002, errors.New("服务名称已存在"))
		return
	}

	//验证端口是否已被占用
	tcpRuleSearch := &dao.TcpRule{
		Port: params.Port,
	}
	tcpRuleSearch.Find(tcpRuleSearch)
	if tcpRuleSearch.ID > 0 {
		middleware.ResponseError(ctx, 2003, errors.New("服务端口已被占用"))
		return
	}

	grpcRuleSearch := &dao.GrpcRule{
		Port: params.Port,
	}
	grpcRuleSearch.Find(grpcRuleSearch)
	if grpcRuleSearch.ID > 0 {
		middleware.ResponseError(ctx, 2004, errors.New("服务端口已被占用"))
		return
	}

	err := lib.DBMySQL.Transaction(func(tx *gorm.DB) error {
		info := dao.ServiceInfo{
			LoadType:    public.LoadTypeGRPC,
			ServiceName: params.ServiceName,
			ServiceDesc: params.ServiceDesc,
		}
		err := info.Save()
		if err != nil {
			return err
		}

		loadBalance := &dao.LoadBalance{
			ServiceID:  info.ID,
			RoundType:  params.RoundType,
			IpList:     params.IpList,
			WeightList: params.WeightList,
			ForbidList: params.ForbidList,
		}
		err = loadBalance.Save()
		if err != nil {
			return err
		}

		grpcRule := dao.GrpcRule{
			ServiceID:      info.ID,
			Port:           params.Port,
			HeaderTransfor: params.HeaderTransfor,
		}
		err = grpcRule.Save()
		if err != nil {
			return err
		}

		accessControl := &dao.AccessControl{
			ServiceID:         info.ID,
			OpenAuth:          params.OpenAuth,
			BlackList:         params.BlackList,
			WhiteList:         params.WhiteList,
			WhiteHostName:     params.WhiteHostName,
			ClientIPFlowLimit: params.ClientIPFlowLimit,
			ServiceFlowLimit:  params.ServiceFlowLimit,
		}
		err = accessControl.Save()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		middleware.ResponseError(ctx, 2009, err)
		return
	}

	middleware.ResponseSuccess(ctx, "")
}

// UpdateGrpc godoc
// @Summary grpc服务更新
// @Description grpc服务更新
// @Tags 服务管理
// @ID /service/service_update_grpc
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateGrpcInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_grpc [post]
func (s *ServiceController) UpdateGrpc(ctx *gin.Context) {
	params := &dto.ServiceUpdateGrpcInput{}
	if err := params.GetValidParams(ctx); err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(ctx, 2002, errors.New("ip列表与权重设置不匹配"))
		return
	}

	service := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail := service.Detail(service)

	err := lib.DBMySQL.Transaction(func(tx *gorm.DB) error {
		info := detail.Info
		info.ServiceDesc = params.ServiceDesc
		err := info.Save()
		if err != nil {
			return err
		}

		loadBalance := &dao.LoadBalance{}
		if detail.LoadBalance != nil {
			loadBalance = detail.LoadBalance
		}
		loadBalance.ServiceID = info.ID
		loadBalance.RoundType = params.RoundType
		loadBalance.IpList = params.IpList
		loadBalance.WeightList = params.WeightList
		loadBalance.ForbidList = params.ForbidList
		err = loadBalance.Save()
		if err != nil {
			return err
		}

		grpcRule := &dao.GrpcRule{}
		if detail.GRPCRule != nil {
			grpcRule = detail.GRPCRule
		}
		grpcRule.ServiceID = info.ID
		grpcRule.HeaderTransfor = params.HeaderTransfor
		err = grpcRule.Save()
		if err != nil {
			return err
		}

		accessControl := &dao.AccessControl{}
		if detail.AccessControl != nil {
			accessControl = detail.AccessControl
		}
		accessControl.ServiceID = info.ID
		accessControl.OpenAuth = params.OpenAuth
		accessControl.BlackList = params.BlackList
		accessControl.WhiteList = params.WhiteList
		accessControl.WhiteHostName = params.WhiteHostName
		accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
		accessControl.ServiceFlowLimit = params.ServiceFlowLimit
		err = accessControl.Save()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		middleware.ResponseError(ctx, 2008, err)
		return
	}

	middleware.ResponseSuccess(ctx, "")
}
