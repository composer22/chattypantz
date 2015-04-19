package server

import (
	"encoding/json"
	"time"
)

// Status contains runtime statistics for the server.
type Status struct {
	Start        time.Time                   `json:"startTime"`    // The start time of the server.
	ReqCount     int64                       `json:"reqCount"`     // How many requests came in to the server.
	ReqBytes     int64                       `json:"reqBytes"`     // Size of the requests in bytes.
	ConnNumAvail int                         `json:"connNumAvail"` // Number of live connections available.
	RoomStats    map[string]map[string]int64 `json:"roomStats"`    // How many requests/bytes came into each room.
}

// StatusNew is a factory function that returns a new instance of Status.
// options is an optional list of functions that initialize the structure
func StatusNew(ops ...func(*Status)) *Status {
	st := &Status{
		Start:        time.Now(),
		ConnNumAvail: -1, // defaults to infinite.
		RoomStats:    make(map[string]map[string]int64),
	}
	for _, f := range ops {
		f(st)
	}
	return st
}

// IncrReqStats increments the stats totals for the server.
func (s *Status) IncrReqStats(b int64) {
	s.ReqCount++
	if b > 0 {
		s.ReqBytes += b
	}
}

// IncrRoomStats increments the stats totals for the chat room.
func (s *Status) IncrRoomStats(rm string, b int64) {
	if _, ok := s.RoomStats[rm]; !ok {
		s.RoomStats[rm] = make(map[string]int64)
	}

	s.RoomStats[rm]["reqCount"]++
	if b > 0 {
		s.RoomStats[rm]["reqBytes"] += b
	}
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (s *Status) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}
