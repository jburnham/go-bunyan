package bunyan

import (
	"os"
)

// StdoutStream defines the Stream and interface to output the logging data.
type StdoutStream struct {
	*Stream
}

// NewStdoutStream creates a new StdoutStream with the specified logging
// level and and potential filters.
func NewStdoutStream(minLogLevel LogLevel, filter StreamFilter) *StdoutStream {
	return &StdoutStream{
		&Stream{
			MinLogLevel: minLogLevel,
			Filter:      filter,
		},
	}
}

// Publish writes the logging data to the stream.
func (s *StdoutStream) Publish(l *LogEntry) {
	if s.shouldPublish(l) {
		os.Stdout.WriteString(l.String())
		os.Stdout.WriteString("\n")
	}
}

// Close flushes and closes any pending writes to the log stream before shutdown.
func (s *StdoutStream) Close() {
	os.Stdout.Sync()
}
