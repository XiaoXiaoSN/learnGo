package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/goburrow/modbus"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
)

var (
	b *backoff.Backoff

	host = "0.tcp.ngrok.io:17408"

	slave   byte          = 1
	timeout time.Duration = time.Second

	modbusClient modbus.Client
)

type Power struct {
	Index int  `yaml:"index"`
	Watt  bool `yaml:"watt"`
	Ac    bool `yaml:"ac"`
}

type Light struct {
	ID       int `yaml:"id"`
	Register struct {
		Switch  int   `yaml:"switch"`
		Dimming int   `yaml:"dimming"`
		Power   Power `yaml:"power"`
	} `yaml:"register"`
}

func NewModbusClient(host string, slave byte, timeout time.Duration) (modbus.Client, error) {
	handler := modbus.NewTCPClientHandler(host)
	handler.Timeout = timeout
	handler.SlaveId = slave

	if err := handler.Connect(); err != nil {
		log.Warn("Modbus Connection Refuse!")
		return nil, err
	}
	defer handler.Close()

	return modbus.NewClient(handler), nil
}

func getLight(modbusClient modbus.Client) {
	// GET raw data
	rawSwitch, err := StableModbus(modbusClient.ReadCoils, 0, 1)
	if err != nil {
		log.Warnf("get switch err: %+v\n", err)
		return
	}

	rawDimming, err := StableModbus(modbusClient.ReadHoldingRegisters, 0, 1)
	if err != nil {
		log.Warnf("get dimming err: %+v\n", err)
		return
	}

	// parsar data format
	realSwitch := true
	if int32(rawSwitch[0]) == 0 {
		realSwitch = false
	}
	realDimming := int32((binary.BigEndian.Uint16(rawDimming)) / 100)

	fmt.Printf("switch:  %+v\n", realSwitch)
	fmt.Printf("dimming: %+v\n", realDimming)
	fmt.Println()
}

func setLightSwitch(modbusClient modbus.Client, lightSwitch bool) {
	// set switch use writeCoil

	var value uint16
	if lightSwitch {
		value = 65280 // 0xFF00
	}

	rawSwitch, err := modbusClient.WriteSingleCoil(0, value)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}

	realSwitch := true
	if int32(rawSwitch[0]) == 0 {
		realSwitch = false
	}
	fmt.Printf("switch:  %+v\n", realSwitch)
}

func setLightDimming(modbusClient modbus.Client, dimming uint16) {
	// set dimming use writeRegisters

	value := dimming * 100

	rawDimming, err := modbusClient.WriteSingleRegister(0, value)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}

	realDimming := int32((binary.BigEndian.Uint16(rawDimming)) / 100)
	fmt.Printf("dimming: %+v\n", realDimming)
}

func StableModbus(fn func(uint16, uint16) ([]byte, error), address, quantity uint16) (results []byte, err error) {
	retry := 3
	for i := 0; i < retry; i++ {
		results, err = fn(address, quantity)
		if err == nil {
			b.Reset()
			return
		}

		// backoff
		d := b.Duration()
		log.Printf("%s, reconnecting in %s", err, d)
		time.Sleep(d)

		// retry to connect
		m, err := NewModbusClient(host, slave, timeout)
		if err == nil {
			log.Warn("connection reseted")
			modbusClient = m
		}
	}

	return
}

func main() {
	b = &backoff.Backoff{
		Min: 100 * time.Millisecond,
		Max: 20 * time.Second,
	}

	m, err := NewModbusClient(host, slave, timeout)
	if err != nil {
		log.Panic(err)
	}
	modbusClient = m

	// setLightSwitch(modbusClient, true)
	// setLightDimming(modbusClient, 60)

	for {
		fmt.Println(time.Now().Format(time.Stamp))
		getLight(modbusClient)
		time.Sleep(3 * time.Second)
	}
}
