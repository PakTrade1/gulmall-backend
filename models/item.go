package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PublicId     int                `json:"publicId,omitempty" bson:"publicId,omitempty"`
	Category     primitive.ObjectID `json:"category,omitempty" bson:"category,omitempty"`
	Country      string             `json:"country,omitempty" bson:"country,omitempty"`
	Currency     string             `json:"currency,omitempty" bson:"currency,omitempty"`
	Images       []Image            `json:"images,omitempty" bson:"images,omitempty"`
	OwnerId      primitive.ObjectID `json:"ownerId,omitempty" bson:"ownerId,omitempty"`
	PlanId       primitive.ObjectID `json:"planId,omitempty" bson:"planId,omitempty"`
	Price        float64            `json:"price,omitempty" bson:"price,omitempty"`
	Qty          int                `json:"qty,omitempty" bson:"qty,omitempty"`
	RemainingQty int                `json:"remainingQty,omitempty" bson:"remainingQty,omitempty"`
	Status       string             `json:"status,omitempty" bson:"status,omitempty"`
	SubCategory  primitive.ObjectID `json:"subCategory,omitempty" bson:"subCategory,omitempty"`
	Title        string             `json:"title,omitempty" bson:"title,omitempty"`
}

type Image struct {
	Image string `json:"image,omitempty" bson:"image,omitempty"`
	Color string `json:"color,omitempty" bson:"color,omitempty"`
}
