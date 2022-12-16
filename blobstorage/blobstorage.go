package blobstorage

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gofrs/uuid"
)

type resp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Id      string `json:"item_id"`
}

const MAX_UPLOAD_SIZE = 30048576 // 1MB 10048576
func UploadBlob() {
	files, errW := walkDir(".")

	if errW != nil {
		fmt.Println("Error has occured:", errW)
	} else {
		var fFiles []string
		for _, fName := range files {
			if strings.Contains(fName, "jpg") {
				fFiles = append(fFiles, fName)
			}
		}

		m := make(map[string][]byte)

		// Read file contents into memory
		for _, fName := range fFiles {
			fmt.Println("Found file:", fName)
			dat, errR := ReadFile(fName)

			if errR != nil {
				fmt.Println("Error reading file:", fName, "Error:", errR)
			} else {
				fmt.Println("Finished reading bytes for file:", fName)
				m[fName] = dat
			}
		}

		// push file contents from memory to Azure
		for _, fName := range fFiles {
			fmt.Println("Started uploading: ", fName)
			u, errU := UploadBytesToBlob(m[fName])
			if errU != nil {
				fmt.Println("Error during upload: ", errU)
			}

			fmt.Println("Finished uploading to: ", u)
			fmt.Println("==========================================================")
		}
	}
}

func walkDir(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func ReadFile(filePath string) ([]byte, error) {
	dat, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	} else {
		return dat, nil
	}
}

func UploadBytesToBlob(b []byte) (string, error) {
	fmt.Println("byte", b)
	azrKey, accountName, endPoint, container := GetAccountInfo()
	fmt.Println("ENDPOINT: ", endPoint, " container: ", container)
	u, _ := url.Parse(fmt.Sprint(endPoint, container, "/", GetBlobName()))
	credential, errC := azblob.NewSharedKeyCredential(accountName, azrKey)
	if errC != nil {
		return "Error in account name or azure key", errC
	}

	blockBlobUrl := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))

	ctx := context.Background()
	o := azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: "image/jpg",
		},
	}
	upload, errU := azblob.UploadBufferToBlockBlob(ctx, b, blockBlobUrl, o)
	fmt.Println("UPLOAD RESP", upload.Response())
	return blockBlobUrl.String(), errU
}

func GetAccountInfo() (string, string, string, string) {
	azrKey := "pIub4yUb20E0RYSDV1BoNwZQZG080Z4O3KdjSBayKBpFeonlq2DfS7HVuG/0dNobfoPqtauE6cBE+AStUkIEuw==/0dNobfoPqtauE6cBE+AStUkIEuw==_key"
	azrBlobAccountName := "appmedia"
	azrPrimaryBlobServiceEndpoint := fmt.Sprintf("https://%s.blob.core.windows.net/", azrBlobAccountName)
	azrBlobContainer := "appmediacontainer"
	return azrKey, azrBlobAccountName, azrPrimaryBlobServiceEndpoint, azrBlobContainer
}

func GetBlobName() string {
	t := time.Now()
	uuid, _ := uuid.NewV4()
	fmt.Printf("%s-%v.jpg", t.Format("20060102"), uuid)
	return fmt.Sprintf("%s-%v.jpg", t.Format("20060102"), uuid)
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Called")
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
	var files_array []*os.File
	for i, fileHeader := range files {
		file, err := files[i].Open()
		new_file, err := os.Create(files[i].Filename)
		fileBytes, err3 := ioutil.ReadAll(file)
		UploadBytesToBlob(fileBytes)
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

}
func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
