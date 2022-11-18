package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

type input_type1 struct {
	Id             string                  `form:"id"`
	Collectionname string                  `form:"collectionname"`
	Containertype  string                  `form:"containertype"`
	Fp             []*multipart.FileHeader `form:"fp" binding:"required"`
}

func main() {
	// setupRoutes()
	r := gin.Default()
	url := "https://paktradegallery.blob.core.windows.net/" //replace <StorageAccountName> with your Azure storage account name
	//ctx := context.Background()
	credential, err1 := azidentity.NewDefaultAzureCredential(nil)
	if err1 != nil {
		log.Fatal("Invalid credentials with error: " + err1.Error())
	}

	r.POST("upload", func(c *gin.Context) {
		var filedata input_type1
		if err := c.ShouldBind(&filedata); err != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}
		blobname1 := filedata.Collectionname + "/" + filedata.Containertype + "/" + filedata.Id
		// formdata := filedata.Fb
		// abc := os.NewFile(0, "temp/abbasi.jpeg")
		// abx, _ := os.Open("")
		// temp, err2 := os.CreateTemp(filedata.Containertype, "*"+filedata.Fp[0].Filename)

		for i, _ := range filedata.Fp {
			// err := c.SaveUploadedFile(filedata.Fp[i], "temp/"+filedata.Fp[i].Filename)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// out, err := os.Create("/tmp/" + files[i].Filename)
			file, err := filedata.Fp[i].Open()
			defer file.Close()
			if err != nil {
				// fmt.Fprintln(w, err)
				return
			}

			// abx, err2 := os.Open("temp/" + filedata.Fp[i].Filename)
			// if err2 != nil {
			// 	fmt.Print(err2)
			// }
			out, err := os.Create("/tmp/" + filedata.Fp[i].Filename)
			handleError(err)

			fileBytes, err3 := ioutil.ReadAll(file)
			if err3 != nil {
				fmt.Println("Error reading the File")

				log.Fatal(err3)
			}
			out.Write(fileBytes)

			client, err := azblob.NewClient(url, credential, nil)
			handleError(err)

			// abc := multipart.File(filedata.Fp)
			// abc = filedata
			// Upload the file to a block blob

			_, err = client.UploadFile(context.TODO(), blobname1, filedata.Fp[i].Filename, out,
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
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data":   filedata,
		})
	})
	r.Run(":8080")

}
