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
	LevelDebug   = 4
	LevelInfo    = 3
	LevelWarning = 2
	LevelError   = 1
	LevelNone    = 0
)

var logLevel = LevelDebug

func init() {
	logger = log.New(os.Stdout, "", log.LUTC)
}

// SetLogLevel set the logging verbosity. 0 is lowest (log nothing at all), 4 is
// highest (log all debug messages)
func SetLogLevel(level int) {
	if level < LevelNone || level > LevelDebug {
		Error("Invalid log level %v", level)
		return
	}
	logLevel = level
}

// Debug logs a debugging message
func Debug(msgFmt string, v ...interface{}) {
	if logLevel < LevelDebug {
		return
	}
	print("DBG", msgFmt, v...)
}

// Info logs an informative message
func Info(msgFmt string, v ...interface{}) {
	if logLevel < LevelInfo {
		return
	}
	print("INF", msgFmt, v...)
}

// Warn logs a warning message
func Warn(msgFmt string, v ...interface{}) {
	if logLevel < LevelWarning {
		return
	}
	print("WRN", msgFmt, v...)
}

// Error logs an error message, and prints an execution stack afterwards
func Error(msgFmt string, v ...interface{}) {
	if logLevel < LevelError {
		return
	}
	print("ERR", msgFmt, v...)
	debug.PrintStack()
}

func print(lvl string, msgFmt string, v ...interface{}) {
	_, fn, line, _ := runtime.Caller(2)
	var cutoff = 30
	if len(fn) < cutoff {
		cutoff = len(fn)
	}

	msg := fmt.Sprintf("[%s] %30s:%-3d %s", lvl, "…"+string(fn[len(fn)-cutoff:]), line, msgFmt)

	if len(v) == 0 {
		logger.Println(msg)
	} else {
		logger.Printf(msg, v...)
	}
}
