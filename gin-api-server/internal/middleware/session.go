package middleware

import (
	"gin-api-server/internal/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SessionMiddleware(cfg *config.HTTPSessionCfg) gin.HandlerFunc {
	store := cookie.NewStore([]byte(cfg.Secret))
	store.Options(sessions.Options{
		MaxAge:   int(cfg.ExpireDuration),
		Path:     cfg.Path,
		HttpOnly: true,
	})
	return sessions.Sessions(cfg.Name, store)
}
