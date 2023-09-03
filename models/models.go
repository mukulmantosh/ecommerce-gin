package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id" json:"_id"`
	FirstName      string             `json:"first_name" validate:"required,min=2,max=30"`
	LastName       string             `json:"last_name" validate:"required,min=2,max=30"`
	Password       string             `json:"password" validate:"required,min=6,max=15"`
	Email          string             `json:"email" validate:"required"`
	Phone          string             `json:"phone" validate:"required"`
	Token          string             `json:"token"`
	RefreshToken   string             `json:"refresh_token"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	UserID         string             `json:"user_id"`
	UserCart       []ProductUser      `json:"user_cart" bson:"user_cart"`
	AddressDetails []Address          `json:"address_details" bson:"address_details"`
	OrderStatus    []Order            `json:"order_status" bson:"order_status"`
}

type Product struct {
	ProductID   primitive.ObjectID `bson:"_id" json:"_id"`
	ProductName string             `json:"product_name" bson:"product_name"`
	Price       uint64             `json:"price"`
	Rating      uint8              `json:"rating"`
	Image       string             `json:"image"`
}

type ProductUser struct {
	ProductID   primitive.ObjectID `bson:"_id" json:"_id"`
	ProductName string             `json:"product_name" bson:"product_name"`
	Price       uint64             `json:"price"`
	Rating      uint8              `json:"rating"`
	Image       string             `json:"image"`
}

type Address struct {
	AddressID primitive.ObjectID `bson:"_id" json:"_id"`
	House     string             `json:"house"`
	Street    string             `json:"street"`
	City      string             `json:"city"`
	PinCode   string             `json:"pin_code"`
}

type Order struct {
	OrderID       primitive.ObjectID `bson:"_id" json:"_id"`
	OrderCart     []ProductUser      `json:"order_list" bson:"order_list"`
	OrderedAt     time.Time          `json:"ordered_at" bson:"ordered_at"`
	Price         uint64             `json:"total_price" bson:"total_price"`
	Discount      int                `json:"discount" bson:"discount"`
	PaymentMethod Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool `json:"digital"`
	COD     bool `json:"cod"`
}
