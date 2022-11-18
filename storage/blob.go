package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	// Your account name and key can be obtained from the Azure Portal.
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	accountKey, ok := os.LookupEnv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY could not be found")
	}
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	// The service URL for blob endpoints is usually in the form: http(s)://<account>.blob.core.windows.net/
	client, err := azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	handleError(err)

	// ===== 1. Create a container =====
	containerName := "testcontainer"
	containerCreateResp, err := client.CreateContainer(context.TODO(), containerName, nil)
	handleError(err)
	fmt.Println(containerCreateResp)

	// ===== 2. Upload and Download a block blob =====
	blobData := "Hello world!"
	blobName := "HelloWorld.txt"
	uploadResp, err := client.UploadStream(context.TODO(),
		containerName,
		blobName,
		strings.NewReader(blobData),
		&azblob.UploadStreamOptions{
			Metadata: map[string]string{"Foo": "Bar"},
			Tags:     map[string]string{"Year": "2022"},
		})
	handleError(err)
	fmt.Println(uploadResp)

	// Download the blob's contents and ensure that the download worked properly
	blobDownloadResponse, err := client.DownloadStream(context.TODO(), containerName, blobName, nil)
	handleError(err)

	// Use the bytes.Buffer object to read the downloaded data.
	// RetryReaderOptions has a lot of in-depth tuning abilities, but for the sake of simplicity, we'll omit those here.
	reader := blobDownloadResponse.Body
	downloadData, err := io.ReadAll(reader)
	handleError(err)
	if string(downloadData) != blobData {
		log.Fatal("Uploaded data should be same as downloaded data")
	}

	err = reader.Close()
	if err != nil {
		return
	}

	// ===== 3. List blobs =====
	// List methods returns a pager object which can be used to iterate over the results of a paging operation.
	// To iterate over a page use the NextPage(context.Context) to fetch the next page of results.
	// PageResponse() can be used to iterate over the results of the specific page.
	pager := client.NewListBlobsFlatPager(containerName, nil)
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		handleError(err)
		for _, v := range resp.Segment.BlobItems {
			fmt.Println(*v.Name)
		}
	}

	// Delete the blob.
	_, err = client.DeleteBlob(context.TODO(), containerName, blobName, nil)
	handleError(err)

	// Delete the container.
	_, err = client.DeleteContainer(context.TODO(), containerName, nil)
	handleError(err)
}
