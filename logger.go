package bunyan

import (
	"fmt"
	"log"
)

// Logger holds the log.Logger, its name, and any streams.
type Logger struct {
	*log.Logger
	name    string
	streams []StreamInterface
}

// NewLogger creates a new logger, given one or more streams
func NewLogger(name string, streams []StreamInterface) *Logger {
	return &Logger{
		name:    name,
		streams: streams,
	}
}

// AddStream adds a stream to the log.
func (l *Logger) AddStream(s StreamInterface) {
	l.streams = append(l.streams, s)
}

// Log logs the logentry.
func (l *Logger) Log(e *LogEntry) {
	e.setLogger(l)

	for _, stream := range l.streams {
		stream.Publish(e)
	}
}

// Logln logs the message at the specified log level.
func (l *Logger) Logln(level LogLevel, message string) *LogEntry {
	e := NewLogEntry(level, message)

	l.Log(e)

	return e
}

// LogF logs the formatted string at the specified log level.
func (l *Logger) LogF(level LogLevel, format string, values ...interface{}) *LogEntry {
	e := NewLogEntry(level, fmt.Sprintf(format, values...))

	l.Log(e)

	return e
}

// Trace logs the message at the Trace level.
func (l *Logger) Trace(message string) *LogEntry {
	e := NewLogEntry(Trace, message)

	l.Log(e)

	return e
}

// TraceF logs the formatted string at the Trace level.
func (l *Logger) TraceF(format string, values ...interface{}) *LogEntry {
	e := NewLogEntry(Trace, fmt.Sprintf(format, values...))

	l.Log(e)

	return e
}

// Debug logs the message at the Debug level.
func (l *Logger) Debug(message string) *LogEntry {
	e := NewLogEntry(Debug, message)

	l.Log(e)

	return e
}

// DebugF logs the formatted string at the Debug level.
func (l *Logger) DebugF(format string, values ...interface{}) *LogEntry {
	e := NewLogEntry(Debug, fmt.Sprintf(format, values...))

	l.Log(e)

	return e
}

// Info logs the message at the Info level.
func (l *Logger) Info(message string) *LogEntry {
	e := NewLogEntry(Info, message)

	l.Log(e)

	return e
}

// InfoF logs the formatted string at the Info level.
func (l *Logger) InfoF(format string, values ...interface{}) *LogEntry {
	e := NewLogEntry(Info, fmt.Sprintf(format, values...))

	l.Log(e)

	return e
}

// Warn logs the message at the Warn level.
func (l *Logger) Warn(message string) *LogEntry {
	e := NewLogEntry(Warn, message)

	l.Log(e)

	return e
}

// WarnF logs the formatted string at the Warn level.
func (l *Logger) WarnF(format string, values ...interface{}) *LogEntry {
	e := NewLogEntry(Warn, fmt.Sprintf(format, values...))

	l.Log(e)

	return e
}

// Error logs the message at the Error level.
func (l *Logger) Error(message string) *LogEntry {
	e := NewLogEntry(Error, message)

	l.Log(e)

	return e
}

// ErrorF logs the formatted string at the Error level.
func (l *Logger) ErrorF(format string, values ...interface{}) *LogEntry {
	e := NewLogEntry(Error, fmt.Sprintf(format, values...))

	l.Log(e)

	return e
}

// Fatal logs the message at the Fatal level.
func (l *Logger) Fatal(message string) *LogEntry {
	e := NewLogEntry(Fatal, message)

	l.Log(e)

	return e
}

// FatalF logs the formatted string at the Fatal level.
func (l *Logger) FatalF(format string, values ...interface{}) *LogEntry {
	e := NewLogEntry(Fatal, fmt.Sprintf(format, values...))

	l.Log(e)

	return e
}

// Close flushes and closes all streams in preparation for process exit.
// Use of this function should guarantee that all logs are persisted in their
// appropriate location.
func (l *Logger) Close() {
	for _, stream := range l.streams {
		stream.Close()
	}
}

// Println is added log.Logger compatibility and acts the same as Logln but
// at the Info level.
func (l *Logger) Println(val interface{}) {
	l.Logln(Info, fmt.Sprintf("%v", val))
}
