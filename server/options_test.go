package server

import (
	"fmt"
	"testing"
)

const (
	testOptionsExpectedJSONResult = `{"name":"Test Options","hostname":"0.0.0.0","port":6661,` +
		`"profPort":6061,"maxConns":1001,"maxRooms":999,"maxHistory":888,"maxIdle":777,` +
		`"maxProcs":666,"debugEnabled":true}`
)

func TestOptionsString(t *testing.T) {
	t.Parallel()
	opts := &Options{
		Name:       "Test Options",
		Hostname:   "0.0.0.0",
		Port:       6661,
		ProfPort:   6061,
		MaxConns:   1001,
		MaxRooms:   999,
		MaxHistory: 888,
		MaxIdle:    777,
		MaxProcs:   666,
		Debug:      true,
	}
	actual := fmt.Sprint(opts)
	if actual != testOptionsExpectedJSONResult {
		t.Errorf("Options not converted to json string.\n\nExpected: %s\n\nActual: %s\n",
			testOptionsExpectedJSONResult, actual)
	}
}
