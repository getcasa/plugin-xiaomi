package devices

// SensorMotionAQ2 define xiaomi sensor motion aq2
type SensorMotionAQ2 struct {
	SID      string `json:"sid"`
	Status   string `json:"status"`
	NoMotion string `json:"no_motion"`
	Lux      string `json:"lux"`
	Voltage  int    `json:"voltage"`
}
