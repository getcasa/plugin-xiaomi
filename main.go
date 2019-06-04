package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/ItsJimi/casa-xiaomi/devices"
)

func main() {}

const (
	ip   = "224.0.0.50"
	port = "9898"
)

var conn *net.UDPConn

// OnStart start UDP server to get Xiaomi data
func OnStart() {
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
		if err != nil {
			log.Println(err)
		}
		return &button
	case "weather.v1":
		data := []byte(res.Data.(string))
		var weather devices.WeatherV1
		err := json.Unmarshal(data, &weather)
		if err != nil {
			log.Println(err)
		}
		return &weather
	case "motion":
		data := []byte(res.Data.(string))
		var motion devices.Motion
		err := json.Unmarshal(data, &motion)
		if err != nil {
			log.Println(err)
		}
		return &motion
	case "sensor_motion.aq2":
		data := []byte(res.Data.(string))
		var motion devices.SensorMotionAQ2
		err := json.Unmarshal(data, &motion)
		if err != nil {
			log.Println(err)
		}
		return &motion
	default:
		return nil
	}
}

// OnStop close connection
func OnStop() {
	conn.Close()
}
