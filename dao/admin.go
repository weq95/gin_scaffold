package dao

import (
	"errors"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/gin_scaffiold/common/lib"
	"github.com/gin_scaffiold/dto"
	"github.com/gin_scaffiold/public"
	"time"
)

type Admin struct {
	Id        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName  string    `json:"user_name" gorm:"column:user_name" description:"管理员用户名"`
	Salt      string    `json:"salt" gorm:"column:salt" description:"盐"`
	Password  string    `json:"password" gorm:"column:password" description:"密码"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (a *Admin) TableName() string {
	return "gateway_admin"
}

func (a *Admin) LoginCheck(input *dto.AdminLoginInput) (*Admin, error) {
	lib.DBMySQL.Debug().Where("user_name = ?", a.UserName).First(a)


	if a.Id <= 0 {
		return nil, errors.New("用户信息不存在")
	}

	saltPwd := public.GenSaltPassword(a.Salt, input.Password)
	if a.Password != saltPwd {
		return nil, errors.New("用户名或密码错误")
	}

	return a, nil
}

func (a *Admin) Save(ctx *gin.Context, db *gorm.DB) error {
	return db.SetCtx(public.GetGinTraceContext(ctx)).Save(a).Error
}
