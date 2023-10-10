package gpio

import (
	"os"
	"testing"
)

func TestGPIOOut(t *testing.T) {
	os.Setenv("RPI_GPIO_PATH_OVERRIDE", "./testdata/sys/class/gpio")
	l, err := NewGPIOLine(1, OUT)
	if err != nil {
		t.Errorf("Error creating GPIOLine: %v", err)
		return
	}
	defer l.Close()
	if err := l.SetState(false); err != nil {
		t.Errorf("Error setting GPIOLine: %v", err)
	}
	if err := l.SetState(true); err != nil {
		t.Errorf("Error setting GPIOLine: %v", err)
	}
}

func TestGPIOIn(t *testing.T) {
	os.Setenv("RPI_GPIO_PATH_OVERRIDE", "./testdata/sys/class/gpio")
	l, err := NewGPIOLine(1, IN)
	if err != nil {
		t.Errorf("Error creating GPIOLine: %v", err)
		return
	}
	defer l.Close()
}
