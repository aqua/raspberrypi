package gpio

import (
	"os"
	"testing"
)

func setupState(t *testing.T, v bool) error {
	w := "0\n"
	if v {
		w = "1\n"
	}
	if err := os.WriteFile("./testdata/sys/class/gpio/gpio1/value", []byte(w), 0666); err != nil {
		t.Fatalf("Error setting up fake GPIO file: %v", err)
	}
	return nil
}

func TestGPIOOut(t *testing.T) {
	setupState(t, false)
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
	setupState(t, false)
	os.Setenv("RPI_GPIO_PATH_OVERRIDE", "./testdata/sys/class/gpio")
	l, err := NewGPIOLine(1, OUT)
	if err != nil {
		t.Errorf("Error creating GPIOLine: %v", err)
		return
	}
	defer l.Close()
	if err := l.SetState(true); err != nil {
		t.Errorf("Error setting up true for GetState: %v", err)
	}
	if v, err := l.GetState(); err != nil {
		t.Errorf("Error in GetState: %v", err)
	} else if v != true {
		t.Errorf("want true from GetState got %v", v)
	}

	if err := l.SetState(false); err != nil {
		t.Errorf("Error setting up false for GetState: %v", err)
	}
	if v, err := l.GetState(); err != nil {
		t.Errorf("Error in GetState: %v", err)
	} else if v != false {
		t.Errorf("want false from GetState got %v", v)
	}
}
