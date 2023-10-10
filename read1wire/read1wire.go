package main

// Dump state of a single DS18S20 1-wire device on an RPi

import (
	"flag"
	"log"
	"time"

	"github.com/aqua/raspberrypi/onewire"
)

var id = flag.String("id", "", "Device to sample")
var delay = flag.Int64("delay", 5000, "Delay n msec between samples")

func main() {
	flag.Parse()

	names, err := onewire.Scan()
	if err != nil {
		log.Printf("Error scanning 1wire names: %v", err)
		return
	}
	log.Printf("scan = %v\n", names)

	devices := make([]*onewire.DS18S20, len(names))
	for i := range names {
		devices[i], err = onewire.NewDS18S20(names[i])
		if err != nil {
			log.Printf("Error opening device %v: %v", devices[i], err)
			return
		}
	}

	if err != nil {
		log.Printf("Error setting up 1wire slave %v: %v", *id, err)
		return
	}

	for {
		for i := range devices {
			log.Printf("attempting read on device %x\n", devices[i].Id)
			value, err := devices[i].Read()
			if err != nil {
				log.Printf("Error on read: %v", err)
			} else {
				log.Printf("Sample: %v°C\n", float64(value)/1000)
			}
			time.Sleep(time.Duration(*delay) * time.Millisecond)
		}
	}
}
