# Enable i2c to your raspberry
1. `sudo raspi-config`
2. go to `5 Interfacing Options`
3. go to `P5 I2C`
4. enable i2c interface
5. `sudo apt-get install -y python-smbus i2c-tools`
6. check the device address `i2cdetect -y 1`