package controllers

import (
	"context"
	"ecommerce-gin/database"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{prodCollection: prodCollection, userCollection: userCollection}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))

			userQueryID := c.Query("userID")
			if userQueryID == "" {
				log.Println("user id is empty")
				c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
				return
			}

			productID, err := primitive.ObjectIDFromHex(productQueryID)

			if err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
			}
			c.IndentedJSON(http.StatusOK, "successfully added")
		}
	}
}

func RemoveItem(app *Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))

			userQueryID := c.Query("userID")
			if userQueryID == "" {
				log.Println("user id is empty")
				c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
				return
			}

			productID, err := primitive.ObjectIDFromHex(productQueryID)

			if err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryID)

			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}
			c.IndentedJSON(http.StatusOK, "Successfully removed item from cart")

		}
	}
}

func GetItemFromCart() gin.HandlerFunc {

}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")

		if userQueryID == "" {
			log.Panic("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "successfully placed the order")

	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))

			userQueryID := c.Query("userID")
			if userQueryID == "" {
				log.Println("user id is empty")
				c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
				return
			}

			productID, err := primitive.ObjectIDFromHex(productQueryID)

			if err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, userQueryID)

			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}
			c.IndentedJSON(http.StatusOK, "successfully placed the order")

		}
	}
}
