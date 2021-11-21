package api

import (
	"gin-api-server/internal/extend/ginx"
	"gin-api-server/internal/service"
	"gin-api-server/pkg/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginAPI struct {
	UserSvc *service.UserSvc
}

// @Summary 用户登录
// @Description 用户登录
// @Accept json
// @Produce json
// @Tags 登录
// @Param body body service.UserLoginRequest true "body"
// @Success 200 {object} object "{"message": "your are logged in"}"
// @Router /login [post]
func (a *LoginAPI) Login(c *gin.Context) {
	var req service.UserLoginRequest
	if err := c.ShouldBind(&req); err != nil {
		ginx.RespondErr(c, err)
		return
	}

	ctx := c.Request.Context()
	user, err := a.UserSvc.Verify(ctx, req.Email, req.Password)
	if err != nil {
		ginx.RespondErr(c, err)
		return
	}
	session := sessions.Default(c)
	session.Set(auth.AuthnUserIDKey, user.ID)
	session.Set(auth.AuthnUserLastActiveAtKey, 0)
	if err := session.Save(); err != nil {
		// logger.WithContext(c.Request.Context()).Error().Err(err).Msg("failed to save session")
		ginx.RespondErr(c, err)
		return
	}
	ginx.RespondMsg(c, "your are logged in")
}

// @Summary 用户登出
// @Description 用户登出
// @Produce json
// @Tags 登录
// @Success 200 {object} object "{"message": "your are logged out"}"
// @Router /logout [post]
func (a *LoginAPI) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		ginx.RespondErr(c, err)
		return
	}
	ginx.RespondMsg(c, "your are logged out")
}
