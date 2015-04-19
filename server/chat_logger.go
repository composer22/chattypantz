package server

import (
	"encoding/json"
	"net/http"

	"github.com/composer22/chattypantz/logger"
)

// ChatLogger is an enhancement over the base logger for application logging patterns.
type ChatLogger struct {
	*logger.Logger
}

// ChatLoggerNew is a factory function that returns a new ChatLogger instance.
func ChatLoggerNew() *ChatLogger {
	return &ChatLogger{
		logger.Logger: logger.New(logger.UseDefault, false),
	}
}

// connectLogEntry is a datastructure for recording initial connection information.
type connectLogEntry struct {
	Method     string      `json:"method"`
	Proto      string      `json:"proto"`
	Host       string      `json:"host"`
	RequestURI string      `json:"requestURI"`
	RemoteAddr string      `json:"remoteAddr"`
	Header     http.Header `json:"header"`
}

// sessionLogEntry is a datastructure for recording general activity between client and server.
type sessionLogEntry struct {
	RemoteAddr string `json:"remoteAddr"`
	Message    string `json:"message"`
}

// LogConnect is used to log request information when the client first connects to the server.
func (l *ChatLogger) LogConnect(r *http.Request) {
	b, _ := json.Marshal(&connectLogEntry{
		Method:     r.Method,
		Proto:      r.Proto,
		Host:       r.Host,
		RequestURI: r.RequestURI,
		RemoteAddr: r.RemoteAddr,
		Header:     r.Header,
	})
	l.Infof(`{"connected":%s}`, string(b))
}

// LogSession is used to record information received during the client's session.
func (l *ChatLogger) LogSession(tp string, addr string, msg string) {
	b, _ := json.Marshal(&sessionLogEntry{
		RemoteAddr: addr,
		Message:    msg,
	})
	l.Infof(`{"%s":%s}`, tp, string(b))
}

// LogError is used to record misc session error information between server and client.
func (l *ChatLogger) LogError(addr string, msg string) {
	b, _ := json.Marshal(&sessionLogEntry{
		RemoteAddr: addr,
		Message:    msg,
	})
	l.Errorf(`{"error":%s}`, string(b))
}
