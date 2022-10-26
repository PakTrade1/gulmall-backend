package main

import (
	"context"
	"example/apies/controllers"
	"example/apies/services"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server          *gin.Engine
	colorservice    services.ColorService
	colorcontroller controllers.ColorController
	ctx             context.Context
	colorcollection *mongo.Collection
	mongoclient     *mongo.Client
	err             error
)

func init() {
	ctx = context.TODO()
	mongoconection := options.Client().ApplyURI("mongodb+srv://developer001:IAmMuslim@cluster0.qeqntol.mongodb.net/test")
	mongoclient, err = mongo.Connect(ctx, mongoconection)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("database conect")
	colorcollection = mongoclient.Database("PakTrade").Collection("color")
	colorservice = services.NewColorService(colorcollection, ctx)
	colorcontroller = controllers.New(colorservice)
	server = gin.Default()

}

func main() {
	defer mongoclient.Disconnect(ctx)
	basepath := server.Group("v1")
	colorcontroller.RegisterColorRouts(basepath)
	log.Fatal(server.Run(":9090"))
 
}
