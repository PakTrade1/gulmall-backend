package docking

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var PakTradeDb *mongo.Database // PakTrade

func PakTradeConnection() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://developer001:IAmMuslim@cluster0.qeqntol.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
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
	url := "https://paktradestorage.blob.core.windows.net/"
	azrKey := "H2imTpv/29KoIVfRMOrbNyRr5MDDUU7AbO99oL+x6L30e53houKz8PEOrhQuSMGM4U7mqdmHZ+My+ASty1aMYg=="
	azrBlobAccountName := "paktradestorage"
	credentialShared, errC := azblob.NewSharedKeyCredential(azrBlobAccountName, azrKey)
	if errC != nil {
		log.Fatal("Invalid credentials with error: " + errC.Error())
	}
	client, client_error := azblob.NewClientWithSharedKeyCredential(url, credentialShared, nil)
	if client_error != nil {
		log.Fatal("Invalid credentials with error: ", client_error.Error())
	}
	return client
}
