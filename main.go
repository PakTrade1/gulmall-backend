package main

import (
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"
	ads "pak-trade-go/api/ads"
	categories "pak-trade-go/api/categories"
	color "pak-trade-go/api/color"
	clothingFilter "pak-trade-go/api/filter"
	gender "pak-trade-go/api/gender"
	"pak-trade-go/api/geolocation"
	keyword "pak-trade-go/api/serchKeyWord"
	"pak-trade-go/api/signin"
	tier "pak-trade-go/api/tier"

	//	"pak-trade-go/api/weight"
	shipping_addres "pak-trade-go/api/address"
	authWhatsapp "pak-trade-go/api/auth"
	cart "pak-trade-go/api/cart"
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

	// r.HandleFunc("/add-keyword", keyword.Serchkeywordinsert)
	r.HandleFunc("/add-cart", cart.AddToCartHandler(docking.PakTradeDb.Collection("cart_mammals")))

	r.HandleFunc("/add-size", size.Add_size)
	r.HandleFunc("/add-address", shipping_addres.Add_shipping_address)
	r.HandleFunc("/add-category", categories.Add_category)
	r.HandleFunc("/add-sub-category", categories.Add_sub_category)
	r.HandleFunc("/add-sub-category-child", categories.Add_sub_child_category)
	r.HandleFunc("/add-mammals_registration", User.Mammals_user_registration)
	r.HandleFunc("/mammals_registration", User.Mammals_user_registration)
	r.HandleFunc("/check-email", User.CheckEmailHandler)
	r.HandleFunc("/check-email-verified", User.CheckEmailVerifiedHandler)
	r.HandleFunc("/check-phone-number", User.CheckPhoneHandler)
	r.HandleFunc("/signin-email", signin.SignInEmailHandler)
	r.HandleFunc("/signin-phone", signin.SignInPhoneHandler)
	r.HandleFunc("/add-item", item.Update_item_wrt_category)
	r.HandleFunc("/create-draft-item", item.Add_item_image)

	//	r.HandleFunc("/add-mammals_registration", User.Mammals_user_registration)

	// UPDATE ROUTE
	r.HandleFunc("/update-category", categories.Update_Category)
	r.HandleFunc("/update-user", User.Mammals_update_one)
	// r.HandleFunc("/update-cart-in", Allcart.Update_cart)
	// r.HandleFunc("/update-item", item.Item_update_one)
	r.HandleFunc("/update-item", item.Add_item_update)
	r.HandleFunc("/update-address", shipping_addres.Address_update_one)

	// DELETE ROUTE
	r.HandleFunc("/delete-address", shipping_addres.Delete_shipping_address)
	r.HandleFunc("/delete-category", categories.Delete_category)
	// r.HandleFunc("/delete-cart", Allcart.Cart_delete)
	r.HandleFunc("/delete-item-by-id", item.Item_delete_by_id)
	// GET ROUTE
	r.HandleFunc("/get-color", color.Color)
	r.HandleFunc("/get-tier", tier.Tier_get)

	r.HandleFunc("/get-plans", tier.Plan_select)
	r.HandleFunc("/get-item-by-keyword", item.Serch_item_by_keyword)

	r.HandleFunc("/get-size-chart", size.Size_select_by_child_id)
	r.HandleFunc("/get-payment-service", payment_service.Get_payment_method)
	r.HandleFunc("/get-item-with-status", item.Items)
	r.HandleFunc("/get-categories", categories.Get_all_categories)
	r.HandleFunc("/get-user", User.Mammals_getall)
	r.HandleFunc("/get-gender", gender.Gender)
	r.HandleFunc("/get-user-by-id", User.Mammals_select_one)
	r.HandleFunc("/get-item-by-id", item.Serch_item_by_id)
	// r.HandleFunc("/get-cart-with-id", Allcart.Get_cart_with_id)

	// r.HandleFunc("/get-order-details", Allcart.Order_with_need_data)
	// r.HandleFunc("/get-all-cart", Allcart.Get_cart_all_with_id_data)
	r.HandleFunc("/get-all-item", item.Get_all_items)
	r.HandleFunc("/get-weight", weight.Weight)
	r.HandleFunc("/get-sub-categories", categories.Sub_Categories_select_by_Cat_id)
	r.HandleFunc("/get-child-categories", categories.Child_Categories_select_by__sub_Cat_id)
	r.HandleFunc("/get-address", shipping_addres.Get_shipping_address_with_mammal_id)
	r.HandleFunc("/get-ads-by-id", ads.Get_ads_user_by_post_id)
	r.HandleFunc("/get-keyword", keyword.Get_all_items_serchkey)
	r.HandleFunc("/get-clothing-filters", clothingFilter.FiltersHandler)
	r.HandleFunc("/create-user", User.CreateUser)
	r.HandleFunc("/get-all-items-by-mamal-id", item.GetUserAndItemsHandler)
	r.HandleFunc("/send-otp", authWhatsapp.SendOTPHandler)
	r.HandleFunc("/verify-otp", authWhatsapp.VerifyOTPHandler)
	r.HandleFunc("/ip", geolocation.IPHandler)
	r.HandleFunc("/update-cart", cart.UpdateOrderHandler(docking.PakTradeDb.Collection("cart_mammals"), docking.PakTradeDb.Collection("cart_audits")))
	// r.HandleFunc("/get-ads-by-id/", ads.Get_ads_user_by_post_id)

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
