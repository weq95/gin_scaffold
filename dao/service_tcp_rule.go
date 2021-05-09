package dao

import "github.com/gin_scaffiold/common/lib"

type TcpRule struct {
	ID        int64 `json:"id" gorm:"primary_key"`
	ServiceID int64 `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	Port      int   `json:"port" gorm:"column:port" description:"端口	"`
}

func (t *TcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}

func (m *TcpRule) Find(search *TcpRule) *TcpRule {

	lib.DBMySQL.Where(search).Find(m)

	return m
}

func (m *TcpRule) Save() error {
	return lib.DBMySQL.Save(m).Error
}