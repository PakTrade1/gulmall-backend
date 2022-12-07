package main

import (
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"
	Allcart "pak-trade-go/api/cart"
	color "pak-trade-go/api/color"
	item "pak-trade-go/api/items"
	User "pak-trade-go/api/mammals"
	payment_service "pak-trade-go/api/payment"
	storage "pak-trade-go/storage"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	r.HandleFunc("/get-payment-service", payment_service.Get_payment_method)

	// r.HandleFunc("/ItemAdd", item.ItemInsertone)
	r.HandleFunc("/update-item", item.Item_update_one)
	r.HandleFunc("/get-item-with-status", item.Items)
	r.HandleFunc("/get-user", User.Mammals_getall)
	r.HandleFunc("/add-user", User.Mammals_insertone)
	r.HandleFunc("/get-user-by-id", User.Mammals_select_one)
	r.HandleFunc("/update-user", User.Mammals_update_one)
	r.HandleFunc("/get-item-by-id", item.Serch_item_by_id)     //item_id
	r.HandleFunc("/delete-item-by-id", item.Item_delete_by_id) // item_id and status
	r.HandleFunc("/get-all-item", item.Get_all_items)          // POST         // get all items
	r.HandleFunc("/upload-file", storage.UploadFile).Methods("POST")
	//r.HandleFunc("/delete-file", item.Delte_item).Methods("POST")
	//r.HandleFunc("/delete-image", storage.Deltefile).Methods("POST")

	r.HandleFunc("/get-all-cart", Allcart.Get_cart_all_with_id_data)
	r.HandleFunc("/get-cart-with-id", Allcart.Get_cart_with_id)

	r.HandleFunc("/get-cart", Allcart.Cart_getall)
	r.HandleFunc("/add-cart", Allcart.Cart_insertone)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)
	fmt.Println("runging server port 9900")
	http.ListenAndServe(":80", handler)
	//fmt.Println("Runging server port 80")
	//http.ListenAndServe(":80", nil)

}
