package server

import (
	"fmt"
	"os"
)

const usageText = `
Description: chattypantz is a chat server allowing clients to text each other within rooms.

Usage: chattypantz [options...]

Server options:
    -N, --name NAME                  NAME of the server (default: empty field).
    -H, --hostname HOSTNAME          HOSTNAME of the server (default: localhost).
    -p, --port PORT                  PORT to listen on (default: 6660).
	-L, --profiler_port PORT         *PORT the profiler is listening on (default: off).
    -n, --connections MAX            *MAX client connections allowed (default: unlimited).
    -r, --rooms MAX                  *MAX chatrooms allowed (default: unlimited).
    -y, --history MAX                *MAX num of history records per room (default: 15).
    -i, --idle MAX                   *MAX idle time in seconds allowed (default: unlimited).
    -X, --procs MAX                  *MAX processor cores to use from the machine.

    -d, --debug                      Enable debugging output (default: false)

     *  Anything <= 0 is no change to the environment (default: 0).

Common options:
    -h, --help                       Show this message
    -V, --version                    Show version

Examples:

    # Server mode activated as "San Francisco" on host 0.0.0.0 port 6661;
	# 10 clients; 50 rooms; one hour idle allowed; 2 processors
    chattypantz -N "San Francisco" -H 0.0.0.0 -p 6661 -n 10 -r 50 -i 3600 -X 2
`

// end help text

// PrintUsageAndExit is used to print out command line options.
func PrintUsageAndExit() {
	fmt.Printf("%s\n", usageText)
	os.Exit(0)
}
