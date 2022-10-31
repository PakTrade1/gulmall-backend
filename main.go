package main

import (
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"
	color "pak-trade-go/api/color"
	item "pak-trade-go/api/items"

	"github.com/gorilla/mux"
)

func main() {
	docking.PakTradeConnection()
	r := mux.NewRouter()
	http.Handle("/", r)
	r.HandleFunc("/color", color.Color)
	r.HandleFunc("/item", item.Items)
	fmt.Println("Runging server port 80")
	http.ListenAndServe(":80", nil)

}
