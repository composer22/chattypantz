// chattypantz is a simple chat server demonstrating some of golangs socket functions and features.
package main

import (
	"flag"
	"runtime"
	"strings"

	"github.com/composer22/chattypantz/server"
)

var (
	log *server.ChatLogger = server.ChatLoggerNew()
)

// main is the main entry point for the application or server launch.
func main() {
	opts := server.Options{}
	var showVersion bool

	flag.StringVar(&opts.Name, "N", "", "Name of the server.")
	flag.StringVar(&opts.Name, "--name", "", "Name of the server.")
	flag.StringVar(&opts.Hostname, "H", server.DefaultHostname, "Hostname of the server")
	flag.StringVar(&opts.Hostname, "--hostname", server.DefaultHostname, "Name of the server")
	flag.IntVar(&opts.Port, "p", server.DefaultPort, "Port to listen on.")
	flag.IntVar(&opts.Port, "--port", server.DefaultPort, "Port to listen on.")
	flag.IntVar(&opts.ProfPort, "L", server.DefaultProfPort, "Profiler port to listen on.")
	flag.IntVar(&opts.ProfPort, "--profiler_port", server.DefaultProfPort, "Profiler port to listen on.")
	flag.IntVar(&opts.MaxConns, "n", server.DefaultMaxConns, "Maximum client connections allowed.")
	flag.IntVar(&opts.MaxConns, "--connections", server.DefaultMaxConns, "Maximum client connections allowed.")
	flag.IntVar(&opts.MaxRooms, "r", server.DefaultMaxRooms, "Maximum chat rooms allowed.")
	flag.IntVar(&opts.MaxRooms, "--rooms", server.DefaultMaxRooms, "Maximum chat rooms allowed.")
	flag.IntVar(&opts.MaxHistory, "y", server.DefaultMaxHistory, "Maximum chat room history allowed.")
	flag.IntVar(&opts.MaxHistory, "--history", server.DefaultMaxHistory, "Maximum chat room history allowed.")
	flag.IntVar(&opts.MaxIdle, "i", server.DefaultMaxIdle, "Maximum client idle allowed.")
	flag.IntVar(&opts.MaxIdle, "--idle", server.DefaultMaxIdle, "Maximum client idle allowed.")
	flag.IntVar(&opts.MaxProcs, "X", server.DefaultMaxProcs, "Maximum processor cores to use.")
	flag.IntVar(&opts.MaxProcs, "--procs", server.DefaultMaxProcs, "Maximum processor cores to use.")
	flag.BoolVar(&opts.Debug, "d", false, "Enable debugging output.")
	flag.BoolVar(&opts.Debug, "--debug", false, "Enable debugging output.")
	flag.BoolVar(&showVersion, "V", false, "Show version.")
	flag.BoolVar(&showVersion, "--version", false, "Show version.")
	flag.Usage = server.PrintUsageAndExit
	flag.Parse()

	// Version flag request?
	if showVersion {
		server.PrintVersionAndExit()
	}

	// Check additional params beyond the flags.
	for _, arg := range flag.Args() {
		switch strings.ToLower(arg) {
		case "version":
			server.PrintVersionAndExit()
		case "help":
			server.PrintUsageAndExit()
		}
	}

	// Set thread and proc usage.
	if opts.MaxProcs > 0 {
		runtime.GOMAXPROCS(opts.MaxProcs)
	}
	log.Infof("NumCPU %d GOMAXPROCS: %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(-1))

	s := server.New(&opts)
	s.Start()
}
