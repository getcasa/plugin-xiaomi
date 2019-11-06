package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/getcasa/plugin-xiaomi/devices"
	"github.com/getcasa/sdk"
)

func main() {}

const (
	ip   = "224.0.0.50"
	port = "9898"
)

// Config set the plugin config
var Config = sdk.Configuration{
	Name:        "xiaomi",
	Version:     "1.0.0",
	Author:      "ItsJimi",
	Description: "xiaomi",
	Main:        "xiaomi",
	FuncData:    "onData",
	Discover:    true,
	Triggers: []sdk.Trigger{
		sdk.Trigger{
			Name:    "switch",
			FieldID: "SID",
			Fields: []sdk.Field{
				sdk.Field{
					Name:          "Status",
					Direct:        true,
					Type:          "string",
					Possibilities: []string{"click", "double_click", "long_click_press", "long_click_release"},
				},
			},
		},
		sdk.Trigger{
			Name:    "weatherv1",
			FieldID: "SID",
			Fields: []sdk.Field{
				sdk.Field{
					Name:   "Temperature",
					Direct: false,
					Type:   "int",
				},
				sdk.Field{
					Name:   "Humidity",
					Direct: false,
					Type:   "int",
				},
				sdk.Field{
					Name:   "Pressure",
					Direct: false,
					Type:   "int",
				},
			},
		},
		sdk.Trigger{
			Name:    "motion",
			FieldID: "SID",
			Fields: []sdk.Field{
				sdk.Field{
					Name:          "Status",
					Direct:        true,
					Type:          "string",
					Possibilities: []string{"motion"},
				},
				sdk.Field{
					Name:   "NoMotion",
					Direct: false,
					Type:   "int",
				},
			},
		},
		sdk.Trigger{
			Name:    "sensormotionaq2",
			FieldID: "SID",
			Fields: []sdk.Field{
				sdk.Field{
					Name:          "Status",
					Direct:        true,
					Type:          "string",
					Possibilities: []string{"motion"},
				},
				sdk.Field{
					Name:   "Lux",
					Direct: false,
					Type:   "int",
				},
				sdk.Field{
					Name:   "NoMotion",
					Direct: false,
					Type:   "int",
				},
			},
		},
		sdk.Trigger{
			Name:    "sensorcubeaqgl01",
			FieldID: "SID",
			Fields: []sdk.Field{
				sdk.Field{
					Name:          "Status",
					Direct:        true,
					Type:          "string",
					Possibilities: []string{"move", "tap_twice", "shake_air", "alert", "flip90", "flip180", "free_fall"},
				},
				sdk.Field{
					Name:   "rotate",
					Direct: true,
					Type:   "string",
				},
			},
		},
	},
	Actions: []sdk.Action{},
}

var conn *net.UDPConn

// Init config and store it to Casa server
func Init() []byte {
	return []byte("")
}

// OnStart start UDP server to get Xiaomi data
func OnStart(config []byte) {
	addr, err := net.ResolveUDPAddr("udp", ip+":"+port)
	if err != nil {
		log.Panic(err)
	}

	conn, err = net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Listening gateway events\n")
}

// OnData get data from xiaomi gateway
func OnData() interface{} {
	if conn == nil {
		log.Panic("No connection")
	}

	buf := make([]byte, 1024)
	var res Event

	for res.SID == "" {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Panic("Can't read udp", err)
		}

		err = json.Unmarshal(buf[0:n], &res)
		if err != nil {
			log.Println(err)
		}
	}

	switch res.Model {
	case "switch":
		data := []byte(res.Data.(string))
		var button devices.Switch
		err := json.Unmarshal(data, &button)
		button.SID = res.SID
		if err != nil {
			log.Println(err)
		}
		return &button
	case "weather.v1":
		data := []byte(res.Data.(string))
		var weather devices.WeatherV1
		err := json.Unmarshal(data, &weather)
		if weather.Temperature != "" {
			nbr, _ := strconv.ParseFloat(weather.Temperature, 32)
			weather.Temperature = fmt.Sprintf("%f", nbr/100)
		}
		if weather.Humidity != "" {
			nbr, _ := strconv.ParseFloat(weather.Humidity, 32)
			weather.Humidity = fmt.Sprintf("%f", nbr/100)
		}
		weather.SID = res.SID
		if err != nil {
			log.Println(err)
		}
		return &weather
	case "motion":
		data := []byte(res.Data.(string))
		var motion devices.Motion
		err := json.Unmarshal(data, &motion)
		motion.SID = res.SID
		if err != nil {
			log.Println(err)
		}
		return &motion
	case "sensor_motion.aq2":
		data := []byte(res.Data.(string))
		var motion devices.SensorMotionAQ2
		err := json.Unmarshal(data, &motion)
		motion.SID = res.SID
		if err != nil {
			log.Println(err)
		}
		return &motion
	// case "gateway":
	// 	data := []byte(res.Data.(string))
	// 	var gateway devices.Gateway
	// 	err := json.Unmarshal(data, &gateway)
	// 	gateway.SID = res.SID
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	return &gateway
	case "sensor_cube.aqgl01":
		data := []byte(res.Data.(string))
		var sensor devices.SensorCubeAqgl01
		err := json.Unmarshal(data, &sensor)
		sensor.SID = res.SID
		if err != nil {
			log.Println(err)
		}
		return &sensor
	default:
		// fmt.Println(res)
		return nil
	}
}

// OnStop close connection
func OnStop() {
	conn.Close()
}
