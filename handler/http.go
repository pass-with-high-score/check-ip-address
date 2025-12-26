package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"checkip/geo"
	"checkip/netutil"
)

func IPHandler(w http.ResponseWriter, r *http.Request) {
	ip := netutil.GetClientIP(r)
	info := geo.LookupIP(ip)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(info); err != nil {
		log.Println("encode error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		http.Error(w, "Missing ip query parameter", http.StatusBadRequest)
		return
	}

	info := geo.LookupIP(ip)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(info); err != nil {
		log.Println("encode error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
