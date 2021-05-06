package dao

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin_scaffiold/public"
	"strings"
	"sync"
)

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"基本信息"`
	HTTPRule      *HttpRule      `json:"http_rule" description:"http_rule"`
	TCPRule       *TcpRule       `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *GrpcRule      `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance" description:"load_balance"`
	AccessControl *AccessControl `json:"access_control" description:"access_control"`
}

var ServiceManagerHandler *ServiceManger

func init() {
	ServiceManagerHandler = NewServiceManger()
}

type ServiceManger struct {
	ServiceMap   map[string]*ServiceDetail
	ServiceSlice []*ServiceDetail
	Locker       sync.RWMutex
	init         sync.Once
	err          error
}

func NewServiceManger() *ServiceManger {
	return &ServiceManger{
		ServiceMap:   map[string]*ServiceDetail{},
		ServiceSlice: []*ServiceDetail{},
		Locker:       sync.RWMutex{},
		init:         sync.Once{},
	}
}

func (s *ServiceManger) GetTcpServiceList() []*ServiceDetail {
	var list []*ServiceDetail

	for _, item := range s.ServiceSlice {
		temp := item
		if temp.Info.LoadType == public.LoadTypeGRPC {
			list = append(list, temp)
		}
	}

	return list
}

func (s *ServiceManger) GetGrpcServiceList() []*ServiceDetail {
	var list []*ServiceDetail

	for _, serviceItem := range s.ServiceSlice {
		temp := serviceItem
		if temp.Info.LoadType == public.LoadTypeGRPC {
			list = append(list, temp)
		}
	}

	return list
}

func (s *ServiceManger) HTTPAccessMode(ctx *gin.Context) (*ServiceDetail, error) {
	//1、前缀匹配 /abc ==> serviceSlice.rule
	//2、域名匹配 www.test.com ==> serviceSlice.rule
	//host c.Request.Host
	//path c.Request.URL.Path

	host := ctx.Request.Host
	host = host[0:strings.Index(host, ":")]
	path := ctx.Request.URL.Path

	for _, serviceItem := range s.ServiceSlice {
		if serviceItem.Info.LoadType != public.LoadTypeHTTP {
			continue
		}

		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule == host {
				return serviceItem, nil
			}
		}

		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				return serviceItem, nil
			}
		}
	}

	return nil, errors.New("not matched service")
}

func (s *ServiceManger) LoadOnce() error {
	s.init.Do(func() {
		/*serviceInfo := &dao.ServiceInfo{}
		c,_ := gin.CreateTestContext(httptest.NewRecorder())
		*/

	})

	return s.err
}
