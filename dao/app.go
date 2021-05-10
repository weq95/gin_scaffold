package dao

import (
	"github.com/gin_scaffiold/common/lib"
	"github.com/gin_scaffiold/dto"
	"sync"
	"time"
)

type App struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	AppID     string    `json:"app_id" gorm:"column:app_id" description:"租户id	"`
	Name      string    `json:"name" gorm:"column:name" description:"租户名称	"`
	Secret    string    `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS  string    `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配"`
	Qpd       int64     `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	Qps       int64     `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间	"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

func (m *App) TableName() string {
	return "gateway_app"
}

func (m *App) First(search *App) error {
	return lib.DBMySQL.Model(search).First(search).Error
}

func (m *App) Save() error {
	return lib.DBMySQL.Model(m).Save(m).Error
}

func (m *App) AppList(params *dto.APPListInput) ([]*App, int64, error) {
	var (
		list   []*App
		count  int64
		offset = (params.PageNo - 1) * params.PageSize
	)

	query := lib.DBMySQL.Model(m).Where("is_delete = 0")
	if len(params.Info) > 0 {
		query = query.Where(" (name LIKE ? OR app_id LIKE ?)", "%"+params.Info+"%"+"%"+params.Info+"%")
	}

	err := query.Limit(params.PageSize).Offset(offset).Order("id DESC").Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

type AppManager struct {
	AppMap   map[string]*App
	AppSlice []*App
	Locker   sync.RWMutex
	init     sync.Once
	err      error
}

var AppManagerHandler *AppManager

func init() {
	AppManagerHandler = NewAppManager()
}

func NewAppManager() *AppManager {
	return &AppManager{
		AppMap:   map[string]*App{},
		AppSlice: []*App{},
		Locker:   sync.RWMutex{},
		init:     sync.Once{},
	}
}

func (m *AppManager) GetAppList() []*App {
	return m.AppSlice
}

func (m *AppManager) LoadOnce() error {
	m.init.Do(func() {
		appInfo := &App{}
		list, _, err := appInfo.AppList(&dto.APPListInput{
			PageNo:   1,
			PageSize: 9999,
		})
		if err != nil {
			m.err = err
			return
		}

		m.Locker.Lock()
		defer m.Locker.Unlock()
		for _, app := range list {
			m.AppMap[app.AppID] = app
			m.AppSlice = append(m.AppSlice, app)
		}
	})

	return m.err
}
