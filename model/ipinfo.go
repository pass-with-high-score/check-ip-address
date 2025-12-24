package model

type IPInfo struct {
	IP      string `json:"ip"`
	ISP     string `json:"isp,omitempty"`
	ASN     uint   `json:"asn,omitempty"`
	ASNOrg  string `json:"asn_org,omitempty"`
	City    string `json:"city,omitempty"`
	Region  string `json:"region,omitempty"`
	Country string `json:"country,omitempty"`
}
