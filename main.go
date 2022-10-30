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

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("Ping...")
	// })
	http.Handle("/", r)
	r.HandleFunc("/color", color.Color)
	r.HandleFunc("/item", item.Items)
	fmt.Println("runging server port 9900")
	http.ListenAndServe(":9900", nil)

}
