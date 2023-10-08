package onewire

import (
	"os"
	"regexp"
	"strings"
)

var linuxW1DevicePath = "/sys/bus/w1/devices/"

var deviceIDRE = regexp.MustCompile(`^[0-9a-f]{2}-[0-9a-f]+$`)

/* Find all 1-wire devices attached to bus masters */
func Scan() ([]string, error) {
	found := map[string]bool{}
	devicedir, err := os.Open(linuxW1DevicePath)
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
			if !found[names[i]] {
				r = append(r, names[i])
				found[names[i]] = true
			}
		}
	}
	return r, nil
}
