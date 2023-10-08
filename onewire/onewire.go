package onewire

import (
	"os"
	"strings"
)

/* Find all 1-wire devices attached to bus masters */
func Scan() ([]string, error) {
	devicedir, err := os.Open("/sys/bus/w1/devices/")
	if err != nil {
		return nil, err
	}
	defer devicedir.Close()
	names, err := devicedir.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	r := make([]string, 0, len(names)-1)
	for i := range names {
		if !strings.Contains(names[i], "w1_bus_master") {
			r = append(r, names[i])
		}
	}
	return r, nil
}
