// Package server implements a chat server for websocket access.
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	// Allow dynamic profiling.
	_ "net/http/pprof"

	"github.com/composer22/chattypantz/logger"
	"golang.org/x/net/netutil"
	"golang.org/x/net/websocket"
)

// connectLogEntry is a datastructure for recording initial connection information.
type connectLogEntry struct {
	Method        string      `json:"method"`
	URL           *url.URL    `json:"url"`
	Proto         string      `json:"proto"`
	Header        http.Header `json:"header"`
	Body          string      `json:"body"`
	ContentLength int64       `json:"contentLength"`
	Host          string      `json:"host"`
	RemoteAddr    string      `json:"remoteAddr"`
	RequestURI    string      `json:"requestURI"`
	Trailer       http.Header `json:"trailer"`
}

// sessionLogEntry is a datastructure for recording general activity between client and server.
type sessionLogEntry struct {
	RemoteAddr string `json:"remoteAddr"`
	Message    string `json:"message"`
}

// Server is the main structure that represents a server instance.
type Server struct {
	info    *Info          // Basic server information used to run the server.
	opts    *Options       // Original options used to create the server.
	stats   *Status        // Server statistics since it started.
	mu      sync.Mutex     // For locking access to server params.
	running bool           // Is the server running?
	log     *logger.Logger // Log instance for recording error and other messages.
	srvr    *http.Server   // HTTP server.
	wg      sync.WaitGroup // Synchronization of channel close.
}

// New is a factory function that returns a new server instance.
func New(ops *Options, addedOpts ...func(*Server)) *Server {
	s := &Server{
		info: InfoNew(func(i *Info) {
			i.Name = ops.Name
			i.Hostname = ops.Hostname
			i.Port = ops.Port
			i.ProfPort = ops.ProfPort
			i.MaxConns = ops.MaxConns
			i.MaxRooms = ops.MaxRooms
			i.MaxHistory = ops.MaxHistory
			i.MaxIdle = ops.MaxIdle
			i.Debug = ops.Debug
		}),
		opts:    ops,
		stats:   StatusNew(),
		log:     logger.New(logger.UseDefault, false),
		running: false,
	}

	if s.info.Debug {
		s.log.SetLogLevel(logger.Debug)
	}

	// Setup the mutext, routes, middleware, and server.
	http.Handle(wsRouteV1Conn, websocket.Handler(s.echoHandler))

	s.srvr = &http.Server{
		Addr: fmt.Sprintf("%s:%d", s.info.Hostname, s.info.Port),
	}

	s.handleSignals() // Evoke trap signals handler

	// Additional hook for specialized custom options.
	for _, f := range addedOpts {
		f(s)
	}
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
	s.running = true
	s.mu.Unlock()

	err = s.srvr.Serve(ln)
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
	s.running = false
	s.mu.Unlock()

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

func (s *Server) echoHandler(ws *websocket.Conn) {
	var reply string
	r := ws.Request()
	s.LogConnect(r)
	remoteAddr := fmt.Sprint(r.RemoteAddr)
	for {
		// Set optional idle timeout.
		if s.info.MaxIdle > 0 {
			ws.SetReadDeadline(time.Now().Add(time.Duration(s.info.MaxIdle) * time.Second))
		}

		if err := websocket.Message.Receive(ws, &reply); err != nil {
			e, ok := err.(net.Error)
			switch {
			case ok && e.Timeout():
				s.LogSession("disconnected", remoteAddr, "Client forced to disconnect due to inactivity.")
			case err.Error() == "EOF":
				s.LogSession("disconnected", remoteAddr, "Client disconnected.")
			default:
				s.LogError(remoteAddr, fmt.Sprintf("Couldn't receive. Error: %s", err.Error()))
			}
			return
		}
		s.LogSession("received", remoteAddr, reply)

		msg := fmt.Sprintf("Received: %s", reply)
		if err := websocket.Message.Send(ws, msg); err != nil {
			switch {
			case err.Error() == "EOF":
				s.LogSession("disconnected", remoteAddr, "Client disconnected.")
			default:
				s.LogError(remoteAddr, fmt.Sprintf("Couldn't receive. Error: %s", err.Error()))
			}
			return
		}
	}
}

// LogConnect logs information from the socket when the client first connects to the server.
func (s *Server) LogConnect(r *http.Request) {
	var cl int64
	if r.ContentLength > 0 {
		cl = r.ContentLength
	}

	b, _ := json.Marshal(&connectLogEntry{
		Method:        r.Method,
		URL:           r.URL,
		Proto:         r.Proto,
		Header:        r.Header,
		ContentLength: cl,
		Host:          r.Host,
		RemoteAddr:    r.RemoteAddr,
		RequestURI:    r.RequestURI,
		Trailer:       r.Trailer,
	})
	s.log.Infof(`{"connected":%s}`, string(b))
}

// LogSession is used to record information received in the client's session.
func (s *Server) LogSession(tp string, addr string, msg string) {
	b, _ := json.Marshal(&sessionLogEntry{
		RemoteAddr: addr,
		Message:    msg,
	})
	s.log.Infof(`{"%s":%s}`, tp, string(b))
}

// LogError is used to record information regarding misc session errors between server and client.
func (s *Server) LogError(addr string, msg string) {
	b, _ := json.Marshal(&sessionLogEntry{
		RemoteAddr: addr,
		Message:    msg,
	})
	s.log.Errorf(`{"error":%s}`, string(b))
}

// isRunning returns a boolean representing whether the server is running or not.
func (s *Server) isRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}
