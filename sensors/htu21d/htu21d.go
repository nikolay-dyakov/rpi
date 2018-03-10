package htu21d

import (
	"os"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/nikolay-dyakov/rpi/i2c"
)

const (
	HTU21D_I2C_ADDR uint8 = 0x40
	HTU21D_TEMP     uint8 = 0xF3
	HTU21D_HUMID    uint8 = 0xF5
)

var conn *i2c.I2C
var err error

func init() {
	conn, err = i2c.New(HTU21D_I2C_ADDR, 1)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func GetHumidity() (float64, error) {
	buf := make([]byte, 2)
	_, err := conn.WriteByte(HTU21D_HUMID)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	time.Sleep(100 * time.Millisecond)
	_, err = conn.Read(buf)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	result := float64((uint16(buf[0])<<8|uint16(buf[1]))&0xFFFC) / 65536.0
	return -6.0 + (125.0 * result), nil
}

func GetTemperature() (float64, error) {
	buf := make([]byte, 2)
	_, err := conn.WriteByte(HTU21D_TEMP)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	time.Sleep(100 * time.Millisecond)
	_, err = conn.Read(buf)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	result := float64((uint16(buf[0])<<8|uint16(buf[1]))&0xFFFC) / 65536.0
	return -46.85 + (175.72 * result), nil
}
