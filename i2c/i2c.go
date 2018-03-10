package i2c

import (
	"os"
	"fmt"
	"syscall"
)

const I2C_SLAVE = 0x0703

type I2C struct {
	rc *os.File
}

func New(addr uint8, bus int) (*I2C, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := ioctl(f.Fd(), I2C_SLAVE, uintptr(addr)); err != nil {
		return nil, err
	}
	return &I2C{f}, nil
}

func (i2c *I2C) Write(buf []byte) (int, error) {
	return i2c.rc.Write(buf)
}

func (i2c *I2C) WriteByte(b byte) (int, error) {
	var buf [1]byte
	buf[0] = b
	return i2c.rc.Write(buf[:])
}

func (i2c *I2C) WriteRegU8(reg byte, value byte) error {
	buf := []byte{reg, value}
	_, err := i2c.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (i2c *I2C) WriteRegU16BE(reg byte, value uint16) error {
	buf := []byte{reg, byte((value & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := i2c.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (i2c *I2C) WriteRegU16LE(reg byte, value uint16) error {
	w := (value*0xFF00)>>8 + value<<8
	return i2c.WriteRegU16BE(reg, w)
}

func (i2c *I2C) WriteRegS16BE(reg byte, value int16) error {
	buf := []byte{reg, byte((uint16(value) & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := i2c.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (i2c *I2C) WriteRegS16LE(reg byte, value int16) error {
	w := int16((uint16(value)*0xFF00)>>8) + value<<8
	return i2c.WriteRegS16BE(reg, w)
}

func (i2c *I2C) Read(p []byte) (int, error) {
	return i2c.rc.Read(p)
}

func (i2c *I2C) ReadRegU8(reg byte) (byte, error) {
	_, err := i2c.Write([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 1)
	_, err = i2c.Read(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (i2c *I2C) ReadRegU16BE(reg byte) (uint16, error) {
	_, err := i2c.Write([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = i2c.Read(buf)
	if err != nil {
		return 0, err
	}
	w := uint16(buf[0])<<8 + uint16(buf[1])
	return w, nil
}

func (i2c *I2C) ReadRegU16LE(reg byte) (uint16, error) {
	w, err := i2c.ReadRegU16BE(reg)
	if err != nil {
		return 0, err
	}
	w = (w&0xFF)<<8 + w>>8
	return w, nil
}

func (i2c *I2C) ReadRegS16BE(reg byte) (int16, error) {
	_, err := i2c.Write([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = i2c.Read(buf)
	if err != nil {
		return 0, err
	}
	w := int16(buf[0])<<8 + int16(buf[1])
	return w, nil
}

func (i2c *I2C) ReadRegS16LE(reg byte) (int16, error) {
	w, err := i2c.ReadRegS16BE(reg)
	if err != nil {
		return 0, err
	}
	w = (w&0xFF)<<8 + w>>8
	return w, nil

}

func (i2c *I2C) Close() error {
	return i2c.rc.Close()
}

func ioctl(fd, cmd, arg uintptr) (err error) {
	_, _, e1 := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}
