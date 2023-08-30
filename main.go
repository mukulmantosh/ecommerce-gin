package main

import (
	"log"
	"os"
)
import (
	"github.com/mukulmantosh/ecommerce-gin/controllers"

	"github.com/gin-gonic/gin"
	"github.com/mukulmantosh/ecommerce-gin/database"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products", database.UserData(database.Client, "Users")))

	router := gin.New()
	router.Use(gin.Logger())
	router.UserRoutes(routes)
	router.use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))

}
