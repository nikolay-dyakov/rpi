package bh1750

import (
	"os"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/nikolay-dyakov/rpi/i2c"
)

const (
	BH1750_ADDR uint8 = 0x23
	BH1750_LUX  uint8 = 0x10
)

var conn *i2c.I2C
var err error

func init() {
	conn, err = i2c.New(BH1750_ADDR, 1)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func GetLux() (int16, error) {
	_, err := conn.WriteByte(BH1750_LUX)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	time.Sleep(500 * time.Millisecond)
	result, err := conn.ReadRegS16BE(0x00)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return result, nil
}
