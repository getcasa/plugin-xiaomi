package devices

// Gateway define xiaomi gateway
type Gateway struct {
	SID          string   `json:"sid"`
	IP           string   `json:"ip"`
	Token        string   `json:"token"`
	Devices      []string `json:"data"`
	RGB          int      `json:"rgb"`
	Illumination int      `json:"illumination"`
}
