package api

import (
	"gin-api-server/internal/extend/ginx"
	"gin-api-server/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrgAPI struct {
	OrgSvc *service.OrgSvc
}

// @Summary 创建组织
// @Description 创建组织
// @Accept json
// @Produce json
// @Tags 组织
// @Param body body service.OrgCreateRequest true "body"
// @Success 200 {object} model.OrgDTO
// @Router /orgs [post]
func (a *OrgAPI) Create(c *gin.Context) {
	var req service.OrgCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ctx := c.Request.Context()
	org, err := a.OrgSvc.Create(ctx, ginx.GetUserID(c), &req)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}
	ginx.RespondData(c, http.StatusCreated, org)
}

// @Summary 更新组织
// @Description 更新组织
// @Accept json
// @Produce json
// @Tags 组织
// @Param id path integer true "org id"
// @Param body body service.OrgUpdateRequest true "body"
// @Success 200 {object} model.OrgDTO
// @Router /orgs/{id} [put]
func (a *OrgAPI) Update(c *gin.Context) {
	var req service.OrgUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ctx := c.Request.Context()
	req.ID = ginx.ParseParamID(c, "id")
	org, err := a.OrgSvc.Update(ctx, ginx.GetUserID(c), &req)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, org)
}

// @Summary 删除组织
// @Description 删除组织
// @Accept json
// @Produce json
// @Tags 组织
// @Param id path integer true "org id"
// @Success 204 {object} object "{}"
// @Router /orgs/{id} [delete]
func (a *OrgAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	if err := a.OrgSvc.Delete(ctx, ginx.GetUserID(c), ginx.ParseParamID(c, "id")); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusNoContent, gin.H{})
}

// @Summary 获取拥有者列表
// @Description 获取拥有者列表
// @Produce json
// @Tags 组织
// @Param id path integer true "org id"
// @Success 200 {object} model.UserDTOList
// @Router /orgs/{id}/owners [get]
func (a *OrgAPI) ListOwners(c *gin.Context) {
	ctx := c.Request.Context()
	owners, err := a.OrgSvc.ListOwners(ctx, ginx.GetUserID(c), ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, owners)
}

// @Summary 获取组织成员列表
// @Description 获取组织成员列表
// @Produce json
// @Tags 组织
// @Param id path integer true "org id"
// @Param app_id query integer false "app id"
// @Success 200 {object} model.UserDTOList
// @Router /orgs/{id}/members [get]
func (a *OrgAPI) ListAppMembers(c *gin.Context) {
	var req service.OrgListAppUserRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ctx := c.Request.Context()
	req.OrgID = ginx.ParseParamID(c, "id")
	users, err := a.OrgSvc.ListAppMembers(ctx, ginx.GetUserID(c), &req)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, users)
}

// @Summary 获取组织应用列表
// @Description 获取组织应用列表
// @Accept json
// @Produce json
// @Tags 组织
// @Success 200 {object} model.UserOrgDTOList
// @Router /orgs/{id}/apps [get]
func (a *OrgAPI) ListApps(c *gin.Context) {
	ctx := c.Request.Context()
	apps, err := a.OrgSvc.ListApps(ctx, ginx.GetUserID(c), ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, apps)
}

// @Summary 查找用户和应用列表
// @Description 查找用户和应用列表
// @Accept json
// @Produce json
// @Tags 组织
// @Param id path integer true "org id"
// @Param email query string true "user email"
// @Success 200 {object} model.UserWithAppDTO
// @Router /orgs/{id}/users-by-query [get]
func (a *OrgAPI) QueryUser(c *gin.Context) {
	var req service.OrgQueryUserRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ctx := c.Request.Context()
	req.OrgID = ginx.ParseParamID(c, "id")
	user, err := a.OrgSvc.QueryAppUser(ctx, &req)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, user)
}

// @Summary 更新组织应用权限
// @Description 更新组织应用权限
// @Accept json
// @Produce json
// @Tags 组织
// @Param id path integer true "org id"
// @Param body body service.OrgUserRoleUpdateRequest true "app roles"
// @Success 200 {object} model.UserOrgDTOList
// @Router /orgs/{id}/roles [put]
func (a *OrgAPI) UpdateUserRoles(c *gin.Context) {
	var req service.OrgUserRoleUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ctx := c.Request.Context()
	req.OrgID = ginx.ParseParamID(c, "id")
	// TODO 后续移除 token 返回
	token, err := a.OrgSvc.UpdateUserRoles(ctx, ginx.GetUserID(c), &req)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, gin.H{"token": token})
}

// @Summary 获取已删除应用列表
// @Description 获取已删除应用列表
// @Produce json
// @Tags 组织
// @Param id path integer true "org id"
// @Success 200 {object} model.AppDTOList
// @Router /orgs/{id}/deleted-apps [get]
func (a *OrgAPI) ListDeletedApps(c *gin.Context) {
	ctx := c.Request.Context()
	list, err := a.OrgSvc.ListDeletedApps(ctx, ginx.GetUserID(c), ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, list)
}

// @Summary 恢复已删除应用
// @Description 恢复已删除应用
// @Produce json
// @Tags 组织
// @Param id path integer true "org id"
// @Param body body service.OrgRestoreDeletedAppsRequest true "body"
// @Success 200 {object} object "{}"
// @Router /orgs/{id}/deleted-apps [put]
func (a *OrgAPI) RestoreDeletedApps(c *gin.Context) {
	var req service.OrgRestoreDeletedAppsRequest
	if err := c.ShouldBind(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ctx := c.Request.Context()
	err := a.OrgSvc.RestoreDeletedApps(ctx, ginx.GetUserID(c), ginx.ParseParamID(c, "id"), &req)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, gin.H{})
}
