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
	fog.High()
	heater.High()
	fan.High()
	light.High()

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
	fan.Low()
	return true
}
func FanOff(interface{}) interface{} {
	fan.High()
	return false
}
func FogOn(interface{}) interface{} {
	fog.Low()
	return true
}
func FogOff(interface{}) interface{} {
	fog.High()
	return false
}
func LightOn(interface{}) interface{} {
	light.Low()
	return true
}
func LightOff(interface{}) interface{} {
	light.High()
	return false
}

func HeaterOn(interface{}) interface{} {
	heater.Low()
	return true
}
func HeaterOff(interface{}) interface{} {
	heater.High()
	return false
}
