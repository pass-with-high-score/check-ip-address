package model

type IPInfo struct {
	IP       IP       `json:"ip"`
	ASN      ASN      `json:"asn,omitempty"`
	Location Location `json:"location,omitempty"`
}
