package api

import (
	"gin-api-server/internal/extend/ginx"
	"gin-api-server/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterAPI struct {
	UserSvc *service.UserSvc
}

// @Summary 用户注册
// @Description 用户注册
// @Accept json
// @Produce json
// @Tags 注册
// @Param body body service.UserRegisterEmailRequest true "body"
// @Success 200 {object} object "{"message": "your are logged in"}"
// @Router /register [post]
func (a *RegisterAPI) RegisterEmail(c *gin.Context) {
	var req service.UserRegisterEmailRequest
	if err := c.ShouldBind(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}
	ctx := c.Request.Context()
	token, err := a.UserSvc.CreateRegisterToken(ctx, req.Email)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, gin.H{
		"message": "",
		"token":   token,
	})
}

// @Summary 用户注册 token 信息
// @Description 用户注册 token 信息
// @Accept json
// @Produce json
// @Tags 注册
// @Param token query string true "register token"
// @Success 200 {object} model.TokenDTO
// @Router /register/token [get]
func (a *RegisterAPI) GetRegisterToken(c *gin.Context) {
	var req service.UserRegisterTokenRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}
	ctx := c.Request.Context()
	token, err := a.UserSvc.GetRegisterToken(ctx, req.Token)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, token)
}

// @Summary 用户完善注册信息
// @Description 用户完善注册信息
// @Accept json
// @Produce json
// @Tags 注册
// @Param body body service.UserRegisterProfileRequest true "body"
// @Success 200 {object} model.UserDTO
// @Router /register/profile [post]
func (a *RegisterAPI) RegisterProfile(c *gin.Context) {
	var req service.UserRegisterProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ctx := c.Request.Context()
	user, err := a.UserSvc.Create(ctx, &req)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ginx.RespondData(c, http.StatusOK, user)
}
