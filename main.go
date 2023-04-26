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

	shipping_addres "pak-trade-go/api/address"
	item "pak-trade-go/api/items"
	User "pak-trade-go/api/mammals"
	payment_service "pak-trade-go/api/payment"
	size "pak-trade-go/api/size"
	weight "pak-trade-go/api/weight"
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
	// ADD ROUTES

	r.HandleFunc("/add-cart", Allcart.Update_cart)
	r.HandleFunc("/add-size", size.Add_size)
	r.HandleFunc("/add-address", shipping_addres.Add_shipping_address)
	r.HandleFunc("/add-category", categories.Add_category)
	r.HandleFunc("/add-sub-category", categories.Add_sub_category)
	r.HandleFunc("/add-sub-category-child", categories.Add_sub_child_category)

	// UPDATE ROUTE
	r.HandleFunc("/update-category", categories.Update_Category)
	r.HandleFunc("/update-user", User.Mammals_update_one)
	r.HandleFunc("/update-cart-in", Allcart.Update_cart)
	r.HandleFunc("/update-item", item.Item_update_one)
	// DELETE ROUTE
	r.HandleFunc("/delete-address", shipping_addres.Delete_shipping_address)
	r.HandleFunc("/delete-category", categories.Delete_category)
	r.HandleFunc("/delete-cart", Allcart.Cart_delete)
	r.HandleFunc("/delete-item-by-id", item.Item_delete_by_id)
	// GET ROUTE
	r.HandleFunc("/get-color", color.Color)
	r.HandleFunc("/get-size-chart", size.Size_select_by_child_id)
	r.HandleFunc("/get-payment-service", payment_service.Get_payment_method)
	r.HandleFunc("/get-item-with-status", item.Items)
	r.HandleFunc("/get-categories", categories.Get_all_categories)
	r.HandleFunc("/get-user", User.Mammals_getall)
	r.HandleFunc("/get-gender", gender.Gender)
	r.HandleFunc("/get-user-by-id", User.Mammals_select_one)
	r.HandleFunc("/get-item-by-id", item.Serch_item_by_id)
	r.HandleFunc("/get-cart-with-id", Allcart.Get_cart_with_id)
	r.HandleFunc("/get-cart", Allcart.Cart_getall)
	r.HandleFunc("/get-order-details", Allcart.Order_with_need_data)
	r.HandleFunc("/get-all-cart", Allcart.Get_cart_all_with_id_data)
	r.HandleFunc("/get-all-item", item.Get_all_items)
	r.HandleFunc("/get-weight", weight.Weight)
	r.HandleFunc("/get-sub-categories", categories.Sub_Categories_select_by_Cat_id)
	r.HandleFunc("/get-child-categories", categories.Child_Categories_select_by__sub_Cat_id)
	r.HandleFunc("/get-address", shipping_addres.Get_shipping_address_with_mammal_id)

	// UPLOAD FILE
	r.HandleFunc("/upload-file", storage.UploadFile).Methods("POST")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)
	fmt.Println("Runging server port ===> 80")
	http.ListenAndServe(":80", handler)
}
