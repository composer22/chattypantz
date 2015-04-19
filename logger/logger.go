// Package logger provides a custom logging abstract over the standard out logging of golang.
// All logging should by go to stdout according to 12-factor principles.
// Logging levels are based on RFC 5424: http://www.rfc-base.org/rfc-5424.html#
package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
)

// Standard labels.
const (
	//  RFC 5424 log levels.
	Emergency = iota
	Alert
	Critical
	Error
	Warning
	Notice
	Info
	Debug

	UseDefault = -1 // Note: literal consts must follow any iota decls else unexpected results.
)

const (
	// ANSI 8 colours.
	foregroundBlack = iota + 30
	foregroundRed
	foregroundGreen
	foregroundYellow
	foregroundBlue
	foregroundMagenta
	foregroundCyan
	foregroundLightGrey
	_
	foregroundDefault

	colourFormat = "[\x1b[%dm%s\x1b[0m] "
)

var (
	// Log labels.
	labels = []string{"[EMERGENCY] ",
		"[ALERT] ",
		"[CRITICAL] ",
		"[ERROR] ",
		"[WARNING] ",
		"[NOTICE] ",
		"[INFO] ",
		"[DEBUG] ",
	}
)

// Wrap the os.Exit() function so we can mock/test or customize exit.
type exiter func(code int)

// Logger provides a datastructure for all logging state.
type Logger struct {
	logger *log.Logger
	level  int
	labels []string
	exit   exiter
}

// New is a factory method to return a new logger instance.
func New(lvl int, clrs bool) *Logger {
	flags := log.Lshortfile | log.Ldate | log.Lmicroseconds
	pre := fmt.Sprintf("[%d] ", os.Getpid())
	if lvl == UseDefault {
		lvl = Info
	}

	l := &Logger{
		logger: log.New(os.Stdout, pre, flags),
		level:  lvl,
		exit:   func(code int) { os.Exit(code) },
	}

	if clrs {
		l.SetColouredLabels()
	} else {
		l.SetPlainLabels()
	}
	return l
}

// SetLogLevel allows a user to set the log level of the logger.
func (l *Logger) SetLogLevel(lvl int) error {
	if lvl < UseDefault || lvl > Debug {
		return errors.New(fmt.Sprintf("%d log level arg is not in valid range.", lvl))
	}

	if lvl == UseDefault {
		lvl = Info
	}
	l.level = lvl
	return nil
}

// SetExitFunc allows a user to set the exit function of the logger.
func (l *Logger) SetExitFunc(e exiter) error {
	if e == nil {
		return errors.New("Exit function is manadatory.")
	}
	l.exit = e
	return nil
}

// GetLogLevel returns the current log level of the logger.
func (l *Logger) GetLogLevel() int {
	return l.level
}

// SetPlainLabels sets the message labels to simple text output.
func (l *Logger) SetPlainLabels() {
	copy(l.labels, labels)
}

// SetColouredLabels sets the message labels to colourized text output.
func (l *Logger) SetColouredLabels() {
	l.labels = make([]string, 0)
	for i, lbl := range labels {
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
		l.labels = append(l.labels, fmt.Sprintf(colourFormat, clr, lbl))
	}
}

// Emergencyf prints an emergency message to the system log,
// This is considered an unrecoverable error and the application also exits, unless dont exit = true.
func (l *Logger) Emergencyf(format string, v ...interface{}) {
	if l.level >= Emergency {
		l.Output(3, labels[Emergency], format, v...)
	}
	l.performExit(l.exit)
}

// Alertf prints an alert message to the system log.
func (l *Logger) Alertf(format string, v ...interface{}) {
	if l.level >= Alert {
		l.Output(3, labels[Alert], format, v...)
	}
}

// Criticalf prints a critical message to the system log.
func (l *Logger) Criticalf(format string, v ...interface{}) {
	if l.level >= Critical {
		l.Output(3, labels[Critical], format, v...)
	}
}

// Errorf prints an error message to the system log.
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level >= Error {
		l.Output(3, labels[Error], format, v...)
	}
}

// Warningf prints a warning message to the system log.
func (l *Logger) Warningf(format string, v ...interface{}) {
	if l.level >= Warning {
		l.Output(3, labels[Warning], format, v...)
	}
}

// Noticef prints a notice message to the system log.
func (l *Logger) Noticef(format string, v ...interface{}) {
	if l.level >= Notice {
		l.Output(3, labels[Notice], format, v...)
	}
}

// Infof prints an informational message to the system log.
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level >= Info {
		l.Output(3, labels[Info], format, v...)
	}
}

// Debugf prints a debug message to the system log.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level >= Debug {
		l.Output(3, labels[Debug], format, v...)
	}
}

// Output prints a message directly into the system log. Normally, you should use level message functions.
// so that level can trap the write.
func (l *Logger) Output(cd int, lbl string, format string, v ...interface{}) error {
	var d int = 2
	if cd > 0 {
		d = cd
	}
	return l.logger.Output(d, fmt.Sprintf(lbl+format, v...))
}

// performExit wraps the application exit point wih a custom closure/anonymous function.
func (l *Logger) performExit(xit exiter) {
	xit(1) // call the exiter function
}
