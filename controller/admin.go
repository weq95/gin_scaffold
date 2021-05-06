package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin_scaffiold/common/lib"
	"github.com/gin_scaffiold/dao"
	"github.com/gin_scaffiold/dto"
	"github.com/gin_scaffiold/middleware"
	"github.com/gin_scaffiold/public"
	"time"
)

type AdminController struct {
}


func AdminLoginRegister(group *gin.RouterGroup) {
	admin := &AdminController{}
	//不需要登录的
	group.POST("/login", admin.Login)
	group.GET("/logout", admin.LogOut)
}

func AdminRegister(group *gin.RouterGroup) {
	admin := &AdminController{}
	group.GET("/admin_info", admin.Info)
	group.POST("/change_pwd", admin.ChangPwd)
}

// Login godoc
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @Id /admins/login
// @Accept json
// @Produce json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admins/login [post]
func (a *AdminController) Login(ctx *gin.Context) {
	params := &dto.AdminLoginInput{}

	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	admin := &dao.Admin{
		UserName: params.UserName,
	}
	admin, err := admin.LoginCheck(params)
	if err != nil {
		middleware.ResponseError(ctx, 2002, err)
		return
	}

	//设置session
	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}

	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(ctx, 2003, err)
		return
	}

	sess := sessions.Default(ctx)
	sess.Set(public.AdminSessionInfoKey, string(sessBts))
	err = sess.Save()
	if err != nil {
		middleware.ResponseError(ctx, 2004, err)
		return
	}

	out := &dto.AdminLoginOutput{
		Token: admin.UserName,
	}

	middleware.ResponseSuccess(ctx, out)
}

// LogOut godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admins/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admins/logout [get]
func (a *AdminController) LogOut(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Delete(public.AdminSessionInfoKey)
	_ = sess.Save()

	middleware.ResponseSuccess(ctx, "")
}

// Info godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (a *AdminController) Info(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sessInof := sess.Get(public.AdminSessionInfoKey)

	adminSessionInfo := &dto.AdminSessionInfo{}
	err := json.Unmarshal([]byte(fmt.Sprint(sessInof)), adminSessionInfo)
	if err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	//1. 读取sessionKey对应json 转换为结构体
	//2. 取出数据然后封装输出结构体

	//Avatar       string    `json:"avatar"`
	//Introduction string    `json:"introduction"`
	//Roles        []string  `json:"roles"`

	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}

	middleware.ResponseSuccess(ctx, out)
}

// ChangePwd godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (a *AdminController) ChangPwd(ctx *gin.Context) {
	params := &dto.ChangePwdInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	//1. session读取用户信息到结构体 sessInfo
	//2. sessInfo.ID 读取数据库信息 adminInfo
	//3. params.password+adminInfo.salt sha256 saltPassword
	//4. saltPassword==> adminInfo.password 执行数据保存

	//session读取用户信息到结构体
	sess := sessions.Default(ctx)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSession := &dto.AdminSessionInfo{}
	err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSession)
	if err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	//数据库获取用户信息
	admin := &dao.Admin{
		UserName: adminSession.UserName,
	}

	lib.DBMySQL.Debug().Where("id = ?", adminSession.ID).Where("user_name = ? ", adminSession.UserName).First(admin)
	if admin.Id <= 0 {
		middleware.ResponseError(ctx, 2002, err)
		return
	}

	//生成新密码
	admin.Password = public.GenSaltPassword(admin.Salt, params.Password)
	lib.DBMySQL.Model(admin).Update("password", "update_at")

	middleware.ResponseSuccess(ctx, "")
}
