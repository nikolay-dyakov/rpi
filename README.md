# Examples
Examples for getting started and testing [ANAVI hardware](http://anavi.technology/) written on Golang. Here you can find C and Python samples:[anavi-examples](https://github.com/AnaviTechnology/anavi-examples)

## require
1. [pigpio](http://abyz.me.uk/rpi/pigpio/) 

## Build for Raspberry Zero W

1. set following env variables `GOOS=linux;GOARCH=arm;GOARM=6`
2. go build

## Build for Raspberry 3

1. set following env variables `GOOS=linux;GOARCH=arm;GOARM=7`
2. go build 

## How to?
1. Run pigpiod on startup `sudo systemctl enable pigpiod`