package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://developer001:IAmMuslim@cluster0.qeqntol.mongodb.net/test")

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Open an aggregation cursor
	coll := client.Database("PakTrade").Collection("cloths")
	_, err = coll.Aggregate(ctx, bson.A{
		bson.D{},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "color"},
					{"localField", "abc"},
					{"foreignField", "_id"},
					{"as", "color"},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
