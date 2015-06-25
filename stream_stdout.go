package bunyan

import (
	"os"
)

type StdoutStream struct {
	*Stream
}

func NewStdoutStream(minLogLevel LogLevel, filter StreamFilter) *StdoutStream {
	return &StdoutStream{
		&Stream{
			MinLogLevel: minLogLevel,
			Filter:      filter,
		},
	}
}

func (s *StdoutStream) Publish(l *LogEntry) {
	if s.shouldPublish(l) {
		os.Stdout.WriteString(l.String())
		os.Stdout.WriteString("\n")
	}
}

func (s *StdoutStream) Flushable() bool {
	return true
}

func (s *StdoutStream) Flush() {
	os.Stdout.Sync()
}
