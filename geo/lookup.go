package geo

import (
	"net"

	"checkip/model"
)

func LookupIP(ipStr string) model.IPInfo {
	result := model.IPInfo{IP: ipStr}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return result
	}

	if CityDB != nil {
		if city, err := CityDB.City(ip); err == nil {
			result.City = city.City.Names["en"]
			result.Country = city.Country.Names["en"]
			if len(city.Subdivisions) > 0 {
				result.Region = city.Subdivisions[0].Names["en"]
			}
		}
	}

	if ISPDB != nil {
		if isp, err := ISPDB.ISP(ip); err == nil {
			result.ISP = isp.ISP
		}
	}

	if ASNDB != nil {
		if asn, err := ASNDB.ASN(ip); err == nil {
			result.ASN = asn.AutonomousSystemNumber
			result.ASNOrg = asn.AutonomousSystemOrganization
		}
	}

	return result
}
