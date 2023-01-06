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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	Id      string `json:"item_id"`
}

type blob_path_struct struct {
	Blobpath string `json:"blobpath"`
}

// call function to login on azure blob
var client = docking.AzureBloblogs()
var insetedid *mongo.InsertOneResult

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
	id1 := r.MultipartForm.Value["mamal_id"]
	typee1 := r.MultipartForm.Value["entity_type"]
	subtype1 := r.MultipartForm.Value["entity_sub_type"]
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
		//	fmt.Print(fileBytes)
		handleError(err)
		files_array = append(files_array, new_file)
		if fileHeader.Size > MAX_UPLOAD_SIZE {
			http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return
		}
		defer os.Remove(new_file.Name())
	}

	uploadToAzureBlob(files_array, id, typee, subtype)
	stringObjectID := insetedid.InsertedID.(primitive.ObjectID).Hex()

	mesage1 := &resp{
		Status:  http.StatusOK,
		Message: "Upload successful",
		Id:      stringObjectID,
	}
	data, err := json.Marshal(mesage1)
	handleError(err)

	fmt.Fprintf(w, string(data))

}

func uploadToAzureBlob(file []*os.File, id string, type_ string, subtype string) {
	var file_path []string

	for i, f := range file {
		blobname1 := "media" + "/" + id + "/" + type_ + "/" + subtype

		// fmt.Println("file upload to this path ", blobname1, time.Now())
		// var blob = client.;

		_, _err := client.UploadFile(context.TODO(), blobname1, f.Name(), file[i],
			&azblob.UploadFileOptions{
				HTTPHeaders: &blob.HTTPHeaders{
					BlobContentType: to.Ptr("image/jpg"),
					// BlobContentDisposition: to.Ptr("attachment"),
				},

				BlockSize:   int64(1024),
				Concurrency: uint16(3),
				// If Progress is non-nil, this function is called periodically as bytes are uploaded.
				Progress: func(bytesTransferred int64) {
					fmt.Println(bytesTransferred)
				},
			})
		handleError(_err)
		////// update path to
		file_path = append(file_path, "https://paktradestorage.blob.core.windows.net/"+blobname1+f.Name())
	}
	// fmt.Println(file_path)
	result := Add_data_to_mongo(file_path)
	insetedid = result
}

//////// delete image form azure blob now for this time code not use
/*
func Deltefile(w http.ResponseWriter, r *http.Request) {

		azblob.ParseURL("")

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var path blob_path_struct
		err := json.NewDecoder(r.Body).Decode(&path)
		handleError(err)
		// for _, a := range path.Blobpath {
		_, err1 := client.DeleteBlob(context.TODO(), "gallerycontainer", path.Blobpath, nil)

		if err1 != nil {
			mesage1 := &resp{
				Status:  http.StatusBadGateway,
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
			data, err := json.Marshal(mesage1.Message)
			handleError(err)

			fmt.Fprintf(w, string(data))
		}
		// }
	}
*/
func Add_data_to_mongo(image_array []string) *mongo.InsertOneResult {
	dumy_array := [1]string{"no image"}

	insertdat := bson.M{"name": bson.M{
		"en": "",
		"ar": "",
	},
		"feature": bson.A{
			"",
		},
		"available_color": bson.A{
			"",
		},
		"size": bson.M{
			"available_size": bson.A{
				"",
				"",
			},
			"size_chart": "",
		},
		"images": bson.A{
			bson.M{
				"low_quility": image_array,
			},
			bson.M{
				"high_quility": dumy_array,
			},
		},
		"price":        0,
		"status":       "pending",
		"gender":       "",
		"category":     "",
		"sub-category": "",
	}

	//fmt.Print(body)
	coll := docking.PakTradeDb.Collection("cloths")

	// // // insert a user

	responceid, err3 := coll.InsertOne(context.TODO(), insertdat)
	if err3 != nil {
		fmt.Print(err3)
	}
	return responceid
}
