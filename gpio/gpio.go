package gpio

import (
	"bufio"
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

func linuxGPIOPath() string {
	if p := os.Getenv("RPI_GPIO_PATH_OVERRIDE"); p != "" {
		return p
	}
	return "/sys/class/gpio"
}

func NewGPIOLine(number uint, direction int) (gpio *GPIOLine, err error) {
	gpio = new(GPIOLine)
	gpio.Number = number

	if err := gpio.enable_export(); err != nil {
		return nil, err
	}
	err = gpio.SetDirection(direction)
	if err != nil {
		return nil, err
	}
	gpio.fd, err = os.OpenFile(fmt.Sprintf("%s/gpio%d/value", linuxGPIOPath(), gpio.Number), os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		return nil, err
	}
	return gpio, nil
}

func (gpio *GPIOLine) enable_export() error {
	_, err := os.Stat(fmt.Sprintf("%s/gpio%d", linuxGPIOPath(), gpio.Number))
	if err == nil {
		// already exported
		return nil
	} else if err != nil && !os.IsNotExist(err) {
		// some other error
		return err
	}
	fd, err := os.OpenFile(fmt.Sprintf("%s/export", linuxGPIOPath()), os.O_WRONLY|os.O_SYNC, 0666)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(fd, "%d\n", gpio.Number)
	return err
}

func (gpio *GPIOLine) SetDirection(direction int) error {
	df, err := os.OpenFile(fmt.Sprintf("%s/gpio%d/direction", linuxGPIOPath(), gpio.Number),
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

func (gpio *GPIOLine) GetState() (bool, error) {
	if _, err := gpio.fd.Seek(0, os.SEEK_SET); err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(gpio.fd)
	t := ""
	for scanner.Scan() {
		t = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return t != "0", nil
}

func (gpio *GPIOLine) Close() {
	gpio.fd.Close()
}
