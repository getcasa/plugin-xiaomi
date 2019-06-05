package devices

// Gateway define xiaomi gateway
type Gateway struct {
	SID          string `json:"sid"`
	IP           string `json:"ip"`
	RGB          int    `json:"rgb"`
	Illumination int    `json:"illumination"`
}
