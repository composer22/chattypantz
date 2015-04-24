package server

import (
	"encoding/json"
	"time"
)

// Stats contains runtime statistics for the server.
type Stats struct {
	Start        time.Time                   `json:"startTime"`    // The start time of the server.
	ReqCount     int64                       `json:"reqCount"`     // How many requests came in to the server.
	ReqBytes     int64                       `json:"reqBytes"`     // Size of the requests in bytes.
	RouteStats   map[string]map[string]int64 `json:"routeStats"`   // How many requests/bytes came into each route.
	ChatterStats []*ChatterStats             `json:"chatterStats"` // Statistics about each logged in chatter.
	RoomStats    []*ChatRoomStats            `json:"roomStats"`    // How many requests etc came into each room.
}

// StatsNew is a factory function that returns a new instance of statistics.
// options is an optional list of functions that initialize the structure
func StatsNew(opts ...func(*Stats)) *Stats {
	s := &Stats{
		Start:        time.Now(),
		RouteStats:   make(map[string]map[string]int64),
		ChatterStats: []*ChatterStats{},
		RoomStats:    []*ChatRoomStats{},
	}
	for _, f := range opts {
		f(s)
	}
	return s
}

// IncrReqStats increments the stats totals for the server.
func (s *Stats) IncrReqStats(b int64) {
	s.ReqCount++
	if b > 0 {
		s.ReqBytes += b
	}
}

// IncrRouteStats increments the stats totals for the route.
func (s *Stats) IncrRouteStats(path string, rb int64) {
	if _, ok := s.RouteStats[path]; !ok {
		s.RouteStats[path] = make(map[string]int64)
	}

	s.RouteStats[path]["requestCount"]++
	if rb > 0 {
		s.RouteStats[path]["requestBytes"] += rb
	}
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (s *Stats) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}
