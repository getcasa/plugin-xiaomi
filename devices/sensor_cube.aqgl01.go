package devices

// SensorCubeAqgl01 define xiaomi sensor cube aqgl01
type SensorCubeAqgl01 struct {
	SID     string `json:"sid"`
	Status  string `json:"status"`
	Rotate  string `json:"rotate"`
	Voltage int    `json:"voltage"`
}
