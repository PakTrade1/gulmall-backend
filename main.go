package main

import (
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"
	Allcart "pak-trade-go/api/cart"
	color "pak-trade-go/api/color"
	item "pak-trade-go/api/items"
	Alluser "pak-trade-go/api/mammals"
	Insertuser "pak-trade-go/api/mammals"

	"github.com/gorilla/mux"
)

func main() {
	docking.PakTradeConnection()
	r := mux.NewRouter()
	http.Handle("/", r)
	r.HandleFunc("/getColor", color.Color)
	r.HandleFunc("/getItem", item.Items)
	r.HandleFunc("/getUser", Alluser.Mammals_getall)
	r.HandleFunc("/addUser", Insertuser.Mammals_insertone)
	r.HandleFunc("/getCart", Allcart.Cart_getall)
	r.HandleFunc("/addCart", Allcart.Mammals_insertone)

	fmt.Println("runging server port 9900")
	http.ListenAndServe(":9901", nil)
	//fmt.Println("Runging server port 80")
	//http.ListenAndServe(":80", nil)

}
