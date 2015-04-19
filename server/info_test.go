package server

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	testInfoExpectedJSONResult = `{"version":"9.8.7","name":"Test Server","hostname":"0.0.0.0",` +
		`"UUID":"ABCDEFGHIJKLMNOPQRSTUVWXYZ","port":6661,"profPort":6061,"maxConns":999,` +
		`"maxRooms":888,"maxHistory":777,"maxIdle":666,"debugEnabled":true}`
)

func TestInfoNew(t *testing.T) {
	info := InfoNew(func(i *Info) {
		i.Version = "9.8.7"
		i.Name = "Test Server"
		i.Hostname = "0.0.0.0"
		i.UUID = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		i.Port = 6661
		i.ProfPort = 6061
		i.MaxConns = 999
		i.MaxRooms = 888
		i.MaxHistory = 777
		i.MaxIdle = 666
		i.Debug = true
	})
	tp := reflect.TypeOf(info)

	if tp.Kind() != reflect.Ptr {
		t.Fatalf("Info not created as a pointer.")
	}

	tp = tp.Elem()
	if tp.Kind() != reflect.Struct {
		t.Fatalf("Info not created as a struct.")
	}
	if tp.Name() != "Info" {
		t.Fatalf("Info struct is not named correctly.")
	}
	if !(tp.NumField() > 0) {
		t.Fatalf("Info struct is empty.")
	}
}

func TestInfoString(t *testing.T) {
	t.Parallel()
	info := InfoNew(func(i *Info) {
		i.Version = "9.8.7"
		i.Name = "Test Server"
		i.Hostname = "0.0.0.0"
		i.UUID = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		i.Port = 6661
		i.ProfPort = 6061
		i.MaxConns = 999
		i.MaxRooms = 888
		i.MaxHistory = 777
		i.MaxIdle = 666
		i.Debug = true
	})
	actual := fmt.Sprint(info)
	if actual != testInfoExpectedJSONResult {
		t.Errorf("Info not converted to json string.\n\nExpected: %s\n\nActual: %s\n",
			testInfoExpectedJSONResult, actual)
	}
}
