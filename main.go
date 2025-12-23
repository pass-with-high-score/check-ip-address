package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/oschwald/geoip2-golang"
)

var (
	cityDB *geoip2.Reader
	ispDB  *geoip2.Reader
	asnDB  *geoip2.Reader
)

/* =========================
   IP DETECTION
========================= */

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

/* =========================
   GEO LOOKUP
========================= */

type IPInfo struct {
	IP      string `json:"ip"`
	ISP     string `json:"isp,omitempty"`
	ASN     uint   `json:"asn,omitempty"`
	ASNOrg  string `json:"asn_org,omitempty"`
	City    string `json:"city,omitempty"`
	Region  string `json:"region,omitempty"`
	Country string `json:"country,omitempty"`
}

func lookupIP(ipStr string) IPInfo {
	result := IPInfo{IP: ipStr}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return result
	}

	// Lookup City data
	if cityDB != nil {
		if city, err := cityDB.City(ip); err == nil {
			result.City = city.City.Names["en"]
			result.Country = city.Country.Names["en"]

			if len(city.Subdivisions) > 0 {
				result.Region = city.Subdivisions[0].Names["en"]
			}
		}
	}

	// Lookup ISP data
	if ispDB != nil {
		if isp, err := ispDB.ISP(ip); err == nil {
			result.ISP = isp.ISP
		}
	}

	// Lookup ASN data
	if asnDB != nil {
		if asn, err := asnDB.ASN(ip); err == nil {
			result.ASN = asn.AutonomousSystemNumber
			result.ASNOrg = asn.AutonomousSystemOrganization
		}
	}

	return result
}

/* =========================
   DB LOAD / RELOAD
========================= */

func loadGeoDB() error {
	var err error

	// Close existing databases
	if cityDB != nil {
		cityDB.Close()
	}
	if ispDB != nil {
		ispDB.Close()
	}
	if asnDB != nil {
		asnDB.Close()
	}

	// Load City database
	cityDB, err = geoip2.Open("/var/lib/GeoIP/GeoLite2-City.mmdb")
	if err != nil {
		log.Printf("Warning: Could not load City DB: %v", err)
	}

	// Load ISP database
	ispDB, err = geoip2.Open("/var/lib/GeoIP/GeoLite2-ISP.mmdb")
	if err != nil {
		log.Printf("Warning: Could not load ISP DB: %v", err)
	}

	// Load ASN database
	asnDB, err = geoip2.Open("/var/lib/GeoIP/GeoLite2-ASN.mmdb")
	if err != nil {
		log.Printf("Warning: Could not load ASN DB: %v", err)
	}

	// Check if at least one database loaded successfully
	if cityDB == nil && ispDB == nil && asnDB == nil {
		return err
	}

	log.Println("GeoIP databases loaded successfully")
	return nil
}

func handleReloadSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)

	go func() {
		for range ch {
			log.Println("Reloading GeoIP databases...")
			if err := loadGeoDB(); err != nil {
				log.Println("Reload failed:", err)
			} else {
				log.Println("Reload completed successfully")
			}
		}
	}()
}

/* =========================
   HTTP HANDLER
========================= */

func handler(w http.ResponseWriter, r *http.Request) {
	ip := getClientIP(r)
	info := lookupIP(ip)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(info); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

/* =========================
   MAIN
========================= */

func main() {
	if err := loadGeoDB(); err != nil {
		log.Fatal("Failed to load GeoIP databases:", err)
	}

	// Ensure cleanup on exit
	defer func() {
		if cityDB != nil {
			cityDB.Close()
		}
		if ispDB != nil {
			ispDB.Close()
		}
		if asnDB != nil {
			asnDB.Close()
		}
	}()

	handleReloadSignal()

	http.HandleFunc("/", handler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
