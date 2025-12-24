package geo

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/oschwald/geoip2-golang"
)

var (
	CityDB *geoip2.Reader
	ISPDB  *geoip2.Reader
	ASNDB  *geoip2.Reader
)

func LoadDB() error {
	var err error

	if CityDB != nil {
		CityDB.Close()
	}
	if ISPDB != nil {
		ISPDB.Close()
	}
	if ASNDB != nil {
		ASNDB.Close()
	}

	CityDB, _ = geoip2.Open("/var/lib/GeoIP/GeoLite2-City.mmdb")
	ISPDB, _ = geoip2.Open("/var/lib/GeoIP/GeoLite2-ISP.mmdb")
	ASNDB, _ = geoip2.Open("/var/lib/GeoIP/GeoLite2-ASN.mmdb")

	if CityDB == nil && ISPDB == nil && ASNDB == nil {
		return err
	}

	log.Println("GeoIP databases loaded")
	return nil
}

func HandleReloadSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)

	go func() {
		for range ch {
			log.Println("Reloading GeoIP databases...")
			if err := LoadDB(); err != nil {
				log.Println("Reload failed:", err)
			}
		}
	}()
}

func CloseDB() {
	if CityDB != nil {
		CityDB.Close()
	}
	if ISPDB != nil {
		ISPDB.Close()
	}
	if ASNDB != nil {
		ASNDB.Close()
	}
}
