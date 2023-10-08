package onewire

import (
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestScan(t *testing.T) {
	os.Setenv("ONEWIRE_BUS_DEVICE_PATH", "testdata/devices/")
	got, err := Scan()
	if err != nil {
		t.Errorf("got error in Scan(): %v", err)
		return
	}
	want := []string{"28-000003d44061", "28-000003d44bad", "28-000003d45c8f"}
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Scan yielded %v, want %v", got, want)
		return
	}
}
