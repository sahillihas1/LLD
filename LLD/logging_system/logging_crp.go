package main

import "fmt"

// LogLevel represents different levels of logging
type LogLevel int

const (
	INFO LogLevel = iota
	DEBUG
	ERROR
)

// Logger interface
type Logger interface {
	SetNext(Logger)
	LogMessage(level LogLevel, message string)
}

// BaseLogger struct (implements Logger partially)
type BaseLogger struct {
	next Logger
}

func (b *BaseLogger) SetNext(next Logger) {
	b.next = next
}

// InfoLogger handles INFO level logs
type InfoLogger struct {
	BaseLogger
}

func (l *InfoLogger) LogMessage(level LogLevel, message string) {
	if level == INFO {
		fmt.Println("[INFO]:", message)
	} else if l.next != nil {
		l.next.LogMessage(level, message)
	}
}

// DebugLogger handles DEBUG level logs
type DebugLogger struct {
	BaseLogger
}

func (l *DebugLogger) LogMessage(level LogLevel, message string) {
	if level == DEBUG {
		fmt.Println("[DEBUG]:", message)
	} else if l.next != nil {
		l.next.LogMessage(level, message)
	}
}

// ErrorLogger handles ERROR level logs
type ErrorLogger struct {
	BaseLogger
}

func (l *ErrorLogger) LogMessage(level LogLevel, message string) {
	if level == ERROR {
		fmt.Println("[ERROR]:", message)
	} else if l.next != nil {
		l.next.LogMessage(level, message)
	}
}

func main() {
	// Create loggers
	infoLogger := &InfoLogger{}
	debugLogger := &DebugLogger{}
	errorLogger := &ErrorLogger{}

	// Set up the chain: INFO → DEBUG → ERROR
	infoLogger.SetNext(debugLogger)
	debugLogger.SetNext(errorLogger)

	// Test logging at different levels
	infoLogger.LogMessage(INFO, "This is an info message.")
	infoLogger.LogMessage(DEBUG, "This is a debug message.")
	infoLogger.LogMessage(ERROR, "This is an error message.")
}
