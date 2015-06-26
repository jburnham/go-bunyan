package bunyan

import (
	"os"
)

// FileStream defines the Stream and interface to output the logging data.
type FileStream struct {
	*Stream
	outputFile *os.File
}

// NewFileStream creates a new FileStream with the specified logging
// level and and potential filters.
func NewFileStream(minLogLevel LogLevel, filter StreamFilter, path string) (result *FileStream, err error) {
	outputFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if err != nil {
		return
	}

	result = &FileStream{
		Stream: &Stream{
			MinLogLevel: minLogLevel,
			Filter:      filter,
		},
		outputFile: outputFile,
	}

	return
}

// Publish writes the logging data to the stream.
func (s *FileStream) Publish(l *LogEntry) {
	if s.shouldPublish(l) {
		s.outputFile.WriteString(l.String())
		s.outputFile.WriteString("\n")
	}
}

// Close flushes and closes any pending writes to the log stream before shutdown.
func (s *FileStream) Close() {
	s.outputFile.Sync()
}
