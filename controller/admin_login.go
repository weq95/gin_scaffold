package controller

import (
	"encoding/json"
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

// AdminLogin godoc
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @Id /admin_login/login
// @Accept json
// @Produce json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func AdminLoginRegister(group *gin.RouterGroup) {
	system := &AdminController{}

	group.POST("/login", system.AdminLogin)
	group.GET("/logout", system.AdminLogOut)
}

func (a *AdminController) AdminLogin(ctx *gin.Context) {
	params := &dto.AdminLoginInput{}

	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}

	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(ctx, db, params)
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
	_ = sess.Save()

	out := &dto.AdminLoginOutput{
		Token: admin.UserName,
	}

	middleware.ResponseSuccess(ctx, out)
}

// AdminLogin godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [get]
func (a *AdminController) AdminLogOut(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Delete(public.AdminSessionInfoKey)
	_ = sess.Save()

	middleware.ResponseSuccess(ctx, "")
}
