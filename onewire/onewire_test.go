package onewire

import (
	"reflect"
	"sort"
	"testing"
)

func TestScan(t *testing.T) {
	linuxW1DevicePath = "testdata/devices/"
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
