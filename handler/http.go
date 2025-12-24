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
