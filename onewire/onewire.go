package onewire

import (
	"os"
	"strings"
)

/* Find all attached 1-wire slave devices */
func ScanSlaves() ([]string, error) {
	devicedir, err := os.Open("/sys/bus/w1/devices/")
	if err != nil {
		return nil, err
	}
	names, err := devicedir.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	devicedir.Close()
	r := make([]string, 0, len(names)-1)
	for i := range names {
		if ! strings.Contains(names[i], "w1_bus_master") {
			r = append(r, names[i])
		}
	}
	return r, nil
}
