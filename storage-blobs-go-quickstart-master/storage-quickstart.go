package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {

	setupRoutes()

}
func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/dirData", perems_data)

	http.ListenAndServe(":8080", nil)
}

type input_type1 struct {
	Id             string `josn:"id"`
	Collectionname string `json:"collectionname"`
	Containertype  string `json:"containertype"`
	// Fp             *os.File
}

var obid = ""
var namecontainer = ""
var typeblob = ""

func perems_data(w http.ResponseWriter, r *http.Request) {
	// conect to blob url
	var Perm_with_file input_type1
	err := json.NewDecoder(r.Body).Decode(&Perm_with_file)
	if err != nil {
		panic(err)
	}

	obid = Perm_with_file.Id
	namecontainer = Perm_with_file.Collectionname
	typeblob = Perm_with_file.Containertype
	fmt.Print(Perm_with_file.Collectionname)

}

func uploadFile(w http.ResponseWriter, r *http.Request) {

	url := "https://paktradegallery.blob.core.windows.net/" //replace <StorageAccountName> with your Azure storage account name
	//ctx := context.Background()
	credential, err1 := azidentity.NewDefaultAzureCredential(nil)
	if err1 != nil {
		log.Fatal("Invalid credentials with error: " + err1.Error())
	}

	blobname1 := namecontainer + "/" + typeblob + "/" + obid

	fmt.Printf("Creating a container name %s\n", namecontainer)

	err4 := r.ParseMultipartForm(200000) // grab the multipart form
	if err4 != nil {
		fmt.Fprintln(w, err4)
		return
	}
	// Perm_with_file.Fp
	formdata := r.MultipartForm      // ok, no problem so far, read the Form data
	files := formdata.File["myFile"] // grab the filenames
	for i, _ := range files {        // loop through the files one by one

		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		out, err := os.Create("/tmp/" + files[i].Filename)

		fileBytes, err3 := ioutil.ReadAll(file)
		if err3 != nil {
			fmt.Println("Error reading the File")

			log.Fatal(err3)
		}

		out.Write(fileBytes)

		client, err := azblob.NewClient(url, credential, nil)
		handleError(err)

		// Upload the file to a block blob
		_, err = client.UploadFile(context.TODO(), blobname1, files[i].Filename, out,
			&azblob.UploadFileOptions{
				BlockSize:   int64(1024),
				Concurrency: uint16(3),
				// If Progress is non-nil, this function is called periodically as bytes are uploaded.
				Progress: func(bytesTransferred int64) {
					fmt.Println(bytesTransferred)
				},
			})
		handleError(err)

		defer os.Remove(out.Name())
	}
	return
}
