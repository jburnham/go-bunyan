package bunyan

import "io"

type IOWriterStream struct {
	*Stream
	writer io.Writer
}

func NewIOWriterStream(w io.Writer, minLogLevel LogLevel, filter StreamFilter) *IOWriterStream {
	return &IOWriterStream{
		&Stream{
			MinLogLevel: minLogLevel,
			Filter:      filter,
		},
		w,
	}
}

func (s *IOWriterStream) Publish(l *LogEntry) {
	if s.shouldPublish(l) {
		s.writer.Write([]byte(l.String()))
		s.writer.Write([]byte("\n"))
	}
}
