package model

type ASN struct {
	Number uint   `json:"number,omitempty"`
	Org    string `json:"org,omitempty"`
	ISP    string `json:"isp,omitempty"`
}
