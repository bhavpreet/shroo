package main

import (
	"fmt"
	"os"

	rpio "github.com/stianeikeland/go-rpio"
)

var (
	fog, fan, light, heater rpio.Pin
)

func Init(periph interface{}) interface{} {

	p := periph.(map[string]interface{})
	err := rpio.Open()
	if err != nil {
		fmt.Println("Some error occured", err)
		os.Exit(-1)
	}

	// defer rpio.Close()

	heater = rpio.Pin(p["heater"].(int))
	heater.Output()
	fog = rpio.Pin(p["fog"].(int))
	fog.Output()
	fan = rpio.Pin(p["fan"].(int))
	fan.Output()
	light = rpio.Pin(p["light"].(int))
	light.Output()

	// All OFF
	fog.Low()
	heater.Low()
	fan.Low()
	light.Low()

	return nil
}

func InitDebug(periph interface{}) interface{} {

	p := periph.(map[string]interface{})

	heater = rpio.Pin(p["heater"].(int))
	fog = rpio.Pin(p["fog"].(int))
	fan = rpio.Pin(p["fan"].(int))
	light = rpio.Pin(p["light"].(int))

	return nil
}

func FanOn(interface{}) interface{} {
	fan.High()
	return true
}
func FanOff(interface{}) interface{} {
	fan.Low()
	return false
}
func FogOn(interface{}) interface{} {
	fog.High()
	return true
}
func FogOff(interface{}) interface{} {
	fog.Low()
	return false
}
func LightOn(interface{}) interface{} {
	light.High()
	return true
}
func LightOff(interface{}) interface{} {
	light.Low()
	return false
}

func HeaterOn(interface{}) interface{} {
	heater.High()
	return true
}
func HeaterOff(interface{}) interface{} {
	heater.Low()
	return false
}
