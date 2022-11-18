package storage

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

const MAX_UPLOAD_SIZE = 20048576 // 1MB 10048576

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

/*func UploadFile(w http.ResponseWriter, r *http.Request) {

	if r.ContentLength > MAX_UPLOAD_SIZE {
		http.Error(w, "The uploaded image is too big. Please use an image less than 1MB in size", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	fmt.Println("WHOLE FILE", file)
	fileName := r.FormValue("file_name")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	fmt.Println("FILE: ", f.Name())
	defer f.Close()
	_, _ = io.WriteString(w, "File "+fileName+" Uploaded successfully")
	_, _ = io.Copy(f, file)
}*/

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
	//[{},{},{}]
	var files_array []*os.File
	for i, fileHeader := range files {
		fmt.Println("FILE: ", files[i].Filename)
		// upload each file to azure blob storage.
		//
		// Restrict the size of each uploaded file to 1MB.
		// To prevent the aggregate size from exceeding
		// a specified value, use the http.MaxBytesReader() method
		// before calling ParseMultipartForm()
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

		//uploadToAzureBlob(files_array)
		// Open the file
		/*file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}*/
		//file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = os.MkdirAll("./uploads", os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	fmt.Fprintf(w, "Upload successful")
}

func uploadToAzureBlob(file []*os.File) {

	url := "https://paktradegallery.blob.core.windows.net/"
	credential, err1 := azidentity.NewDefaultAzureCredential(nil)
	if err1 != nil {
		log.Fatal("Invalid credentials with error: " + err1.Error())
	}
	blobname1 := "gallerycontainer" + "/" + "Test001" + "/" + "Test003"
	client, client_error := azblob.NewClient(url, credential, nil)
	handleError(client_error)
	for i, f := range file {
		fmt.Println("Called")
		_, _err := client.UploadFile(context.TODO(), blobname1, f.Name(), file[i],
			&azblob.UploadFileOptions{
				BlockSize:   int64(1024),
				Concurrency: uint16(3),
				// If Progress is non-nil, this function is called periodically as bytes are uploaded.
				Progress: func(bytesTransferred int64) {
					fmt.Println(bytesTransferred)
				},
			})
		handleError(_err)
	}

}

/*
	func UploadFile() {
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
*/
