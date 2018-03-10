package bmp180

import (
	"os"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/nikolay-dyakov/rpi/i2c"
	"math"
)

var conn *i2c.I2C
var err error

const (
	BMP180_I2CADDR         = 0x77
	BMP180_ULTRALOWPOWER   = 0
	BMP180_STANDARD        = 1
	BMP180_HIGHRES         = 2
	BMP180_ULTRAHIGHRES    = 3
	BMP180_CAL_AC1         = 0xAA // R   Calibration data (16 bits)
	BMP180_CAL_AC2         = 0xAC // R   Calibration data (16 bits)
	BMP180_CAL_AC3         = 0xAE // R   Calibration data (16 bits)
	BMP180_CAL_AC4         = 0xB0 // R   Calibration data (16 bits)
	BMP180_CAL_AC5         = 0xB2 // R   Calibration data (16 bits)
	BMP180_CAL_AC6         = 0xB4 // R   Calibration data (16 bits)
	BMP180_CAL_B1          = 0xB6 // R   Calibration data (16 bits)
	BMP180_CAL_B2          = 0xB8 // R   Calibration data (16 bits)
	BMP180_CAL_MB          = 0xBA // R   Calibration data (16 bits)
	BMP180_CAL_MC          = 0xBC // R   Calibration data (16 bits)
	BMP180_CAL_MD          = 0xBE // R   Calibration data (16 bits)
	BMP180_CONTROL         = 0xF4
	BMP180_TEMPDATA        = 0xF6
	BMP180_PRESSUREDATA    = 0xF6
	BMP180_READTEMPCMD     = 0x2E
	BMP180_READPRESSURECMD = 0x34
	seaLevelPressure       = 101325
)

var ac1, ac2, ac3, b1, b2, mb, mc, md int16
var ac4, ac5, ac6 uint16
var oversampling uint

func init() {
	conn, err = i2c.New(BMP180_I2CADDR, 1)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if 0 > begin() {
		log.Error("BMP180 sensor module not found")
		os.Exit(1)
	}
}

func begin() int {
	oversampling = BMP180_ULTRAHIGHRES
	tmp, _ := conn.ReadRegU8(0xD0)
	if 0x55 != tmp {
		return -1
	}
	// Read calibration data
	ac1, _ = conn.ReadRegS16BE(BMP180_CAL_AC1)
	ac2, _ = conn.ReadRegS16BE(BMP180_CAL_AC2)
	ac3, _ = conn.ReadRegS16BE(BMP180_CAL_AC3)
	ac4, _ = conn.ReadRegU16BE(BMP180_CAL_AC4)
	ac5, _ = conn.ReadRegU16BE(BMP180_CAL_AC5)
	ac6, _ = conn.ReadRegU16BE(BMP180_CAL_AC6)

	b1, _ = conn.ReadRegS16BE(BMP180_CAL_B1)
	b2, _ = conn.ReadRegS16BE(BMP180_CAL_B2)

	mb, _ = conn.ReadRegS16BE(BMP180_CAL_MB)
	mc, _ = conn.ReadRegS16BE(BMP180_CAL_MC)
	md, _ = conn.ReadRegS16BE(BMP180_CAL_MD)

	return 0
}

func computeB5(ut uint) int32 {
	x1 := ((int32(ut) - int32(ac6)) * int32(ac5)) >> 15
	x2 := (int32(mc) << 11) / (int32(x1) + int32(md))
	return x1 + x2
}

func readRawTemperature() uint {
	conn.WriteRegU8(BMP180_CONTROL, BMP180_READTEMPCMD)
	time.Sleep(5 * time.Millisecond)
	rawTemp, _ := conn.ReadRegU16BE(BMP180_TEMPDATA)
	return uint(rawTemp)
}

func readRawPressure() int {
	conn.WriteRegU8(BMP180_CONTROL, byte(BMP180_READPRESSURECMD+(oversampling<<6)))
	switch oversampling {
	case BMP180_ULTRALOWPOWER:
		time.Sleep(5 * time.Millisecond)
	case BMP180_STANDARD:
		time.Sleep(8 * time.Millisecond)
	case BMP180_HIGHRES:
		time.Sleep(14 * time.Millisecond)
	default:
		time.Sleep(26 * time.Millisecond)
	}

	firstMeasure, _ := conn.ReadRegU16BE(BMP180_PRESSUREDATA)
	secondMeasure, _ := conn.ReadRegU8(byte(BMP180_PRESSUREDATA + 2))
	pressure := int(firstMeasure)
	pressure <<= 8
	pressure |= int(secondMeasure)
	pressure >>= (8 - oversampling)
	return pressure
}

func readPressure() int32 {
	ut := readRawTemperature()
	up := readRawPressure()
	b5 := computeB5(ut)
	b6 := b5 - 4000
	x1 := (int32(b2) * int32(b6*b6) >> 12) >> 11
	x2 := (int32(ac2) * int32(b6)) >> 11
	x3 := x1 + x2
	b3 := (((int32(ac1)*4 + x3) << oversampling) + 2) >> 2
	x1 = (int32(ac3) * b6) >> 13
	x2 = (int32(b1) * ((b6 * b6) >> 12)) >> 16
	x3 = ((x1 + x2) + 2) >> 2
	b4 := (uint32(ac4) * uint32(x3+32768)) >> 15
	var p int32
	b7 := ((uint32(up) - uint32(b3)) * (50000 >> oversampling))
	if b7 < 0x80000000 {
		p = int32((b7 << 1) / b4)
	} else {
		p = int32((b7 / b4) << 1)
	}
	x1 = (p >> 8) * (p >> 8)
	x1 = (x1 * 3038) >> 16
	x2 = (-7357 * p) >> 16
	p += (x1 + x2 + 3791) >> 4
	return p
}

func GetPressure() float64 {
	return float64(readPressure()) / 100.0
}

func GetAltitude() float64 {
	pressure := readPressure()
	return 44330 * (1 - math.Pow(float64(pressure)/seaLevelPressure, 0.190295))
}

func GetTemperature() float64 {
	ut := readRawTemperature()
	compensate := computeB5(ut)
	rawTemperature := ((compensate + 8) >> 4)
	return float64(rawTemperature) / 10
}
