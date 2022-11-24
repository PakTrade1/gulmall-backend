package main

import (
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"
	Allcart "pak-trade-go/api/cart"
	color "pak-trade-go/api/color"
	item "pak-trade-go/api/items"
	User "pak-trade-go/api/mammals"
	storage "pak-trade-go/storage"

	"github.com/gorilla/mux"
)

func main() {
	// DOCKING WITH AZURE BLOB STORAGE.
	docking.PakTradeConnection()
	docking.AzureBloblogs()
	// ROUTERS
	r := mux.NewRouter()
	http.Handle("/", r)
	// API ENDPOINTS
	r.HandleFunc("/get-color", color.Color)
	// r.HandleFunc("/ItemAdd", item.ItemInsertone)
	r.HandleFunc("/update-item", item.Item_update_one)
	r.HandleFunc("/get-item", item.Items)
	r.HandleFunc("/get-user", User.Mammals_getall)
	r.HandleFunc("/add-user", User.Mammals_insertone)
	r.HandleFunc("/get-user-by-id", User.Mammals_select_one)
	r.HandleFunc("/update-user", User.Mammals_update_one)
	r.HandleFunc("/upload-file", storage.UploadFile).Methods("POST")
	r.HandleFunc("/delete-file", storage.Deltefile).Methods("POST")
	r.HandleFunc("/get-cart", Allcart.Cart_getall)
	r.HandleFunc("/add-cart", Allcart.Cart_insertone)

	fmt.Println("runging server port 9900")
	http.ListenAndServe(":9900", nil)
	//fmt.Println("Runging server port 80")
	//http.ListenAndServe(":80", nil)

}
