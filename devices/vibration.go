package devices

// Vibration define xiaomi vibration
type Vibration struct {
	SID            string `json:"sid"`
	Status         string `json:"status"`
	BedActivity    string `json:"bed_activity"`
	Coordination   string `json:"coordination"`
	FinalTiltAngle string `json:"final_tilt_angle"`
	Voltage        int    `json:"voltage"`
}
