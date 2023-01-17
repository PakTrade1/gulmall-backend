package main

import (
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"
	Allcart "pak-trade-go/api/cart"
	categories "pak-trade-go/api/categories"
	color "pak-trade-go/api/color"
	gender "pak-trade-go/api/gender"

	//	"pak-trade-go/api/weight"

	item "pak-trade-go/api/items"
	User "pak-trade-go/api/mammals"
	payment_service "pak-trade-go/api/payment"
	size "pak-trade-go/api/size"
	weight "pak-trade-go/api/weight"
	blobstorage "pak-trade-go/blobstorage"
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
	// RETURNS LIST OF ALL PRE-DEFIEND COLORS.
	r.HandleFunc("/get-color", color.Color)
	// GET SIZE CHART RETURNS ALL
	r.HandleFunc("/get-size-chart", size.Size_select_by_child_id)
	r.HandleFunc("/get-payment-service", payment_service.Get_payment_method)
	r.HandleFunc("/update-item", item.Item_update_one)
	r.HandleFunc("/get-item-with-status", item.Items)
	// RETURNS LIST OF CATEGORIES.
	r.HandleFunc("/get-categories", categories.Get_all_categories)
	// ADD _____ Update___ Delete_______
	r.HandleFunc("/add-category", categories.Add_category)
	r.HandleFunc("/update-category", categories.Update_Category)
	r.HandleFunc("/delete-category", categories.Delete_category)
	// ADD _________ Size
	r.HandleFunc("/add-size", size.Add_size)
	// Get  _________ Weightt
	r.HandleFunc("/get-weight", weight.Weight)

	// RETURNS LIST OF ALL SUB-CAT.
	r.HandleFunc("/get-sub-categories", categories.Sub_Categories_select_by_Cat_id)
	// RETURNS LIST OF ALL SUB-CAT-CHILD.
	r.HandleFunc("/get-child-categories", categories.Child_Categories_select_by__sub_Cat_id)
	// RETURNS LIST OF ALL PRE-DEFIEND GENDERS.
	r.HandleFunc("/get-gender", gender.Gender)
	// RETURNS LIST OF ALL USERS.
	r.HandleFunc("/get-user", User.Mammals_getall)
	r.HandleFunc("/add-user", User.Mammals_insertone)
	r.HandleFunc("/get-user-by-id", User.Mammals_select_one)
	r.HandleFunc("/update-user", User.Mammals_update_one)
	// RETURNS SINGLE ITEM.
	r.HandleFunc("/get-item-by-id", item.Serch_item_by_id)     //item_id
	r.HandleFunc("/delete-item-by-id", item.Item_delete_by_id) // item_id and status
	r.HandleFunc("/get-all-item", item.Get_all_items)          // POST         // get all items
	r.HandleFunc("/upload-file", storage.UploadFile).Methods("POST")
	r.HandleFunc("/get-all-cart", Allcart.Get_cart_all_with_id_data)
	r.HandleFunc("/get-cart-with-id", Allcart.Get_cart_with_id)
	r.HandleFunc("/delete-cart", Allcart.Cart_delete)
	r.HandleFunc("/update-cart-in", Allcart.Update_cart)
	r.HandleFunc("/upload-filen", blobstorage.UploadFile).Methods("POST")
	r.HandleFunc("/get-cart", Allcart.Cart_getall)
	r.HandleFunc("/add-cart", Allcart.Update_cart)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)
	fmt.Println("Runging server port ===> 80")
	http.ListenAndServe(":80", handler)
}
