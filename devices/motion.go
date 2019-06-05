package devices

// Motion define xiaomi motion sensor
type Motion struct {
	SID      string `json:"sid"`
	Status   string `json:"status"`
	NoMotion string `json:"no_motion"`
	Voltage  int    `json:"voltage"`
}
