// This is a simple logging package which can be used anywhere without any
// configuration. This package only logs to stdout and is supposed to be used in
// conjunction with an external system logger, like systemd-journal.

package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
)

var logger *log.Logger

// Log level, higher number is more verbosity
const (
	// LevelTrace is used for printing verbose network communications
	LevelTrace = 5

	// LevelDebug is used for printing the results of actions which are
	// essential to the succeeding of the program, but not important enough to
	// show up in the logs of every server which runs this program
	LevelDebug = 4

	// LevelInfo is used for printing messages which are useful to for knowing
	// the current state of the application. These messages should not show up
	// too often, and they should be useful for system administrators to see
	LevelInfo = 3

	// LevelWarning is used for notifying the administrator that something is
	// wrong, but the application can still function properly
	LevelWarning = 2

	// LevelError is used to indicate that something important broke and the
	// application is now in an inconsistent state. The problem should be fixed
	// before restarting the application
	LevelError = 1

	// LevelNone makes the application log no messages at all
	LevelNone = 0
)

var logLevel = LevelDebug

// defaultLevel is the level used for messages sent to the Write function
var defaultLevel = LevelDebug

// Colours controls if the log package should print ANSI colour codes depending
// on the log level of the logged message. Defaults to false
var Colours = false

func init() {
	logger = log.New(os.Stdout, "", log.LUTC)
}

// SetLogLevel set the logging verbosity. 0 is lowest (log nothing at all), 5 is
// highest (log all debug and trace messages)
func SetLogLevel(level int) {
	if level < LevelNone || level > LevelTrace {
		Error("Invalid log level %v", level)
		return
	}
	logLevel = level
}

// SetDefaultLevel set the log level for the log messages written to the Write
// function. This is useful for integrating packages which use Go's native
// logger interface into this log package.
//
// An example of this is the very verbose http logger which tends to spam logs
// with messages which can otherwise not be silenced.
func SetDefaultLevel(level int) {
	if level < LevelNone || level > LevelDebug {
		Error("Invalid log level %v", level)
		return
	}
	defaultLevel = level
}

// Trace logs a tracing message
func Trace(msgFmt string, v ...interface{}) {
	if logLevel < LevelTrace {
		return
	}
	print("95", "TRC", msgFmt, v...)
}

// Debug logs a debugging message
func Debug(msgFmt string, v ...interface{}) {
	if logLevel < LevelDebug {
		return
	}
	print("96", "DBG", msgFmt, v...)
}

// Info logs an informative message
func Info(msgFmt string, v ...interface{}) {
	if logLevel < LevelInfo {
		return
	}
	print("92", "INF", msgFmt, v...)
}

// Warn logs a warning message
func Warn(msgFmt string, v ...interface{}) {
	if logLevel < LevelWarning {
		return
	}
	print("93", "WRN", msgFmt, v...)
}

// Error logs an error message, and prints an execution stack afterwards
func Error(msgFmt string, v ...interface{}) {
	if logLevel < LevelError {
		return
	}
	print("91", "ERR", msgFmt, v...)
	debug.PrintStack()
}

type writer int

// Write can be used as a logging destination in the log.Logger interface. It
// logs a message to the default logging level, this level can be set with the
// SetDefaultLevel function
func (writer) Write(p []byte) (n int, err error) {
	switch defaultLevel {
	case LevelDebug:
		Debug(string(p))
	case LevelInfo:
		Info(string(p))
	case LevelWarning:
		Warn(string(p))
	case LevelError:
		Error(string(p))
	}
	return len(p), nil
}

// Logger is an instance of the standard Go log.Logger which can be used by Go
// packages to log to the Default log level
var Logger = log.New(writer(0), "", 0)

func print(colour string, lvl string, msgFmt string, v ...interface{}) {
	// Get the file name and line number
	_, fn, line, _ := runtime.Caller(2)

	// Maximum length of the file path which is printed
	var cutoff = 30
	if len(fn) < cutoff {
		cutoff = len(fn)
	}

	// If colour codes are enabled we add some ANSI magic to the mix
	if Colours {
		lvl = "\x1b[1m\x1b[" + colour + "m" + lvl + "\x1b[0m"
	}

	// Format the message to print. First the log level, then the source file
	// name, line number and the message
	msg := fmt.Sprintf(
		"[%s] %30s:%-3d %s",
		lvl,
		"…"+string(fn[len(fn)-cutoff:]),
		line,
		msgFmt,
	)

	// If variadic arguments were passed we parse them with Printf, else we just
	// print the message normally
	if len(v) == 0 {
		logger.Println(msg)
	} else {
		logger.Printf(msg, v...)
	}
}
