package main

import (
	"fmt"
	"net/http"
	docking "pak-trade-go/services/Docking"
)

func main() {
	docking.Dbconnect()

	http.HandleFunc("/color", docking.Color)
	fmt.Println("runging server port 9900")
	http.ListenAndServe(":9900", nil)
}
