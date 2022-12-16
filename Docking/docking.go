package docking

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var PakTradeDb *mongo.Database // PakTrade

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

	pakTrade := client.Database("PakTrade")
	PakTradeDb = pakTrade
}

func AzureBloblogs() *azblob.Client {
	url := "https://appmedia.blob.core.windows.net/"
	credential, err1 := azidentity.NewDefaultAzureCredential(nil)
	if err1 != nil {
		log.Fatal("Invalid credentials with error: " + err1.Error())
	}
	client, client_error := azblob.NewClient(url, credential, nil)
	if client_error != nil {
		log.Fatal("Invalid credentials with error: ", client_error.Error())
	}
	return client
}
