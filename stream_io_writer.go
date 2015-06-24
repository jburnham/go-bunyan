package bunyan

import (
	"io"
	"log"
)

type IOWriterStream struct {
	*Stream
	writer io.Writer
}

func NewIOWriterStream(w io.Writer, minLogLevel LogLevel, filter StreamFilter) *IOWriterStream {
	log.Printf("writer: %T", w)
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
		log.Printf("Publish %T", s.writer)
		s.writer.Write([]byte(l.String()))
		s.writer.Write([]byte("\n"))
	}
}
