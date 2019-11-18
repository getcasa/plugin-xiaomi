package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
	"time"

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
	Description: "Control xiaomi ecosystem",
	Devices: []sdk.Device{
		sdk.Device{
			Name:           "switch",
			DefaultTrigger: "",
			DefaultAction:  "",
			Triggers: []sdk.Trigger{
				sdk.Trigger{
					Name:          "Status",
					Direct:        true,
					Type:          "string",
					Possibilities: []string{"click", "double_click", "long_click_press", "long_click_release"},
				},
			},
			Actions: []string{},
		},
		sdk.Device{
			Name:           "sensorht",
			DefaultTrigger: "",
			DefaultAction:  "",
			Triggers: []sdk.Trigger{
				sdk.Trigger{
					Name:   "Temperature",
					Direct: false,
					Type:   "int",
				},
				sdk.Trigger{
					Name:   "Humidity",
					Direct: false,
					Type:   "int",
				},
			},
			Actions: []string{},
		},
		sdk.Device{
			Name:           "weatherv1",
			DefaultTrigger: "",
			DefaultAction:  "",
			Triggers: []sdk.Trigger{
				sdk.Trigger{
					Name:   "Temperature",
					Direct: false,
					Type:   "int",
				},
				sdk.Trigger{
					Name:   "Humidity",
					Direct: false,
					Type:   "int",
				},
				sdk.Trigger{
					Name:   "Pressure",
					Direct: false,
					Type:   "int",
				},
			},
			Actions: []string{},
		},
		sdk.Device{
			Name:           "motion",
			DefaultTrigger: "",
			DefaultAction:  "",
			Triggers: []sdk.Trigger{
				sdk.Trigger{
					Name:          "Status",
					Direct:        true,
					Type:          "string",
					Possibilities: []string{"motion"},
				},
				sdk.Trigger{
					Name:   "NoMotion",
					Direct: false,
					Type:   "int",
				},
			},
		},
		sdk.Device{
			Name:           "sensormagnetaq2",
			DefaultTrigger: "",
			DefaultAction:  "",
			Triggers: []sdk.Trigger{
				sdk.Trigger{
					Name:          "Status",
					Direct:        true,
					Type:          "string",
					Possibilities: []string{"open", "close"},
				},
			},
		},
		sdk.Device{
			Name:           "sensormotionaq2",
			DefaultTrigger: "",
			DefaultAction:  "",
			Triggers: []sdk.Trigger{
				sdk.Trigger{
					Name:          "Status",
					Direct:        true,
					Type:          "string",
					Possibilities: []string{"motion"},
				},
				sdk.Trigger{
					Name:   "Lux",
					Direct: false,
					Type:   "int",
				},
				sdk.Trigger{
					Name:   "NoMotion",
					Direct: false,
					Type:   "int",
				},
			},
		},
		sdk.Device{
			Name:           "sensorcubeaqgl01",
			DefaultTrigger: "",
			DefaultAction:  "",
			Triggers: []sdk.Trigger{
				sdk.Trigger{
					Name:          "Status",
					Direct:        true,
					Type:          "string",
					Possibilities: []string{"move", "tap_twice", "shake_air", "alert", "flip90", "flip180", "free_fall"},
				},
				sdk.Trigger{
					Name:   "rotate",
					Direct: true,
					Type:   "string",
				},
			},
		},
		sdk.Device{
			Name:           "vibration",
			DefaultTrigger: "",
			DefaultAction:  "",
			Triggers:       []sdk.Trigger{},
		},
	},
	Actions: []sdk.Action{},
}

var conn *net.UDPConn
var gateways []devices.Gateway
var devs []sdk.DiscoveredDevice
var addr *net.UDPAddr

type xiaomi struct {
	SID            string   `json:"sid"`
	Status         string   `json:"status"`
	IP             string   `json:"ip"`
	Token          string   `json:"token"`
	Devices        []string `json:"data"`
	RGB            int      `json:"rgb"`
	Illumination   int      `json:"illumination"`
	Rotate         string   `json:"rotate"`
	NoMotion       string   `json:"no_motion"`
	Lux            string   `json:"lux"`
	Voltage        int      `json:"voltage"`
	BedActivity    string   `json:"bed_activity"`
	Coordination   string   `json:"coordination"`
	FinalTiltAngle string   `json:"final_tilt_angle"`
	Temperature    string   `json:"temperature"`
	Humidity       string   `json:"humidity"`
	Pressure       string   `json:"pressure"`
}

// OnStart start UDP server to get Xiaomi data
func OnStart(config []byte) {
	var err error
	addr, err = net.ResolveUDPAddr("udp", ip+":"+port)
	if err != nil {
		log.Panic(err)
	}

	conn, err = net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Listening gateway events\n")
}

// Discover return array of all found devices
func Discover() []sdk.DiscoveredDevice {
	return devs
}

// OnData get data from xiaomi gateway
func OnData() []sdk.Data {
	var datas []sdk.Data
	if conn == nil {
		log.Panic("No connection")
	}

	buf := make([]byte, 1024)
	var res Event
	var err error
	var n int

	for res.SID == "" {
		n, _, err = conn.ReadFromUDP(buf)
		if err != nil {
			log.Panic("Can't read udp", err)
		}

		err = json.Unmarshal(buf[0:n], &res)
		if err != nil {
			log.Println(err)
		}
	}

	var newData sdk.Data
	physicalName := strings.Replace(strings.Replace(strings.ToLower(res.Model), ".", "", -1), "_", "", -1)

	switch res.CMD {
	case "get_id_list_ack":
		var datas []string

		if findGatewayFromSID(res.SID) != nil {
			break
		}
		err := json.Unmarshal([]byte(res.Data.(string)), &datas)
		if err != nil {
			fmt.Println(err)
		}
		gateways = append(gateways, devices.Gateway{
			SID:     res.SID,
			Token:   res.Token,
			Devices: datas,
		})

		go func() {
			for _, data := range datas {
				for findDeviceFromSID(data) == nil {
					_, err = conn.WriteToUDP([]byte(`{"cmd": "read", "sid": "`+data+`"}`), addr)
					if err != nil {
						log.Println(err)
					}
					time.Sleep(500 * time.Millisecond)
				}
			}
		}()
		return nil
	case "read_ack":
		if res.Model == "" {
			break
		}
		devs = append(devs, sdk.DiscoveredDevice{
			Name:         "",
			PhysicalID:   res.SID,
			PhysicalName: physicalName,
			Plugin:       Config.Name,
		})
		return nil
	}

	if res.Model == "gateway" && findGatewayFromSID(res.SID) == nil {
		conn.WriteToUDP([]byte(`{"cmd": "get_id_list", "sid": "`+res.SID+`"}`), addr)
		if err != nil {
			log.Println(err)
		}
		return nil
	}

	if res.Model != "" && sdk.FindDevicesFromName(Config.Devices, physicalName).Name != "" {
		data := []byte(res.Data.(string))
		device := new(xiaomi)
		err := json.Unmarshal(data, &device)
		if err != nil {
			log.Println(err)
		}

		newData = sdk.Data{
			Plugin:       Config.Name,
			PhysicalName: physicalName,
			PhysicalID:   res.SID,
		}
		for _, field := range sdk.FindDevicesFromName(Config.Devices, physicalName).Triggers {
			newData.Values = append(newData.Values, sdk.Value{
				Name:  field.Name,
				Value: []byte(reflect.ValueOf(device).Elem().FieldByName(field.Name).String()),
				Type:  field.Type,
			})
		}
		datas = append(datas, newData)
		return datas
	}

	return nil
}

// OnStop close connection
func OnStop() {
	conn.Close()
}

func findGatewayFromSID(sid string) *devices.Gateway {
	if len(gateways) == 0 {
		return nil
	}
	for _, gateway := range gateways {
		if gateway.SID == sid {
			return &gateway
		}
	}
	return nil
}

func findDeviceFromSID(sid string) *sdk.DiscoveredDevice {
	if len(devs) == 0 {
		return nil
	}
	for _, dev := range devs {
		if dev.PhysicalID == sid {
			return &dev
		}
	}
	return nil
}
