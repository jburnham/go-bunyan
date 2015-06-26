package bunyan

import (
	"encoding/json"
	"errors"

	"github.com/bugsnag/bugsnag-go"
)

var (
	pendingReports chan *LogEntry
	doneBugsnag    chan bool
)

// BugsnagStream defines the Stream and interface to output the logging data.
type BugsnagStream struct {
	*Stream
}

// NewBugsnagStream creates a new BugsnagStream with the specified logging
// level and and potential filters.
func NewBugsnagStream(minLogLevel LogLevel, filter StreamFilter) (result *BugsnagStream) {

	pendingReports = make(chan *LogEntry)
	doneBugsnag = make(chan bool)

	go dequeueBugsnags()
	return &BugsnagStream{
		Stream: &Stream{
			MinLogLevel: minLogLevel,
			Filter:      filter,
		},
	}
}

func enqueueBugsnag(l *LogEntry) {
	pendingReports <- l
}

func dequeueBugsnags() {
	for {
		logEntry, more := <-pendingReports
		if more {
			publishBugsnags(logEntry)
		} else {
			doneBugsnag <- true
			return
		}
	}
}

func publishBugsnags(l *LogEntry) {
	b, err := json.Marshal(l)

	if err != nil {
		return
	}

	data := map[string]interface{}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return
	}

	err = l.Error

	if err == nil {
		err = errors.New(l.Message)
	}

	metadata := &bugsnag.MetaData{}
	metadata.AddStruct("Log", data)

	bugsnag.Notify(err, *metadata)
}

// Publish writes the logging data to the stream.
func (s *BugsnagStream) Publish(l *LogEntry) {
	if s.shouldPublish(l) {
		enqueueBugsnag(l)
	}
}

// Close flushes and closes any pending writes to the log stream before shutdown.
func (s *BugsnagStream) Close() {
	close(pendingReports)
	<-doneBugsnag
}
