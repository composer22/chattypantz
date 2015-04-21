package server

import (
	"runtime"
	"testing"
	"time"
)

const ()

var (
	testSrvr *Server
)

func TestServerStartup(t *testing.T) {
	opts := &Options{
		Name:       "Test Server",
		Hostname:   "localhost",
		Port:       6660,
		ProfPort:   6060,
		MaxConns:   1000,
		MaxRooms:   1000,
		MaxHistory: 15,
		MaxIdle:    0,
		MaxProcs:   1,
		Debug:      true,
	}
	runtime.GOMAXPROCS(1)
	testSrvr = New(opts)
	go func() { testSrvr.Start() }()
}

func TestServerPrintVersion(t *testing.T) {
	t.Parallel()
	t.Skip("Exit cannot be covered.")
}

func TestServerIsRunning(t *testing.T) {
	time.Sleep(2 * time.Second) // Make sure we are all ready.
	if !testSrvr.isRunning() {
		t.Errorf("Server should be runnning.")
	}
}

func TestServerTakeDown(t *testing.T) {
	testSrvr.Shutdown()
	if testSrvr.isRunning() {
		t.Errorf("Server should have shut down.")
	}
	testSrvr = nil
}
