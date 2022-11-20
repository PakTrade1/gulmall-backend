package storage

//USER_ID==> Container name/
//MEDIA_TYPE==> Type of media which could be Item,Personal,Doc
//MEDIA_SUB_TYPE ==> Type of media that could be dress, shoes, profile,
//123/item/dress/abc.jpeg

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	docking "pak-trade-go/Docking"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

const MAX_UPLOAD_SIZE = 30048576 // 1MB 10048576

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

type resp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type blob_path_struct struct {
	Blobpath string `json:"blobpath"`
}

// call function to login on azure blob
var client = docking.AzureBloblogs()

func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File["file"]
	id1 := r.MultipartForm.Value["id"]
	typee1 := r.MultipartForm.Value["type"]
	subtype1 := r.MultipartForm.Value["subtype"]
	id := strings.Join(id1, " ")
	typee := strings.Join(typee1, " ")
	subtype := strings.Join(subtype1, " ")
	var files_array []*os.File
	for i, fileHeader := range files {
		file, err := files[i].Open()
		new_file, err := os.Create(files[i].Filename)
		fileBytes, err3 := ioutil.ReadAll(file)
		if err3 != nil {
			fmt.Println("Error reading the File")

			log.Fatal(err3)
		}
		if err3 != nil {
			fmt.Println("Error reading the File")

			log.Fatal(err3)
		}
		new_file.Write(fileBytes)

		handleError(err)
		files_array = append(files_array, new_file)
		if fileHeader.Size > MAX_UPLOAD_SIZE {
			http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return
		}
		defer os.Remove(new_file.Name())
	}

	uploadToAzureBlob(files_array, id, typee, subtype)
	mesage1 := &resp{
		Status:  http.StatusOK,
		Message: "Upload successful",
	}
	data, err := json.Marshal(mesage1)
	handleError(err)

	fmt.Fprintf(w, string(data))

}

func uploadToAzureBlob(file []*os.File, id string, type_ string, subtype string) {
	for i, f := range file {
		blobname1 := "gallerycontainer" + "/" + id + "/" + type_ + "/" + subtype
		fmt.Println("file upload to this path ", blobname1, time.Now())
		_, _err := client.UploadFile(context.TODO(), blobname1, f.Name(), file[i],
			&azblob.UploadFileOptions{
				BlockSize:   int64(1024),
				Concurrency: uint16(3),
				// If Progress is non-nil, this function is called periodically as bytes are uploaded.
				Progress: func(bytesTransferred int64) {
					// fmt.Println(bytesTransferred)
				},
			})
		handleError(_err)
	}

}

func Deltefile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var path blob_path_struct
	err := json.NewDecoder(r.Body).Decode(&path)
	handleError(err)
	_, err1 := client.DeleteBlob(context.TODO(), "gallerycontainer", path.Blobpath, nil)
	if err1 != nil {
		mesage1 := &resp{
			Status:  http.StatusOK,
			Message: "image not found",
		}
		data, err := json.Marshal(mesage1)
		handleError(err)

		fmt.Fprintf(w, string(data))

	} else {
		mesage1 := &resp{
			Status:  http.StatusOK,
			Message: "delete Image successful",
		}
		data, err := json.Marshal(mesage1)
		handleError(err)

		fmt.Fprintf(w, string(data))
	}
}
