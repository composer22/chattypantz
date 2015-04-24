// Package server implements a chat server for websocket access.
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	// Allow dynamic profiling.
	_ "net/http/pprof"

	"github.com/composer22/chattypantz/logger"
	"golang.org/x/net/netutil"
	"golang.org/x/net/websocket"
)

// Server is the main structure that represents a server instance.
type Server struct {
	info     *Info             // Basic server information used to run the server.
	opts     *Options          // Original options used to create the server.
	stats    *Stats            // Server statistics since it started.
	mu       sync.Mutex        // For locking access to server attributes.
	running  bool              // Is the server running?
	log      *ChatLogger       // Log instance for recording error and other messages.
	roomMngr *ChatRoomManager  // Manager of chat rooms.
	chatters map[*Chatter]bool // A list of chatters connected to the server.
	done     chan bool         // A channel to signal to web socked to close.
	srvr     *http.Server      // HTTP server.
	wg       sync.WaitGroup    // Synchronization of channel close.
}

// New is a factory function that returns a new server instance.
func New(ops *Options) *Server {
	s := &Server{
		info: InfoNew(func(i *Info) {
			i.Name = ops.Name
			i.Hostname = ops.Hostname
			i.Port = ops.Port
			i.ProfPort = ops.ProfPort
			i.MaxConns = ops.MaxConns
			i.MaxRooms = ops.MaxRooms
			i.MaxIdle = ops.MaxIdle
			i.Debug = ops.Debug
		}),
		opts:     ops,
		stats:    StatsNew(),
		log:      ChatLoggerNew(),
		chatters: map[*Chatter]bool{},
		running:  false,
	}

	if s.info.Debug {
		s.log.SetLogLevel(logger.Debug)
	}

	// Setup the routes.
	http.Handle(wsRouteV1Conn, websocket.Handler(s.chatHandler))
	http.HandleFunc(httpRouteV1Alive, s.aliveHandler)
	http.HandleFunc(httpRouteV1Stats, s.statsHandler)
	s.srvr = &http.Server{
		Addr: fmt.Sprintf("%s:%d", s.info.Hostname, s.info.Port),
	}

	s.roomMngr = ChatRoomManagerNew(s.info.MaxRooms, s.log) // Set the manager of the chat rooms
	s.handleSignals()                                       // Evoke trap signals handler
	return s
}

// PrintVersionAndExit prints the version of the server then exits.
func PrintVersionAndExit() {
	fmt.Printf("chattypantz version %s\n", version)
	os.Exit(0)
}

// Start spins up the server to accept incoming connections.
func (s *Server) Start() error {
	if s.isRunning() {
		return errors.New("Server already started.")
	}

	s.log.Infof("Starting chattypantz version %s\n", version)

	// Construct listener
	ln, err := net.Listen("tcp", s.srvr.Addr)
	if err != nil {
		s.log.Errorf("Cannot create net.listener: %s", err.Error())
		return err
	}
	// If we want to limit connections, created a special listener with a throttle.
	if s.info.MaxConns > 0 {
		ln = netutil.LimitListener(ln, s.info.MaxConns)
	}

	s.mu.Lock()

	// Pprof http endpoint for the profiler.
	if s.info.ProfPort > 0 {
		s.StartProfiler()
	}

	s.stats.Start = time.Now()
	s.done = make(chan bool)
	s.running = true
	s.mu.Unlock()
	err = s.srvr.Serve(ln)

	// Done.
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
	if err != nil {
		s.log.Emergencyf("Listen and Server Error: %s", err.Error())
	}
	return nil
}

// StartProfiler is called to enable dynamic profiling.
func (s *Server) StartProfiler() {
	s.log.Infof("Starting profiling on http port %d", s.opts.ProfPort)
	hp := fmt.Sprintf("%s:%d", s.info.Hostname, s.info.ProfPort)
	go func() {
		err := http.ListenAndServe(hp, nil)
		if err != nil {
			s.log.Emergencyf("Error starting profile monitoring service: %s", err)
		}
	}()
}

// Shutdown takes down the server gracefully back to an initialize state.
func (s *Server) Shutdown() {
	if !s.isRunning() {
		return
	}
	s.log.Infof("BEGIN server service stop.")

	s.mu.Lock()
	s.log.Infof("\tShutting down chatters...")
	close(s.done)
	s.running = false
	s.mu.Unlock()
	s.wg.Wait()
	s.log.Infof("\tShutting down rooms...")
	s.roomMngr.shutDownRooms()
	s.log.Infof("END server service stop.")
}

// handleSignals responds to operating system interrupts such as application kills.
func (s *Server) handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			s.log.Infof("Server received signal: %v\n", sig)
			s.Shutdown()
			s.log.Infof("Server exiting.")
			os.Exit(0)
		}
	}()
}

// chatHandler is the main entry point to handle chat connections to the client.
func (s *Server) chatHandler(ws *websocket.Conn) {
	s.log.LogConnect(ws.Request())
	s.incrementStats(ws.Request())
	c := ChatterNew(s, ws)
	s.mu.Lock()
	s.chatters[c] = true // register chatter
	s.mu.Unlock()
	c.Run()
	s.mu.Lock()
	delete(s.chatters, c) // unregister chatter
	s.mu.Unlock()
}

// aliveHandler handles a client http:// "is the server alive?" request.
func (s *Server) aliveHandler(w http.ResponseWriter, r *http.Request) {
	s.log.LogConnect(r)
	s.incrementStats(r)
	s.initResponseHeader(w)
}

// statsHandler handles a client request for server information and statistics.
func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	s.log.LogConnect(r)
	s.incrementStats(r)
	s.initResponseHeader(w)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.ChatterStats = []*ChatterStats{}
	for c := range s.chatters {
		s.stats.ChatterStats = append(s.stats.ChatterStats, c.ChatterStatsNew())
	}
	s.stats.RoomStats = s.roomMngr.getRoomStats()
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)
	b, _ := json.Marshal(
		&struct {
			Info    *Info             `json:"info"`
			Options *Options          `json:"options"`
			Stats   *Stats            `json:"stats"`
			Memory  *runtime.MemStats `json:"memStats"`
		}{
			Info:    s.info,
			Options: s.opts,
			Stats:   s.stats,
			Memory:  mStats,
		})
	w.Write(b)
}

// initResponseHeader sets up the common http response headers for the return of all json calls.
func (s *Server) initResponseHeader(w http.ResponseWriter) {
	h := w.Header()
	h.Add("Content-Type", "application/json;charset=utf-8")
	h.Add("Date", time.Now().UTC().Format(time.RFC1123Z))
	if s.info.Name != "" {
		h.Add("Server", s.info.Name)
	}
	h.Add("X-Request-ID", createV4UUID())
}

// incrementStats increments the statistics for the request being handled by the server.
func (s *Server) incrementStats(r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.IncrReqStats(r.ContentLength)
	s.stats.IncrRouteStats(r.URL.Path, r.ContentLength)
}

// isRunning returns a boolean representing whether the server is running or not.
func (s *Server) isRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}
