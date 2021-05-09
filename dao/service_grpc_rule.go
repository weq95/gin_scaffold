package dao

import "github.com/gin_scaffiold/common/lib"

type GrpcRule struct {
	ID             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	Port           int    `json:"port" gorm:"column:port" description:"端口	"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue"`
}

func (t *GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}

func (m *GrpcRule) Find(search *GrpcRule) *GrpcRule {

	lib.DBMySQL.Where(search).Find(m)

	return m
}

func (m *GrpcRule) Save() error {
	return lib.DBMySQL.Save(m).Error
}
