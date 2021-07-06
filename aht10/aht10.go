package main

import (
	"log"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
	"time"
)

var (
	dev  *i2c.Dev
	defs map[string]interface{}
)

type envt struct {
	Temp float64
	RH   float64
	PS   float64
}

var (
	rete     envt
	isInited bool
)

func initDevice() bool {
	if isInited {
		return isInited
	}

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	b, err := i2creg.Open("")
	if err != nil {
		log.Printf("Unable to to i2creg.Open, err: ", err)
		return false
	}

	// Dev is a valid conn.Conn.
	dev = &i2c.Dev{Addr: Address, Bus: b}

	// Initialize
	_, err = dev.Write([]byte{CMD_INITIALIZE, 0x00, 0x00})
	if err != nil {
		log.Printf("Failed to d.Write, on CMD_INITIALIZE, err: %v", err)
		return false
	}
	time.Sleep(500 * time.Millisecond)

	_, err = dev.Write([]byte{CMD_INITIALIZE, 0x08, 0x00})
	if err != nil {
		log.Printf("Failed to d.Write, on CMD_INITIALIZE, err: %v", err)
		return false
	}
	isInited = true
	time.Sleep(450 * time.Millisecond)
	return isInited
}

func Init(_defs interface{}) interface{} {
	defs = _defs.(map[string]interface{})
	initDevice()
	return nil
}

func InitDebug(_defs interface{}) interface{} {
	// save args
	defs = _defs.(map[string]interface{})
	return nil
}

func Update(args interface{}) interface{} {

	// Init device if it was previously in bad state
	if !initDevice() {
		log.Print("Unable to initialize the device!!")
		return rete
	}

	dev.Tx([]byte{CMD_TRIGGER, 0x00, 0x00}, nil)
	data := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	for retry := 0; retry < 3; retry++ {
		time.Sleep(80 * time.Millisecond)
		err := dev.Tx(nil, data)
		if err != nil {
			log.Printf("Faild at dev.Tx, err: %v", err)
			isInited = false
			break
		}

		if data[0]&0x68 == 0x08 {
			rete.RH = float64((uint32(data[1])<<12 | uint32(data[2])<<4 | (uint32(data[3])&0xf0)>>4) * 100.0 / (1 << 20))
			rete.Temp = float64(((uint32(data[3])&0xf)<<16|uint32(data[4])<<8|uint32(data[5]))*200.0/(1<<20) - 50)
			// log.Printf("Humidity = %v, temp = %v", humidity, temp)
			break
		}
	}

	rete.Temp = rete.Temp + defs["tempCorrection"].(float64)

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
