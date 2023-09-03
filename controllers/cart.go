package controllers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mukulmantosh/ecommerce-gin/database"
	"github.com/mukulmantosh/ecommerce-gin/models"
	"go.mongodb.org/mongo-driver/bson"
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

func (app *Application) RemoveItem() gin.HandlerFunc {
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
	return func(c *gin.Context) {
		userId := c.Query("id")

		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}

		userInfo, _ := primitive.ObjectIDFromHex(userId)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var fillCart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userInfo}}).Decode(&fillCart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(404, "not found")
			return
		}

		filterMatch := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userInfo}}}}
		Unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$user_cart"}}}}
		Grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$user_cart.price"}}}}}}

		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filterMatch, Unwind, Grouping})
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		if err = pointCursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, fillCart.UserCart)
		}
		ctx.Done()
	}
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
