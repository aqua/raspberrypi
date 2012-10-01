package gpio

import (
	"fmt"
	"os"
)

type GPIOLine struct {
	Number uint
	fd     *os.File
}

const (
	IN = iota
	OUT
)

func NewGPIOLine(number uint, direction int) (gpio *GPIOLine, err error) {
	gpio = new(GPIOLine)
	gpio.Number = number

	if err := gpio.SetDirection(direction); err != nil {
		return nil, err
	}
	gpio.fd, err = os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", gpio.Number), os.O_WRONLY|os.O_SYNC, 0666)
	if err != nil {
		return nil, err
	}
	return gpio, nil
}

func (gpio *GPIOLine) SetDirection(direction int) error {
	df, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", gpio.Number),
		os.O_WRONLY|os.O_SYNC, 0666)
	if err != nil {
		return err
	}
	fmt.Fprintln(df, "out")
	df.Close()
	return nil
}

func (gpio *GPIOLine) SetState(state bool) error {
	v := "0"
	if state {
		v = "1"
	}
	_, err := fmt.Fprintln(gpio.fd, v)
	return err
}

func (gpio *GPIOLine) Close() {
	gpio.fd.Close()
}
