package onewire

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"encoding/binary"
	"strconv"
)

type DS18S20 struct {
	Id     uint64		// Actual 64-bit ID burned into chip
	Name   string		// Name assigned by linux w1 drivers
	fd     *os.File
	reader *bufio.Reader
}

func NewDS18S20(name string) (*DS18S20, error) {
	device := new(DS18S20)
	device.Name = name

	var err error
	if device.Id, err = read_device_id(device); err != nil {
		return nil, err
	}
	device.fd, err = os.OpenFile(fmt.Sprintf("/sys/bus/w1/devices/%v/w1_slave", device.Name), os.O_RDONLY|os.O_SYNC, 0666)
	if err != nil {
		return nil, err
	}
	device.reader = bufio.NewReader(device.fd)
	return device, nil
}

func read_device_id(device *DS18S20) (uint64, error) {
	fn := fmt.Sprintf("/sys/bus/w1/devices/%v/id", device.Name)
	id_fd, err := os.OpenFile(fn, os.O_RDONLY, 0666)
	if err != nil {
		return 0, err
	}
	defer id_fd.Close()
	var ret uint64
	err = binary.Read(id_fd, binary.LittleEndian, &ret)
	if err != nil {
		return 0, fmt.Errorf("Error decoding %v device id: %v", fn, err)
	}
	return ret, nil
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
