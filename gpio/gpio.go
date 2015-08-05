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

  if err := gpio.enable_export(); err != nil {
    return nil, err
  }
  err = gpio.SetDirection(direction)
  if err != nil {
    return nil, err
  }
  gpio.fd, err = os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", gpio.Number), os.O_RDWR|os.O_SYNC, 0666)
  if err != nil {
    return nil, err
  }
  return gpio, nil
}

func (gpio *GPIOLine) enable_export() error {
  _, err := os.Stat(fmt.Sprintf("/sys/class/gpio/gpio%d", gpio.Number))
  if err == nil {
    // already exported
    return nil
  } else if err != nil && !os.IsNotExist(err) {
    // some other error
    return err
  }
  fd, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY|os.O_SYNC, 0666)
  if err != nil {
    return err
  }
  _, err = fmt.Fprintf(fd, "%d\n", gpio.Number)
  return err
}

func (gpio *GPIOLine) SetDirection(direction int) error {
  df, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", gpio.Number),
    os.O_WRONLY|os.O_SYNC, 0666)
  if err != nil {
    return err
  }
  if direction > 0 {
    fmt.Fprintln(df, "out")
  } else {
    fmt.Fprintln(df, "in")
  }
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

func (gpio *GPIOLine) GetState() (bool) {
  data := make([]byte, 1)
  val := false
  _, err := gpio.fd.Read(data)
  err = nil
  if err != nil {
    fmt.Println("Error means nothing in this case, Needed for Read issue")
  }
  if data[0] != 49 {
    val = false
  } else {
    val = true
  }
  return val
}

func (gpio *GPIOLine) Close() {
  gpio.fd.Close()
}
