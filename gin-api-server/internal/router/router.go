package router

import (
	"gin-api-server/internal/api"
	"gin-api-server/internal/config"
	"gin-api-server/internal/middleware"
	"gin-api-server/pkg/auth"
	"reflect"
	"strings"

	_ "gin-api-server/internal/swagger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Router struct {
	Engine      *gin.Engine
	Auther      auth.Auther
	LoginAPI    *api.LoginAPI
	RegisterAPI *api.RegisterAPI
	UserAPI     *api.UserAPI
	OrgAPI      *api.OrgAPI
	AppAPI      *api.AppAPI
	RoleAPI     *api.RoleAPI
}

func (a *Router) Setup(cfg *config.HTTPServerCfg) {
	a.newEngine(cfg)
	a.register(cfg)
}

// MEMO 设置验证错误字段来源于 json or form tag，都没有设置的话取结构体字段名称
func setValidatorTagName() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			if name != "" {
				return name
			}
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func (a *Router) newEngine(cfg *config.HTTPServerCfg) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	gin.SetMode(cfg.RunMode)

	setValidatorTagName()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(middleware.SessionMiddleware(&cfg.Session))
	r.Use(middleware.AuthMiddleware(a.Auther, middleware.AllowPathPrefixSkipper("/api/v1/login", "/api/v1/register", "/api/swagger")))

	if cfg.CORS.Enable {
		r.Use(middleware.CORSMiddleware(&cfg.CORS))
	}

	if cfg.Swagger.Enable {
		r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)))
	}

	a.Engine = r
}

func (a *Router) register(cfg *config.HTTPServerCfg) {
	api := a.Engine.Group("api")
	v1 := api.Group("v1")
	{
		register := v1.Group("register")
		{
			register.POST("", a.RegisterAPI.RegisterEmail)
			register.GET("token", a.RegisterAPI.GetRegisterToken)
			register.POST("profile", a.RegisterAPI.RegisterProfile)
		}

		v1.POST("login", a.LoginAPI.Login)
		v1.POST("logout", a.LoginAPI.Logout)

		current := v1.Group("current")
		{
			current.GET("", a.UserAPI.Current)
			current.PUT("password", a.UserAPI.UpdateCurrentPassword)
			current.PUT("profile", a.UserAPI.UpdateCurrentProfile)
			current.GET("orgs", a.UserAPI.ListUserOrgs)
		}

		org := v1.Group("orgs")
		{
			org.POST("", a.OrgAPI.Create)
			org.PUT(":id", a.OrgAPI.Update)
			org.DELETE(":id", a.OrgAPI.Delete)
			org.GET(":id/owners", a.OrgAPI.ListOwners)
			org.GET(":id/apps", a.OrgAPI.ListApps)
			org.GET(":id/members", a.OrgAPI.ListAppMembers)
			org.GET(":id/users-by-query", a.OrgAPI.QueryUser)
			org.PUT(":id/roles", a.OrgAPI.UpdateUserRoles)
			org.GET(":id/deleted-apps", a.OrgAPI.ListDeletedApps)
			org.PUT(":id/deleted-apps", a.OrgAPI.RestoreDeletedApps)
		}

		app := v1.Group("apps")
		{
			app.POST("", a.AppAPI.Create)
			app.GET(":id", a.AppAPI.Get)
			app.PUT(":id", a.AppAPI.Update)
			app.DELETE(":id", a.AppAPI.Delete)
		}

		role := v1.Group("roles")
		{
			role.GET("", a.RoleAPI.List)
		}
	}
}
