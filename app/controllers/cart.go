package controllers

import (
	"encoding/json"
	"log"
	//	"html/template"
	//	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunwukonga/qor-example/app/models"
	"github.com/sunwukonga/qor-example/config"
	"github.com/sunwukonga/qor-example/config/admin"
	"github.com/sunwukonga/qor-example/config/auth" // for sessionStore
)

func CartShow(ctx *gin.Context) {
	var (
		cart           []models.OrderItem
		sessionStorer  *auth.SessionStorer
		colorVariation models.ColorVariation
	)

	DB(ctx).Preload("Product").Preload("Color").Preload("SizeVariations.Size").First(&colorVariation)

	sessionStorer = auth.NewSessionStorer(ctx.Writer, ctx.Request).(*auth.SessionStorer)
	if cartText, ok := sessionStorer.Get(SessionCartKey); ok {
		log.Printf("Cart from session: %v", cartText)
		json.Unmarshal([]byte(cartText), &cart)
	} else {
		log.Printf("Could not get key from session.")
	}
	if yomumma, ok := sessionStorer.Get(YomummaKey); ok {
		log.Printf("We got the message: %v", yomumma)
	}

	config.View.Funcs(funcsMap(ctx)).Execute(
		"cart_show",
		gin.H{
			"ActionBarTag":   admin.ActionBar.Render(ctx.Writer, ctx.Request),
			"Cart":           cart,
			"ColorVariation": colorVariation,
			"CurrentUser":    CurrentUser(ctx),
			"CurrentLocale":  CurrentLocale(ctx),
		},
		ctx.Request,
		ctx.Writer,
	)
}

// func funcsMap(ctx *gin.Context) template.FuncMap {
