package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/composer22/chattypantz/logger"
)

const (
	testChatLogExpCnt = `{"connected":{"method":"GET","proto":"HTTP/1.1","host":"example.com",` +
		`"requestURI":"ws://ladeda.com/v1.0/chat","remoteAddr":"127.8.9.10","header":{}}}`
	testChatLogExpSess = `{"disconnected":{"remoteAddr":"127.8.9.10","message":"Client disconnected."}}`
	testChatLogExpErr  = `{"error":{"remoteAddr":"127.8.9.10","message":"Couldn't receive. Error: Tester"}}`
)

func TestLogConnect(t *testing.T) {
	t.Parallel()
	testLbl := logger.Labels[logger.Info]
	r := &http.Request{
		Method:     "GET",
		Proto:      "HTTP/1.1",
		Header:     make(map[string][]string),
		Host:       "example.com",
		RequestURI: "ws://ladeda.com/v1.0/chat",
		RemoteAddr: "127.8.9.10",
	}
	expectOutput(t, func() {
		l := ChatLoggerNew()
		l.LogConnect(r)
	}, fmt.Sprintf("%s%s\n", testLbl, testChatLogExpCnt))
}

func TestLogSession(t *testing.T) {
	t.Parallel()
	testLbl := logger.Labels[logger.Info]
	expectOutput(t, func() {
		l := ChatLoggerNew()
		l.LogSession("disconnected", "127.8.9.10", "Client disconnected.")
	}, fmt.Sprintf("%s%s\n", testLbl, testChatLogExpSess))
}

func TestLogError(t *testing.T) {
	t.Parallel()
	testLbl := logger.Labels[logger.Error]
	expectOutput(t, func() {
		l := ChatLoggerNew()
		l.LogError("127.8.9.10", "Couldn't receive. Error: Tester")
	}, fmt.Sprintf("%s%s\n", testLbl, testChatLogExpErr))
}

// expectOutput is a helper function that repipes or mocks out stdout and allows error messages to be tested
// against the pipe.
func expectOutput(t *testing.T, f func(), expected string) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	os.Stdout.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC
	if !strings.Contains(out, expected) {
		t.Errorf("Expected '%s', received '%s'.", expected, out)
	}
}
