package tokens

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mukulmantosh/ecommerce-gin/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var UserData *mongo.Collection = database.UserData(database.Client, "Users")

var SecretKey string = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UID       string
	jwt.RegisteredClaims
}

func TokenGenerator(email string, firstname string, lastname string, uid string) (signedToken string, signedRefreshToken string, err error) {
	timer := time.Now().Local().Add(time.Hour * time.Duration(2)).Unix()

	claims := &SignedDetails{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		UID:       uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(timer, 0)),
		},
	}

	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		UID:       uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(timer, 0)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SecretKey))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SecretKey))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("The token is invalid")
		msg = err.Error()
		return
	}
	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{"AccessToken", signedToken})
	updateObj = append(updateObj, bson.E{"RefreshToken", signedRefreshToken})

	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"UpdatedAt", UpdatedAt})
	upsert := true

	filter := bson.M{"UserID": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := UserData.UpdateOne(ctx, filter,
		bson.D{
			{"$set", updateObj},
		}, &opt)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
	return

}
