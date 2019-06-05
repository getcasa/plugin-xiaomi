package devices

// Switch define xiaomi switch button
type Switch struct {
	SID     string `json:"sid"`
	Status  string `json:"status"`
	Voltage int    `json:"voltage"`
}
