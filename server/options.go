package server

import "encoding/json"

// Options represents parameters that are passed to the application to be used in constructing
// the server.
type Options struct {
	Name     string `json:"name"`         // The name of the server.
	Hostname string `json:"hostname"`     // The hostname of the server.
	Port     int    `json:"port"`         // The default port of the server.
	ProfPort int    `json:"profPort"`     // The profiler port of the server.
	MaxConns int    `json:"maxConns"`     // The maximum concurrent clients accepted.
	MaxRooms int    `json:"maxRooms"`     // The maximum number of chat rooms allowed.
	MaxIdle  int    `json:"maxIdle"`      // The maximum client idle time in seconds before disconnect.
	MaxProcs int    `json:"maxProcs"`     // The maximum number of processor cores available.
	Debug    bool   `json:"debugEnabled"` // Is debugging enabled in the application or server.
}

// String is an implentation of the Stringer interface so the structure is returned as a string
// to fmt.Print() etc.
func (o *Options) String() string {
	b, _ := json.Marshal(o)
	return string(b)
}
