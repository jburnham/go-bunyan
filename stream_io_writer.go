package bunyan

type IOWriterStream struct {
	*Stream
	writer FlushableWriter
}

type FlushableWriter interface {
	Write(p []byte) (int, error)
	Flush() error
}

func NewIOWriterStream(w FlushableWriter, minLogLevel LogLevel, filter StreamFilter) *IOWriterStream {
	return &IOWriterStream{
		&Stream{
			MinLogLevel: minLogLevel,
			Filter:      filter,
			Flushable:   true,
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

func (s *IOWriterStream) Flushable() bool {
	return s.Stream.Flushable
}
func (s *IOWriterStream) Flush() {
	s.writer.Flush()
}
