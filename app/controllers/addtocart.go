package controllers

import (
	"encoding/json"
	"log"
	"strconv"
	//	"html/template"
	//	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunwukonga/qor-example/app/models"
	"github.com/sunwukonga/qor-example/config/auth" // for sessionStore
)

const SessionCartKey string = "cart_key"
const YomummaKey string = "yomumma"

func AddToCart(ctx *gin.Context) {

	var (
		cart          []models.OrderItem
		sessionStorer *auth.SessionStorer
		productId     uint
		product       models.Product
	)

	// Get cart from session. Automatically creates empty cart even if none in session.
	sessionStorer = auth.NewSessionStorer(ctx.Writer, ctx.Request).(*auth.SessionStorer)
	if cartText, ok := sessionStorer.Get(SessionCartKey); ok {
		json.Unmarshal([]byte(cartText), &cart)
	} else {
		log.Printf("Could not get key from session.")
	}
	//Test session cookie
	sessionStorer.Put(YomummaKey, "Kings are made, not born.")

	productIdStr := ctx.Param("id")
	log.Printf("Product ID passed: %v", productIdStr)
	productIduint64, _ := strconv.ParseUint(productIdStr, 10, 32)
	productId = uint(productIduint64)
	log.Printf("Product ID after strconv: %v", productId)

	// If product already exists in cart, ignore, else add to cart.
	// Quantity is never incremented.
	for _, orderItem := range cart {
		log.Println("Checking session cart for product ...")
		if orderItem.ProductID == productId {
			// signal that product should not be added to cart as it already exists
			productId = 0
			log.Println("We found the product in the cart already.")
			break
		}
	}

	if productId > 0 {
		// Fetch product details from the database.
		DB(ctx).First(&product, productId)
		log.Printf("Product %v fetched from database", productId)
		// add a new product to the cart
		cart = append(cart, models.OrderItem{ProductID: productId, Product: product, Quantity: 1, Price: product.Price})
		log.Println(cart)
		log.Println("------")
		cartBytes, err := json.Marshal(&cart)
		if err != nil {
			log.Printf("Marshalling error: %v", err)
		}
		log.Println(string(cartBytes))
		log.Println("------")
		sessionStorer.Put(SessionCartKey, string(cartBytes))

	}

	redirectBack(ctx.Writer, ctx.Request)

}
