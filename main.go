package main

import (
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"
	color "pak-trade-go/api"
)

func main() {

	docking.PakTradeConnection()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})
	http.HandleFunc("/color", color.Color)
	fmt.Println("runging server port 9900")
	http.ListenAndServe(":9900", nil)
}
