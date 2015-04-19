package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestSetAndGetLogLevel(t *testing.T) {
	l := New(Debug, false)
	if l.GetLogLevel() != Debug {
		t.Errorf("Invalid Log Level Set.")
	}
	err := l.SetLogLevel(Info)
	if err != nil {
		t.Errorf("Set log level func should have been called correctly for value Info.")
	}
	if l.level != Info {
		t.Errorf("Set log level func should have set new value correctly.")
	}

	err = l.SetLogLevel(UseDefault)
	if err != nil {
		t.Errorf("Set log level func should have been called correctly for value UseDefault.")
	}

	if l.level != Info {
		t.Errorf("Set default log level should have set new value correctly.")
	}

	err = l.SetLogLevel(UseDefault - 1)
	if err == nil {
		t.Errorf("Low param value was not tested properly.")
	}
	err = l.SetLogLevel(Debug + 1)
	if err == nil {
		t.Errorf("High param value was not tested properly.")
	}
}

func TestDefaultSetLogLevel(t *testing.T) {
	l := New(UseDefault, false)
	if l.GetLogLevel() != Info {
		t.Errorf("Invalid default Log Level Set.")
	}
}

func TestSetErrorFunc(t *testing.T) {
	l := New(Debug, false)
	if err := l.SetExitFunc(nil); err == nil {
		t.Errorf("Invalid set exit function with nil.")
	}

	if err := l.SetExitFunc(func(code int) {}); err != nil {
		t.Errorf("Invalid set exit function with vald value.")
	}
}

func TestSetColourLabels(t *testing.T) {
	l := New(UseDefault, true)
	for i, actual := range l.labels {
		var clr int
		switch i {
		case Emergency, Alert, Critical, Error:
			clr = foregroundRed
		case Warning:
			clr = foregroundYellow
		case Notice:
			clr = foregroundGreen
		case Debug:
			clr = foregroundBlue
		default:
			clr = foregroundDefault
		}
		expected := fmt.Sprintf(colourFormat, clr, labels[i])
		if expected != actual {
			t.Errorf("Invalid colour label\nExpected:%s\nActual:%s", expected, actual)
		}
	}
}

func TestEmergencyf(t *testing.T) {
	t.Parallel()
	testMsg := "Emergencyf"
	expectOutput(t, func() {
		l := New(Debug, false) // Mock the exit so coverage can complete.
		l.exit = func(code int) {}
		l.Emergencyf(testMsg)
	}, fmt.Sprintf("%s%s\n", labels[Emergency], testMsg))
}

func TestAlertf(t *testing.T) {
	t.Parallel()
	testMsg := "Alertf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Alertf(testMsg)
	}, fmt.Sprintf("%s%s\n", labels[Alert], testMsg))
}

func TestCriticalf(t *testing.T) {
	t.Parallel()
	testMsg := "Criticalf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Criticalf(testMsg)
	}, fmt.Sprintf("%s%s\n", labels[Critical], testMsg))
}

func TestErrorf(t *testing.T) {
	t.Parallel()
	testMsg := "Errorf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Errorf(testMsg)
	}, fmt.Sprintf("%s%s\n", labels[Error], testMsg))
}

func TestWarningf(t *testing.T) {
	t.Parallel()
	testMsg := "Warningf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Warningf(testMsg)
	}, fmt.Sprintf("%s%s\n", labels[Warning], testMsg))
}

func TestNoticef(t *testing.T) {
	t.Parallel()
	testMsg := "Noticef"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Noticef(testMsg)
	}, fmt.Sprintf("%s%s\n", labels[Notice], testMsg))
}

func TestInfof(t *testing.T) {
	t.Parallel()
	testMsg := "Infof"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Infof(testMsg)
	}, fmt.Sprintf("%s%s\n", labels[Info], testMsg))
}

func TestDebugf(t *testing.T) {
	t.Parallel()
	testMsg := "Debugf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Debugf(testMsg)
	}, fmt.Sprintf("%s%s\n", labels[Debug], testMsg))
}

func TestOutputf(t *testing.T) {
	t.Parallel()
	testLbl := "[OUTPUT] "
	testMsg := "Output"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Output(-1, testLbl, testMsg)
	}, fmt.Sprintf("%s%s\n", testLbl, testMsg))
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
