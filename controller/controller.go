package controller

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/naiba/nbdomain"
	"github.com/naiba/nbdomain/controller/cat"
	"github.com/naiba/nbdomain/controller/domain"
	"github.com/naiba/nbdomain/controller/mibiao"
	"github.com/naiba/nbdomain/controller/offer"
	"github.com/naiba/nbdomain/controller/panel"
	"github.com/naiba/nbdomain/controller/user"
	"github.com/naiba/nbdomain/controller/whois"
	"github.com/naiba/nbdomain/model"
	"github.com/naiba/nbdomain/pkg/mygin"
)

// Web start
func Web() {
	var mode string
	if nbdomain.CF.Debug {
		mode = gin.DebugMode
	} else {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"toLower": strings.ToLower,
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
		},
	})
	r.LoadHTMLGlob("theme/template/**/*")
	r.Static("static", "theme/static")
	r.Static("upload", "data/upload")

	if nbdomain.CF.Debug {
		conf := cors.DefaultConfig()
		conf.AllowAllOrigins = true
		conf.AddAllowMethods("DELETE")
		conf.AddAllowHeaders("Authorization")
		r.Use(cors.New(conf))
	}

	r.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "http://"+c.Request.URL.Hostname())
	})

	panelRouter := r.Group("/")
	{
		panelRouter.GET("", mibiao.Index)
		panelRouter.GET("offer/:domain", mibiao.Offer)
		panelRouter.POST("offer/:domain", mibiao.Offer)
		panelRouter.GET("offer/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/")
		})
		panelRouter.GET("allowed", mibiao.Allow)
	}

	api := r.Group("/api")
	{
		api.POST("login", user.Login)
		authUser := api.Group("")
		{
			authUser.Use(mygin.Authorize(mygin.AuthOption{NeedUser: true}))
			// 用户
			authUser.GET("user", user.GET)
			authUser.PUT("user", user.Settings)
			authUser.POST("logout", user.Logout)
			// 米表
			authUser.GET("panel", panel.List)
			authUser.DELETE("panel/:id", panel.Delete)
			authUser.POST("panel", panel.Edit)
			authUser.POST("panel/import", panel.Import)
			authUser.GET("panel/:id/export", panel.Export)
			// 分类
			authUser.GET("cat", cat.List)
			authUser.POST("cat", cat.Edit)
			authUser.DELETE("cat/:id", cat.Delete)
			// 域名
			authUser.GET("domain", domain.List)
			authUser.POST("domain", domain.Edit)
			authUser.DELETE("domain/:id", domain.Delete)
			// 销售
			authUser.GET("offer", offer.List)
			authUser.DELETE("offer/:id", offer.Delete)
			// 其他
			authUser.GET("panel_option", func(c *gin.Context) {
				var r model.Response
				r.Code = http.StatusOK
				r.Result = gin.H{
					"themes":         model.ThemeList,
					"offer_themes":   model.OfferThemeList,
					"analysis_types": model.AnalysisTypes,
				}
				c.JSON(http.StatusOK, r)
			})
			authUser.GET("whois/:domain", whois.Whois)
		}
	}
	go r.Run(nbdomain.CF.Web.Addr)
}
