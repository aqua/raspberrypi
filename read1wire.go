package main

// Dump state of a single DS18S20 1-wire device on an RPi

import (
	"flag"
	"fmt"
	"rpi/onewire"
	"time"
)

var id = flag.String("id", "", "Device to sample")
var delay = flag.Int64("delay", 3000, "Delay n msec between samples")

func main() {
	flag.Parse()

	device, err := onewire.NewDS18S20(*id)
	if err != nil {
		fmt.Printf("Error setting up 1wire slave %v: %v", *id, err)
		return
	}

	for {
		fmt.Printf("attempting read\n")
		value, err := device.Read()
		if err != nil {
			fmt.Printf("Error on read: %v", err)
		} else {
			fmt.Printf("Sample: %v°C\n", float64(value)/1000)
		}
		time.Sleep(time.Duration(*delay) * time.Millisecond)
	}
}

