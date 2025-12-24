package geo

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
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

	dbDir := os.Getenv("GEOIP_DB_DIR")
	if dbDir == "" {
		dbDir = "./db"
	}

	CityDB, _ = geoip2.Open(filepath.Join(dbDir, "GeoLite2-City.mmdb"))
	ISPDB, _ = geoip2.Open(filepath.Join(dbDir, "GeoLite2-ISP.mmdb"))
	ASNDB, _ = geoip2.Open(filepath.Join(dbDir, "GeoLite2-ASN.mmdb"))

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
