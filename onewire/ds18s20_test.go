package onewire

import (
	"os"
	"testing"
)

var testDevices = []struct {
	name         string
	model        uint8
	humanModel   string
	humanID      string
	failed       bool
	millidegrees int64
}{
	{"28-000003d44061", MODEL_DS18B20, "ds18b20", "ds18b20-000003d44061", false, 24812},
	{"28-000003d44bad", MODEL_DS18B20, "ds18b20", "ds18b20-000003d44bad", true, 0},
	{"28-000003d45c8f", MODEL_DS18B20, "ds18b20", "ds18b20-000003d45c8f", false, 25500},
}

func TestSetup(t *testing.T) {
	os.Setenv("ONEWIRE_BUS_DEVICE_PATH", "testdata/devices/")

	for _, c := range testDevices {
		d, err := NewDS18S20(c.name)
		if err != nil {
			t.Errorf("NewDS18S20 on %s returned err %v", c.name, err)
			return
		}
		defer d.Close()
		if got := d.HumanId(); got != c.humanID {
			t.Errorf("%s: want human ID %q, got %q", c.name, c.name, got)
		}
		if d.FamilyCode != c.model {
			t.Errorf("%s: want family code %d, got %d", c.name, c.model, d.FamilyCode)
		}
		if got := d.Model(); got != c.humanModel {
			t.Errorf("%s: want model %q, got %q", c.name, c.humanModel, got)
		}
	}
}

func TestRead(t *testing.T) {
	os.Setenv("ONEWIRE_BUS_DEVICE_PATH", "testdata/devices/")

	for _, c := range testDevices {
		d, err := NewDS18S20(c.name)
		if err != nil {
			t.Errorf("NewDS18S20 on %s returned err %v", c.name, err)
			return
		}
		defer d.Close()
		got, err := d.Read()
		if err != nil && !c.failed {
			t.Errorf("%s: error reading: %v", c.name, err)
			return
		}
		if c.failed {
			if err == nil {
				t.Errorf("%s: expected CRC error but got success", c.name)
			}
			continue
		}
		if got != c.millidegrees {
			t.Errorf("%s: want temp read of %d millidegrees, got %d", c.name, c.millidegrees, got)
		}
	}
}
