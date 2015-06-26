package bunyan

// StreamInterface is the interface that all Streams must implement.
// This allows the logger to publish all logs and close them when the
// program is complete.
type StreamInterface interface {
	Publish(*LogEntry)
	Close()
}

type flushableStream interface {
	Flush() error
}

// StreamFilter defines a filter function that determines if a LogEntry
// should be published.
type StreamFilter func(*LogEntry) bool

// Stream defines the stream with a minimum log level and any filters.
type Stream struct {
	MinLogLevel LogLevel
	Filter      StreamFilter
}

func (s *Stream) shouldPublish(l *LogEntry) bool {
	if l.Level < s.MinLogLevel {
		return false
	}

	if s.Filter != nil {
		return s.Filter(l)
	}

	return true
}
