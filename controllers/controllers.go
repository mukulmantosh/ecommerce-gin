package controllers

import "C"
import (
	"context"
	"ecommerce-gin/database"
	"ecommerce-gin/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	}
	return valid, msg

}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.Bind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists!"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		if err != nil {
			log.Panic(err)
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this phone number is already in use"})
			return
		}

		password := HashPassword(user.Password)
		user.Password = password
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()

		token, refreshToken, _ := generate.TokenGenerator(user.Email, user.FirstName, user.LastName, user.UserID)

		user.Token = &token
		user.RefreshToken = &refreshToken
		user.UserCart = make([]models.ProductUser, 0)
		user.AddressDetails = make([]models.Address, 0)
		user.OrderStatus = make([]models.Order, 0)
		_, insertErr := UserCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user"})
			return
		}

		c.JSON(http.StatusCreated, "Successfully signed in!")

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}

		PasswordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)

		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			log.Println(msg)
			return
		}
		token, refreshToken, _ := generate.TokenGenerator(foundUser.Email, foundUser.FirstName,
			foundUser.LastName, foundUser.UserID)

		generate.UpdateAllTokens(token, refreshToken, foundUser.UserID)

		c.JSON(http.StatusFound, foundUser)
		return

	}
}

func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something went wrong,"+
				"please try after sometime.")
			return
		}

		err = cursor.All(ctx, &productList)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid")
			return
		}

		c.IndentedJSON(http.StatusOK, productList)

	}
}

func SearchProductByQuery() gin.HandlerFunc {

}
