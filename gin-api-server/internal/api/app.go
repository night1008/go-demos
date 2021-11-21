package api

import (
	"gin-api-server/internal/extend/ginx"
	"gin-api-server/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppAPI struct {
	AppSvc *service.AppSvc
}

// @Summary 创建应用
// @Description 创建应用
// @Accept json
// @Produce json
// @Tags 应用
// @Param body body service.AppCreateRequest true "body"
// @Success 200 {object} model.AppDTO
// @Router /apps [post]
func (a *AppAPI) Create(c *gin.Context) {
	var req service.AppCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ctx := c.Request.Context()
	org, err := a.AppSvc.Create(ctx, ginx.GetUserID(c), &req)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}
	ginx.RespondData(c, http.StatusCreated, org)
}

// @Summary 获得应用
// @Description 获得应用
// @Produce json
// @Tags 应用
// @Param id path integer true "app id"
// @Success 200 {object} model.AppDTO
// @Router /apps/{id} [get]
func (a *AppAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	app, err := a.AppSvc.Get(ctx, ginx.GetUserID(c), ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, app)
}

// @Summary 更新应用
// @Description 更新应用
// @Accept json
// @Produce json
// @Tags 应用
// @Param id path integer true "app id"
// @Param body body service.AppUpdateRequest true "body"
// @Success 200 {object} model.AppDTO
// @Router /apps/{id} [put]
func (a *AppAPI) Update(c *gin.Context) {
	var req service.AppUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ctx := c.Request.Context()
	req.ID = ginx.ParseParamID(c, "id")
	org, err := a.AppSvc.Update(ctx, ginx.GetUserID(c), &req)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, org)
}

// @Summary 删除应用
// @Description 删除应用
// @Produce json
// @Tags 应用
// @Param id path integer true "app id"
// @Success 204 {object} object "{}"
// @Router /apps/{id} [delete]
func (a *AppAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	if err := a.AppSvc.Delete(ctx, ginx.GetUserID(c), ginx.ParseParamID(c, "id")); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusNoContent, gin.H{})
}
