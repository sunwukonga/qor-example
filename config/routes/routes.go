package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qor/qor"
	"github.com/qor/qor/utils"
	"github.com/qor/wildcard_router"
	"github.com/sunwukonga/qor-example/app/controllers"
	"github.com/sunwukonga/qor-example/config"
	"github.com/sunwukonga/qor-example/config/auth"
	"github.com/sunwukonga/qor-example/db"
)

var rootMux *http.ServeMux
var WildcardRouter *wildcard_router.WildcardRouter

func Router() *http.ServeMux {
	if rootMux == nil {
		router := gin.Default()
		router.Use(func(ctx *gin.Context) {
			if locale := utils.GetLocale(&qor.Context{Request: ctx.Request, Writer: ctx.Writer}); locale != "" {
				ctx.Set("DB", db.DB.Set("l10n:locale", locale))
			}
		})
		gin.SetMode(gin.DebugMode)

		router.GET("/", controllers.HomeIndex)
		router.GET("/products/:code", controllers.ProductShow)
		router.GET("/cart", controllers.CartShow)
		router.GET("/switch_locale", controllers.SwitchLocale)

		router.GET("/addtocart/:id", controllers.AddToCart)

		rootMux = http.NewServeMux()
		rootMux.Handle("/auth/", auth.Auth.NewRouter())
		publicDir := http.Dir(strings.Join([]string{config.Root, "public"}, "/"))
		rootMux.Handle("/dist/", http.FileServer(publicDir))
		rootMux.Handle("/vendors/", http.FileServer(publicDir))
		rootMux.Handle("/images/", http.FileServer(publicDir))
		rootMux.Handle("/fonts/", http.FileServer(publicDir))

		WildcardRouter = wildcard_router.New()
		WildcardRouter.MountTo("/", rootMux)
		WildcardRouter.AddHandler(router)
	}
	return rootMux
}
