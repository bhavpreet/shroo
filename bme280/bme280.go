package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"
	"periph.io/x/periph/host"
)

var (
	dev  *bmxx80.Dev
	defs map[string]interface{}
)

type envt struct {
	Temp float64
	RH   float64
	PS   float64
}

var rete envt

func Init(_defs interface{}) interface{} {
	var err error

	// save args
	defs = _defs.(map[string]interface{})

	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open a handle to the first available I²C bus:
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	// defer bus.Close()

	// Open a handle to a bme280/bmp280 connected on the I²C bus using default
	// settings:
	dev, err = bmxx80.NewI2C(bus, 0x77, &bmxx80.Opts{
		Temperature: bmxx80.O8x,
		Pressure:    bmxx80.O8x,
		Humidity:    bmxx80.O8x,
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil
	// defer dev.Halt()
}

func InitDebug(_defs interface{}) interface{} {
	// save args
	defs = _defs.(map[string]interface{})
	return nil
}

func Update(args interface{}) interface{} {

	var err error
	env := new(physic.Env)
	// Read temperature from the sensor:
	if err = dev.Sense(env); err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("Unable to read sensor bme280")
	}

	rete.RH, err = strconv.ParseFloat(strings.Split(env.Humidity.String(), "%rH")[0], 64)
	rete.Temp, err = strconv.ParseFloat(strings.Split(env.Temperature.String(), "°C")[0], 64)
	rete.Temp = rete.Temp + defs["tempCorrection"].(float64)
	rete.PS, err = strconv.ParseFloat(strings.Split(env.Pressure.String(), "kPa")[0], 64)

	// log.Printf("%8s %10s %9s\n", st.env.Temperature, st.env.Pressure, st.env.Humidity)
	return rete
}

func UpdateDebug(args interface{}) interface{} {
	rete.Temp = 10.0 + defs["tempCorrection"].(float64)
	rete.RH = 10.0
	rete.PS = 10.0
	// log.Printf("%8s %10s %9s\n", st.env.Temperature, st.env.Pressure, st.env.Humidity)
	return rete
	// return nil
}
