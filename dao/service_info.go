package dao

import (
	"github.com/gin_scaffiold/common/lib"
	"github.com/gin_scaffiold/dto"
	"sync"
	"time"
)

type ServiceInfo struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"负载类型 0=http 1=tcp 2=grpc"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	UpdatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"更新时间"`
	CreatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"添加时间"`
	IsDelete    int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

var (
	wg sync.WaitGroup
)

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (m *ServiceInfo) PageList(input *dto.ServiceListInput) ([]*ServiceInfo, int64) {
	var (
		total int64                     //总条数
		list  = make([]*ServiceInfo, 0) //列表
	)

	query := lib.DBMySQL.Model(m).Debug().Where("is_delete = 0")

	if input.Info != "" {
		query.Where("(service_name LIKE ? or service_desc LIKE ?)", "%"+input.Info+"%", "%"+input.Info+"%")
	}

	query.Count(&total)
	query.Limit(input.PageSize).Offset((input.PageNo - 1) * input.PageSize).Order("id DESC").Find(&list)

	return list, total
}

func (m *ServiceInfo) FindByWhere() *ServiceInfo {

	lib.DBMySQL.Where(m).First(m)

	return m
}

func (m *ServiceInfo) Find(search *ServiceInfo) *ServiceInfo {

	lib.DBMySQL.Where(search).Find(m)

	return m
}

func (m *ServiceInfo) Detail(search *ServiceInfo) *ServiceDetail {
	httpRule := &HttpRule{
		ServiceID: search.ID,
	}

	tcpRule := &TcpRule{
		ServiceID: search.ID,
	}

	grpcRule := &GrpcRule{
		ServiceID: search.ID,
	}

	accessControl := &AccessControl{
		ServiceID: search.ID,
	}

	loadBalance := &LoadBalance{
		ServiceID: search.ID,
	}

	wg.Add(5)

	go func() {
		defer wg.Done()
		lib.DBMySQL /*.Model(new(HttpRule))*/ .Find(httpRule)
	}()

	go func() {
		defer wg.Done()
		lib.DBMySQL /*.Model(new(HttpRule))*/ .Find(tcpRule)
	}()

	go func() {
		defer wg.Done()
		lib.DBMySQL /*.Model(new(HttpRule))*/ .Find(grpcRule)
	}()

	go func() {
		defer wg.Done()
		lib.DBMySQL /*.Model(new(HttpRule))*/ .Find(accessControl)
	}()

	go func() {
		defer wg.Done()
		lib.DBMySQL /*.Model(new(HttpRule))*/ .Find(loadBalance)
	}()

	wg.Wait()

	return &ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
}

func (m *ServiceInfo) Save() error {
	return lib.DBMySQL.Save(m).Error
}

func (m *ServiceInfo) Update() error {
	return lib.DBMySQL.Debug().Model(m).Updates(*m).Error
}

func (m *ServiceInfo) Delete() error {
	return lib.DBMySQL.Delete(m).Error
}

func (m *ServiceInfo) GroupByLoadType() ([]*dto.DashServiceStatItemOutput, error) {
	list := []*dto.DashServiceStatItemOutput{}

	err := lib.DBMySQL.Model(m).Where("is_delete=0").
		Select("load_type, count(*) as value").
		Group("load_type").
		Scan(&list).Error

	return list, err
}
