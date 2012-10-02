package onewire

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type DS18S20 struct {
	Id string
	fd *os.File
	reader *bufio.Reader
}

func NewDS18S20(id string) (*DS18S20, error) {
	device := new(DS18S20)
	device.Id = id

	var err error
	device.fd, err = os.OpenFile(fmt.Sprintf("/sys/bus/w1/devices/%v/w1_slave", device.Id), os.O_RDONLY|os.O_SYNC, 0666)
	if err != nil {
		return nil, err
	}
	device.reader = bufio.NewReader(device.fd)
	return device, nil
}

var __CRC_CHECK_REGEX *regexp.Regexp = regexp.MustCompile(`crc=\w+\s(YES|NO)`)
var __TEMP_SAMPLE_REGEX *regexp.Regexp = regexp.MustCompile(`.*\st=(\d+)`)

func (device *DS18S20) Read() (millidegrees int64, err error) {
	device.fd.Seek(0, 0)
	for {
		line, err := device.reader.ReadString('\n')
		if err != nil {
			return 0., err
		} else if len(line) == 0 {
			return 0., fmt.Errorf("EOF without data from w1")
		} else {
			// fmt.Printf("read from %v: %v", device.Id, line)
			matches := __CRC_CHECK_REGEX.FindStringSubmatch(line)
			if len(matches) > 0 && matches[1] != "YES" {
				return 0., fmt.Errorf("CRC mismatch on read")
			}

			matches = __TEMP_SAMPLE_REGEX.FindStringSubmatch(line)
			if len(matches) > 0 {
				// fmt.Printf("matched for temp: %v", matches)
				v, err := strconv.ParseInt(matches[1], 10, 64)
				if err != nil {
					return 0., err
				}
				return v, nil
			}
		}
	}
	return 0., fmt.Errorf("Malformed/empty output from w1 slave")
}
