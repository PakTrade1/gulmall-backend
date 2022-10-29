package main

import (
	"fmt"
	"net/http"
	docking "pak-trade-go/Docking"
	color "pak-trade-go/api/color"
)

func main() {
	docking.PakTradeConnection()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Ping...")
	})
	http.HandleFunc("/color", color.Color)
	fmt.Println("runging server port 8080")
	http.ListenAndServe(":80", nil)
}
