package server

import "encoding/json"

// Info provides basic config information to/about the running server.
type Info struct {
	Version  string `json:"version"`      // Version of the server.
	Name     string `json:"name"`         // The name of the server.
	Hostname string `json:"hostname"`     // The hostname of the server.
	UUID     string `json:"UUID"`         // Unique ID of the server.
	Port     int    `json:"port"`         // Port the server is listening on.
	ProfPort int    `json:"profPort"`     // Profiler port the server is listening on.
	MaxConns int    `json:"maxConns"`     // The maximum concurrent clients accepted.
	MaxRooms int    `json:"maxRooms"`     // The maximum number of chat rooms allowed.
	MaxIdle  int    `json:"maxIdle"`      // The maximum client idle time in seconds before disconnect.
	Debug    bool   `json:"debugEnabled"` // Is debugging enabled on the server.
}

// InfoNew is a factory function that returns a new instance of Info.
// opts is an optional list of functions that initialize the structure
func InfoNew(opts ...func(*Info)) *Info {
	inf := &Info{
		Version: version,
		UUID:    createV4UUID(),
	}
	for _, f := range opts {
		f(inf)
	}
	return inf
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (i *Info) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}
