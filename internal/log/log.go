// Package log implements a custom logger for Zorra
package log

import "fmt"

// Verbose is a flag specifying whether messages will be logged at the debug level
var Verbose bool

// Info logs a message at the info level
func Info(msg string) {
	log(msg, green+"INFO")
}

// Infof logs a message with formatting at the info level
func Infof(format string, args ...any) {
	Info(fmt.Sprintf(format, args...))
}

// Debug logs a message at the debug level
func Debug(msg string) {
	if !Verbose {
		return
	}

	log(msg, blue+"DEBUG")
}

// Debugf logs a message with formatting at the debug level
func Debugf(format string, args ...any) {
	Debug(fmt.Sprintf(format, args...))
}

// Warn logs a message at the warning level
func Warn(msg string) {
	log(msg, yellow+"WARN")
}

// Warnf logs a message with formatting at the warning level
func Warnf(format string, args ...any) {
	Warn(fmt.Sprintf(format, args...))
}

// Error logs a message at the error level
func Error(msg string) {
	log(msg, red+"ERROR")
}

// Errorf logs a message with formatting at the error level
func Errorf(format string, args ...any) {
	Error(fmt.Sprintf(format, args...))
}

// log is a helper function to log the message with a given prefix
func log(msg, prefix string) {
	fmt.Println("[" + prefix + reset + "] " + msg)
}
