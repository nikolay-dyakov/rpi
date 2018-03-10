package light

import (
	"fmt"
	"os/exec"
	log "github.com/sirupsen/logrus"
)

const (
	redPin   int8   = 9
	bluePin  int8   = 10
	greenPin int8   = 11
	pigs     string = "/usr/bin/pigs"
)

func SetRGB(r, g, b uint8) {
	ClearRGB()
	arg := fmt.Sprintf("p %d %d p %d %d p %d %d", redPin, r, greenPin, g, bluePin, b)
	cmd := exec.Command(pigs, arg)
	err := cmd.Run()
	if err != nil {
		log.Error(err)
		return
	}
}

func ClearRGB() {
	arg := fmt.Sprintf("p %d 0 p %d 0 p %d 0", redPin, greenPin, bluePin)
	cmd := exec.Command(pigs, arg)
	err := cmd.Run()
	if err != nil {
		log.Error(err)
		return
	}
}
