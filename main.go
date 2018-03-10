package main

import (
	"os"
	"time"
	"syscall"
	"os/signal"
	"math/rand"
	log "github.com/sirupsen/logrus"
	_ "github.com/nikolay-dyakov/rpi/light"
	"github.com/nikolay-dyakov/rpi/sensors/htu21d"
	"github.com/nikolay-dyakov/rpi/sensors/bh1750"
	"github.com/nikolay-dyakov/rpi/sensors/bmp180"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		//light.ClearRGB()
		log.Info("\nTurn off")
		os.Exit(0)
	}()
	for {

		//min1 := uint8(rand.Uint32() % 10)
		//min2 := uint8(rand.Uint32() % 20)
		//max := uint8(rand.Uint32() % 255)
		//switch rand.Uint32() % 2 {
		//case 0:
		//	light.SetRGB(max, min1, min2)
		//case 1:
		//	light.SetRGB(min1, max, min2)
		//default:
		//	light.SetRGB(min1, min2, max)
		//}
		humidity, _ := htu21d.GetHumidity()
		temperature, _ := htu21d.GetTemperature()
		lux, _ := bh1750.GetLux()
		log.Infof("%5.2f%%rh, %5.2fC, %d Lux", humidity, temperature, lux)

		pressure := bmp180.GetPressure()
		altitude := bmp180.GetAltitude()
		temp := bmp180.GetTemperature()
		log.Infof("%5.2fhPa, %5.2f, %5.2fC", pressure, altitude, temp)
		time.Sleep(time.Second)
	}
	os.Exit(0)
}
