package middleware

import (
	"gin-api-server/internal/extend/ginx"
	"gin-api-server/internal/extend/svcerrorsx"
	"gin-api-server/internal/extend/timex"
	"gin-api-server/internal/model"
	"gin-api-server/pkg/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(a auth.Auther, skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		userID := ginx.GetUserID(c)
		if userID == 0 {
			session := sessions.Default(c)
			userIDStr := session.Get(auth.AuthnUserIDKey)
			if userIDStr != nil {
				if id, ok := userIDStr.(uint64); ok {
					userID = id
				}
			}

			// 记录用户最后活跃时间
			if userID > 0 {
				lastActiveAtStr := session.Get(auth.AuthnUserLastActiveAtKey)
				lastActiveAt, ok := lastActiveAtStr.(int64)
				now := timex.GetNowUTCMilli()
				if !ok || now > lastActiveAt+10*60*1000 {
					if err := a.GetDB().Model(&model.User{}).Where("id = ?", userID).Update("last_active_at", now).Error; err != nil {
						c.Abort()
						ginx.RespondErr(c, err)
						return
					}
					session.Set(auth.AuthnUserLastActiveAtKey, now)
					if err := session.Save(); err != nil {
						c.Abort()
						ginx.RespondErr(c, err)
						return
					}
				}
			}
		}

		if userID == 0 {
			c.Abort()
			ginx.RespondErr(c, svcerrorsx.ErrUserUnauthorized)
			return
		}

		ginx.SetUserID(c, userID)
		c.Next()
	}
}
