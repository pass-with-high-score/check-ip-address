package main

import (
	"log"
	"net/http"

	"checkip/geo"
	"checkip/handler"
)

func main() {
	if err := geo.LoadDB(); err != nil {
		log.Fatal("Failed to load GeoIP DB:", err)
	}
	defer geo.CloseDB()

	geo.HandleReloadSignal()

	http.HandleFunc("/", handler.IPHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
