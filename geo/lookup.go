package geo

import (
	"net"

	"checkip/model"
)

func LookupIP(ipStr string) model.IPInfo {
	result := model.IPInfo{
		IP: model.IP{
			Address: ipStr,
		},
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return result
	}

	// Location (City / Country / Timezone)
	if CityDB != nil {
		if city, err := CityDB.City(ip); err == nil {
			result.Location = model.Location{
				City:        city.City.Names["en"],
				Country:     city.Country.Names["en"],
				CountryCode: city.Country.IsoCode,
				Timezone:    city.Location.TimeZone,
				Latitude:    city.Location.Latitude,
				Longitude:   city.Location.Longitude,
			}

			if len(city.Subdivisions) > 0 {
				result.Location.Region = city.Subdivisions[0].Names["en"]
			}
		}
	}

	// ISP (network info)
	if ISPDB != nil {
		if isp, err := ISPDB.ISP(ip); err == nil {
			result.ASN.ISP = isp.ISP
		}
	}

	// ASN
	if ASNDB != nil {
		if asn, err := ASNDB.ASN(ip); err == nil {
			result.ASN.Number = asn.AutonomousSystemNumber
			result.ASN.Org = asn.AutonomousSystemOrganization
		}
	}

	return result
}
