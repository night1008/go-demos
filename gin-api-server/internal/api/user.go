package api

import (
	"gin-api-server/internal/extend/ginx"
	"gin-api-server/internal/service"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	UserSvc *service.UserSvc
}

// @Summary 获取当前用户
// @Description 获取当前用户
// @Produce json
// @Tags 用户
// @Success 200 {object} model.UserDTO
// @Router /current [get]
func (a *UserAPI) Current(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := a.UserSvc.Get(ctx, ginx.GetUserID(c))
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}
	ginx.RespondData(c, http.StatusOK, user)
}

// @Summary 修改当前用户密码
// @Description 修改当前用户密码
// @Accept json
// @Produce json
// @Tags 用户
// @Param body body service.UserUpdatePasswordRequest true "body"
// @Success 200 {object} object "{}"
// @Router /current/password [put]
func (a *UserAPI) UpdateCurrentPassword(c *gin.Context) {
	var req service.UserUpdatePasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ctx := c.Request.Context()
	if err := a.UserSvc.UpdatePassword(ctx, ginx.GetUserID(c), &req); err != nil {
		ginx.RespondErr(c, err)
		return
	} else {
		// 更新密码成功后用户需要重新登录
		session := sessions.Default(c)
		session.Clear()
		if err := session.Save(); err != nil {
			// logger.WithContext(c.Request.Context()).Error().Err(err).Msg("failed to save session")
			ginx.RespondErr(c, err)
			return
		}
	}
	ginx.RespondMsg(c, "update success")
}

// @Summary 修改当前用户信息
// @Description 修改当前用户信息
// @Accept json
// @Produce json
// @Tags 用户
// @Param body body service.UserUpdateProfileRequest true "body"
// @Success 200 {object} model.UserDTO
// @Router /current/profile [put]
func (a *UserAPI) UpdateCurrentProfile(c *gin.Context) {
	var req service.UserUpdateProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ctx := c.Request.Context()
	user, err := a.UserSvc.UpdateProfile(ctx, ginx.GetUserID(c), &req)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}
	ginx.RespondMsg(c, user)
}

// @Summary 获取用户组织列表
// @Description 获取用户组织列表
// @Accept json
// @Produce json
// @Tags 用户
// @Success 200 {object} model.UserOrgDTOList
// @Router /current/orgs [get]
func (a *UserAPI) ListUserOrgs(c *gin.Context) {
	ctx := c.Request.Context()
	list, err := a.UserSvc.ListUserOrgs(ctx, ginx.GetUserID(c))
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}
	ginx.RespondData(c, http.StatusOK, list)
}
