package docking

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var PakTradeDb *mongo.Database // PakTrade
var ItemDb *mongo.Database     // Items
var CartsDb *mongo.Database    // Items

func PakTradeConnection() {
	fmt.Print("\nCalling fun db connect\n")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://developer001:IAmMuslim@cluster0.qeqntol.mongodb.net/test"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// DATABASE NAME

	pakTrade := client.Database("PakTrade")
	items := client.Database("Item")
	carts := client.Database("carts")
	ItemDb = items
	PakTradeDb = pakTrade
	CartsDb = carts

}
