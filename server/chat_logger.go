package server

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/composer22/chattypantz/logger"
)

// ChatLogger is an enhancement over the base logger for application logging patterns.
type ChatLogger struct {
	*logger.Logger
}

// ChatLoggerNew is a factory function that returns a new ChatLogger instance.
func ChatLoggerNew() *ChatLogger {
	return &ChatLogger{
		logger.New(logger.UseDefault, false),
	}
}

// connectLogEntry is a datastructure for recording initial connection information.
type connectLogEntry struct {
	Method     string      `json:"method"`
	URL        *url.URL    `json:"url"`
	Proto      string      `json:"proto"`
	Header     http.Header `json:"header"`
	Host       string      `json:"host"`
	RemoteAddr string      `json:"remoteAddr"`
	RequestURI string      `json:"requestURI"`
}

// sessionLogEntry is a datastructure for recording general activity between client and server.
type sessionLogEntry struct {
	RemoteAddr string `json:"remoteAddr"`
	Message    string `json:"message"`
}

// LogConnect is used to log request information when the client first connects to the server.
func (l *ChatLogger) LogConnect(r *http.Request) {
	if l.GetLogLevel() >= logger.Info {
		b, _ := json.Marshal(&connectLogEntry{
			Method:     r.Method,
			URL:        r.URL,
			Proto:      r.Proto,
			Header:     r.Header,
			Host:       r.Host,
			RemoteAddr: r.RemoteAddr,
			RequestURI: r.RequestURI,
		})
		l.Output(3, logger.Labels[logger.Info], `{"connected":%s}`, string(b))
	}
}

// LogSession is used to record information received during the client's session.
func (l *ChatLogger) LogSession(tp string, addr string, msg string) {
	if l.GetLogLevel() >= logger.Info {
		b, _ := json.Marshal(&sessionLogEntry{
			RemoteAddr: addr,
			Message:    msg,
		})
		l.Output(3, logger.Labels[logger.Info], `{"%s":%s}`, tp, string(b))
	}
}

// LogError is used to record misc session error information between server and client.
func (l *ChatLogger) LogError(addr string, msg string) {
	if l.GetLogLevel() >= logger.Error {
		b, _ := json.Marshal(&sessionLogEntry{
			RemoteAddr: addr,
			Message:    msg,
		})
		l.Output(3, logger.Labels[logger.Error], `{"error":%s}`, string(b))
	}
}
