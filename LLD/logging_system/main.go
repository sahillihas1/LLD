package main

import (
	"fmt"
	"os"
	"time"
)

type Logger interface {
	Log(message string)
}

type ILogDecorator interface {
	DecorateLog(message string) string
}

type INFODecorator struct {
}

func (i *INFODecorator) DecorateLog(message string) string {
	return ""
}

type ERRORDecorator struct{}

func (e *ERRORDecorator) DecorateLog(message string) string {
	return ""
}

type WARNINGDecorator struct{}

func (w *WARNINGDecorator) DecorateLog(message string) string {
	return ""
}

func getDecoratorFactory(logType string) ILogDecorator {
	switch logType {
	case "INFO":
		return &INFODecorator{}
	case "ERROR":
		return &ERRORDecorator{}
	case "WARNING":
		return &WARNINGDecorator{}
	}
	return nil
}

type ConsoleLogger struct {
}

func (c ConsoleLogger) Log(message string) {
	message = getDecoratorFactory("INFO").DecorateLog(message)
	fmt.Println(message)
}

// Concrete FileLogger - Logs messages to a file
type FileLogger struct {
	file *os.File
}

func NewFileLogger(filename string) *FileLogger {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	return &FileLogger{file: file}
}

func (f FileLogger) Log(message string) {
	f.file.WriteString(message + "\n")
}

// Factory Pattern - Logger Factory
func LoggerFactory(loggerType, filename string) Logger {
	if loggerType == "console" {
		return ConsoleLogger{}
	} else if loggerType == "file" {
		return NewFileLogger(filename)
	}
	return nil
}

// Decorator Pattern - Adding Timestamp
type TimestampLogger struct {
	logger Logger
}

func (t TimestampLogger) Log(message string) {
	t.logger.Log(time.Now().Format("2006-01-02 15:04:05") + " - " + message)
}

// Decorator Pattern - JSON Format Logger
type JSONLogger struct {
	logger Logger
}

func (j JSONLogger) Log(message string) {
	j.logger.Log(fmt.Sprintf("{\"timestamp\":\"%s\", \"message\":\"%s\"}", time.Now().Format("2006-01-02 15:04:05"), message))
}

// Builder Pattern - Configuring a Logger with multiple decorators
type LoggerBuilder struct {
	logger Logger
}

func NewLoggerBuilder() *LoggerBuilder {
	return &LoggerBuilder{}
}

func (b *LoggerBuilder) SetLogger(logger Logger) *LoggerBuilder {
	b.logger = logger
	return b
}

func (b *LoggerBuilder) AddTimestamp() *LoggerBuilder {
	b.logger = TimestampLogger{b.logger}
	return b
}

func (b *LoggerBuilder) AddJSONFormat() *LoggerBuilder {
	b.logger = JSONLogger{b.logger}
	return b
}

func (b *LoggerBuilder) Build() Logger {
	return b.logger
}

// Main function
func main() {
	// Using Factory Pattern to create a logger
	consoleLogger := LoggerFactory("console", "")
	fileLogger := LoggerFactory("file", "log.txt")

	// Simple loggers
	consoleLogger.Log("This is a console log message")
	fileLogger.Log("This is a file log message")

	// Using Builder Pattern to configure a logger
	logger := NewLoggerBuilder().SetLogger(LoggerFactory("console", "")).AddTimestamp().AddJSONFormat().Build()
	logger.Log("This is a structured log message")
}
