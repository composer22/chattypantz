package server

import "time"

const (
	version           = "0.1.0"     // Application and server version.
	DefaultHostname   = "localhost" // The hostname of the server.
	DefaultPort       = 6660        // Port to receive requests: see IANA Port Numbers.
	DefaultProfPort   = 0           // Profiler port to receive requests. *
	DefaultMaxConns   = 0           // Maximum number of connections allowed. *
	DefaultMaxRooms   = 0           // Maximum number of chat rooms allowed. *
	DefaultMaxHistory = 15          // Maximum number of chat history records per room.
	DefaultMaxIdle    = 0           // Maximum idle seconds per user connection. *
	DefaultMaxProcs   = 0           // Maximum number of computer processors to utilize. *

	// * zeros = no change or no limitations or not enabled.

	// Listener and connections.
	TCPKeepAliveTimeout = 3 * time.Minute  // deprecated
	TCPReadTimeout      = 10 * time.Second // deprecated
	TCPWriteTimeout     = 30 * time.Second // deprecated

	// http: routes.
	wsRouteV1Conn = "/v1.0/chat"
)
