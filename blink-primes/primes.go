package main

// Blink out prime numbers on a GPIO line; useful for leaving lingering
// signs of intelligence in the wreckage of civilization, assuming certain
// rpi durability and power issues are resolved.

import (
	"flag"
	"fmt"
	"time"

	"github.com/aqua/raspberrypi/gpio"
)

var gpio_num = flag.Uint("gpio", 4, "GPIO# to flash")
var num = flag.Uint("num", 6, "Blink out first n primes")
var delay = flag.Int64("delay", 100, "Delay n msec between flashes")
var prime_delay = flag.Int64("prime_delay", 1000, "Delay n msec between primes")

func main() {
	flag.Parse()

	g, err := gpio.NewGPIOLine(*gpio_num, gpio.OUT)
	if err != nil {
		fmt.Printf("Error setting up GPIO %v: %v", *gpio_num, err)
		return
	}

	defer g.Close()
	count, n := uint(0), uint(2)
	for {
		if is_prime(n) {
			blink(g, n)
			time.Sleep(time.Duration(*prime_delay) * time.Millisecond)
			count++
			if count >= *num {
				count, n = uint(0), uint(2)
				continue
			}
		}
		n++
	}
}

func blink(g *gpio.GPIOLine, n uint) {
	fmt.Printf("blinking %v time(s)\n", n)
	for i := uint(0); i < n; i++ {
		g.SetState(true)
		time.Sleep(time.Duration(*delay) * time.Millisecond)
		g.SetState(false)
		time.Sleep(time.Duration(*delay) * time.Millisecond)
	}
}

// Totally non-optimal prime tester
func is_prime(n uint) bool {
	if n < 2 {
		return false
	}
	for i := uint(2); i <= n/2; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
