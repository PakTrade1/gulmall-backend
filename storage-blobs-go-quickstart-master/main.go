package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// Azure Storage Quickstart Sample - Demonstrate how to upload, list, download, and delete blobs.
//
// Documentation References:
// - What is a Storage Account - https://docs.microsoft.com/azure/storage/common/storage-create-storage-account
// - Blob Service Concepts - https://docs.microsoft.com/rest/api/storageservices/Blob-Service-Concepts
// - Blob Service Go SDK API - https://godoc.org/github.com/Azure/azure-storage-blob-go
// - Blob Service REST API - https://docs.microsoft.com/rest/api/storageservices/Blob-Service-REST-API
// - Scalability and performance targets - https://docs.microsoft.com/azure/storage/common/storage-scalability-targets
// - Azure Storage Performance and Scalability checklist https://docs.microsoft.com/azure/storage/common/storage-performance-checklist
// - Storage Emulator - https://docs.microsoft.com/azure/storage/common/storage-use-emulator

func randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Int())
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	url := "https://paktradegallery.blob.core.windows.net/" //replace <StorageAccountName> with your Azure storage account name
	ctx := context.Background()
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	containerName := "abbasi/img"
	fmt.Printf("Creating a container named %s\n", containerName)
	fmt.Printf("Creating a dummy file to test the upload and download\n")

	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, _, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	// fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	// fmt.Printf("File Size: %+v\n", handler.Size)
	// fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(fileBytes)
	//data := []byte("\nhello world this is a blob\n")

	blobName := "mamals"
	blobClient, err := azblob.NewBlockBlobClient(url+containerName+"/"+blobName, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Upload to data to blob storage
	_, err = blobClient.UploadBufferToBlockBlob(ctx, fileBytes, azblob.HighLevelUploadToBlockBlobOption{})

	if err != nil {
		log.Fatalf("Failure to upload to blob: %+v", err)
	}

	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}

func main() {
	// fmt.Printf("Azure Blob storage quick start sample\n")
	// url := "https://paktradegallery.blob.core.windows.net/" //replace <StorageAccountName> with your Azure storage account name
	// ctx := context.Background()

	// // Create a default request pipeline using your storage account name and account key.
	// credential, err := azidentity.NewDefaultAzureCredential(nil)
	// if err != nil {
	// 	log.Fatal("Invalid credentials with error: " + err.Error())
	// }

	// // serviceClient, err := azblob.NewServiceClient(url, credential, nil)
	// // if err != nil {
	// // 	log.Fatal("Invalid credentials with error: " + err.Error())
	// // }

	// // Create the container
	// containerName := "abbasi/abc"
	// fmt.Printf("Creating a container named %s\n", containerName)
	// // containerClient := serviceClient.NewContainerClient(containerName)
	// // _, err = containerClient.Create(ctx, nil)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }

	// fmt.Printf("Creating a dummy file to test the upload and download\n")

	// data := []byte("\nhello world this is a blob\n")
	// blobName := "abbasi" + "-" + randomString()

	// blobClient, err := azblob.NewBlockBlobClient(url+containerName+"/"+blobName, credential, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Upload to data to blob storage
	// _, err = blobClient.UploadBufferToBlockBlob(ctx, data, azblob.HighLevelUploadToBlockBlobOption{})

	// if err != nil {
	// 	log.Fatalf("Failure to upload to blob: %+v", err)
	// }

	// List the blobs in the container
	//fmt.Println("Listing the blobs in the container:")

	// pager := containerClient.ListBlobsFlat(nil)

	// for pager.NextPage(ctx) {
	// 	resp := pager.PageResponse()

	// 	for _, v := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
	// 		fmt.Println(*v.Name)
	// 	}
	// }

	// if err = pager.Err(); err != nil {
	// 	log.Fatalf("Failure to list blobs: %+v", err)
	// }

	// Download the blob
	// get, err := blobClient.Download(ctx, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// downloadedData := &bytes.Buffer{}
	// reader := get.Body(azblob.RetryReaderOptions{})
	// _, err = downloadedData.ReadFrom(reader)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = reader.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(downloadedData.String())

	// fmt.Printf("Press enter key to delete the blob fils, example container, and exit the application.\n")
	// bufio.NewReader(os.Stdin).ReadBytes('\n')
	// fmt.Printf("Cleaning up.\n")

	// // Delete the blob
	// fmt.Printf("Deleting the blob " + blobName + "\n")

	// _, err = blobClient.Delete(ctx, nil)
	// if err != nil {
	// 	log.Fatalf("Failure: %+v", err)
	// }

	// Delete the container
	// fmt.Printf("Deleting the blob " + containerName + "\n")
	// _, err = containerClient.Delete(ctx, nil)

	// if err != nil {
	// 	log.Fatalf("Failure: %+v", err)
	// }
	fmt.Println("Hello World")
	setupRoutes()

}
