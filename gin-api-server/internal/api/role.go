package api

import (
	"gin-api-server/internal/extend/ginx"
	"gin-api-server/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleAPI struct {
	RoleSvc *service.RoleSvc
}

// @Summary 获取全部角色列表
// @Description 获取全部角色列表
// @Produce json
// @Tags 角色
// @Success 200 {object} model.RoleDTOList
// @Router /roles [get]
func (a *RoleAPI) List(c *gin.Context) {
	ctx := c.Request.Context()
	list, err := a.RoleSvc.ListAll(ctx)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, list)
}

// func (a *RoleAPI) Create(c *gin.Context) {
// 	var req service.RoleCreateRequest
// 	if err := c.ShouldBind(&req); err != nil {
// 		ginx.RespondErr(c, err)
// 		return
// 	}

// 	ctx := c.Request.Context()
// 	org, err := a.RoleSvc.Create(ctx, ginx.GetUserID(c), &req)
// 	if err != nil {
// 		ginx.RespondErr(c, err)
// 		return
// 	}

// 	ginx.RespondData(c, http.StatusCreated, org)
// }

// func (a *RoleAPI) Update(c *gin.Context) {
// 	var req service.RoleUpdateRequest
// 	if err := c.ShouldBind(&req); err != nil {
// 		ginx.RespondErr(c, err)
// 		return
// 	}

// 	ctx := c.Request.Context()
// 	req.ID = ginx.ParseParamID(c, "id")
// 	org, err := a.RoleSvc.Update(ctx, ginx.GetUserID(c), &req)
// 	if err != nil {
// 		ginx.RespondErr(c, err)
// 		return
// 	}

// 	ginx.RespondData(c, http.StatusOK, org)
// }

// func (a *RoleAPI) Delete(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	if err := a.RoleSvc.Delete(ctx, ginx.GetUserID(c), ginx.ParseParamID(c, "id")); err != nil {
// 		ginx.RespondErr(c, err)
// 		return
// 	}

// 	ginx.RespondData(c, http.StatusNoContent, gin.H{})
// }
