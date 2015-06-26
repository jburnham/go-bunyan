package bunyan

import "io"

// IOWriterStream defines the Stream and interface to output the logging data.
type IOWriterStream struct {
	*Stream
	writer io.Writer
}

// NewIOWriterStream creates a new IOWriterStream with the specified logging
// level and and potential filters.
func NewIOWriterStream(w io.Writer, minLogLevel LogLevel, filter StreamFilter) *IOWriterStream {
	return &IOWriterStream{
		&Stream{
			MinLogLevel: minLogLevel,
			Filter:      filter,
		},
		w,
	}
}

// Publish writes the logging data to the stream.
func (s *IOWriterStream) Publish(l *LogEntry) {
	if s.shouldPublish(l) {
		s.writer.Write([]byte(l.String()))
		s.writer.Write([]byte("\n"))
	}
}

// Close flushes and closes any pending writes to the log stream before shutdown.
func (s *IOWriterStream) Close() {
	if flusher, ok := s.writer.(flushableStream); ok {
		flusher.Flush()
	}
}
