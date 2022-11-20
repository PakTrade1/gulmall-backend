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
	docking.PakTradeConnection()
	docking.AzureBloblogs()
	r := mux.NewRouter()
	http.Handle("/", r)
	r.HandleFunc("/getColor", color.Color)
	r.HandleFunc("/getItem", item.Items)
	r.HandleFunc("/getUser", User.Mammals_getall)
	r.HandleFunc("/addUser", User.Mammals_insertone)
	r.HandleFunc("/searchUser", User.Mammals_select_one)
	r.HandleFunc("/updateUser", User.Mammals_update_one)
	r.HandleFunc("/upload-file", storage.UploadFile).Methods("POST")
	r.HandleFunc("/delete-file", storage.Deltefile).Methods("POST")

	r.HandleFunc("/getCart", Allcart.Cart_getall)
	r.HandleFunc("/addCart", Allcart.Cart_insertone)

	fmt.Println("runging server port 9900")
	http.ListenAndServe(":9900", nil)
	//fmt.Println("Runging server port 80")
	//http.ListenAndServe(":80", nil)

}
