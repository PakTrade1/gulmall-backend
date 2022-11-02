package main

import (
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"
	Allcart "pak-trade-go/api/cart"
	color "pak-trade-go/api/color"
	item "pak-trade-go/api/items"
	Alluser "pak-trade-go/api/mammals"

	"github.com/gorilla/mux"
)

func main() {
	docking.PakTradeConnection()
	r := mux.NewRouter()
	http.Handle("/", r)
	r.HandleFunc("/color", color.Color)
	r.HandleFunc("/item", item.Items)
	r.HandleFunc("/alluser", Alluser.Mammals_getall)
	//r.HandleFunc("/adduser", Insertuser.Mammals_insertone)
	r.HandleFunc("/allcart", Allcart.Cart_getall)
	//r.HandleFunc("/addtocart", Allcart.Cart_insertone)

	fmt.Println("runging server port 9900")
	http.ListenAndServe(":9901", nil)
	fmt.Println("Runging server port 80")
	http.ListenAndServe(":80", nil)

}
