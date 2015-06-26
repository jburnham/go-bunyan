package bunyan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// LogLevel defines the logging level.
type LogLevel int

const (
	// Fatal is the most severe logging level, programs should not continue
	// running if there is an error of this level.
	Fatal LogLevel = 60
	// Error is the normal error logging, programs can probably continue running.
	Error = 50
	// Warn is a logging level that shows the log entry as maybe or maybe not
	// being a problem.
	Warn = 40
	// Info is a normal informational logging level.
	Info = 30
	// Debug is a normal verbose level of logging.
	Debug = 20
	// Trace is the highly verbose level of logging.
	Trace = 10
)

// Request contains the elements of an http request.
type Request struct {
	Method        string      `json:"method"`
	URL           string      `json:"url"`
	Headers       http.Header `json:"headers"`
	RemoteAddress string      `json:"remoteAddress"`
	Body          interface{} `json:"body,omitempty"`
}

// Response contains the elements of an http response.
type Response struct {
	StatusCode int         `json:"statusCode,omitempty"`
	Headers    http.Header `json:"headers,omitempty"`
	Body       interface{} `json:"body,omitempty"`
}

// LogEntry contains the elemts of a logging message.
type LogEntry struct {
	Data        interface{} `json:"data,omitempty"`
	Error       error       `json:"error,omitempty"`
	Hostname    string      `json:"hostname"`
	Level       LogLevel    `json:"level"`
	Message     string      `json:"msg"`
	Name        string      `json:"name"`
	ProcessID   int         `json:"pid"`
	Request     *Request    `json:"req,omitempty"`
	Response    *Response   `json:"res,omitempty"`
	StackTrace  string      `json:"trace,omitempty"`
	Time        time.Time   `json:"time"`
	CompletedIn string      `json:"completed_in,omitempty"`
	Version     int         `json:"v"`
}

func hostname() string {
	result, err := os.Hostname()

	if err != nil {
		panic(fmt.Sprintf("Error retrieving hostname %v", err))
	}

	return result
}

var logEntryTemplate = LogEntry{
	Version:   0,
	Hostname:  hostname(),
	ProcessID: os.Getpid(),
}

// NewLogEntry creates a new LogEntry with the specified log level and message.
func NewLogEntry(level LogLevel, message string) *LogEntry {
	result := logEntryTemplate

	result.Level = level
	result.Message = message
	result.Time = time.Now()

	return &result
}

// SetData sets the data for the LogEntry.
func (l *LogEntry) SetData(data interface{}) {
	l.Data = data
}

// SetRequest sets the Request field for the LogEntry.
func (l *LogEntry) SetRequest(r *http.Request) {
	l.Request = &Request{
		Method:        r.Method,
		URL:           r.URL.RequestURI(),
		Headers:       r.Header,
		RemoteAddress: r.RemoteAddr,
	}
}

// SetRequestBody sets the Request.Body field for the LogEntry.
func (l *LogEntry) SetRequestBody(body []byte) {
	if json.Unmarshal(body, &l.Request.Body) != nil {
		l.Request.Body = string(body)
	}
}

// SetResponseStatusCode sets the Response.StatusCode for the LogEntry.
func (l *LogEntry) SetResponseStatusCode(statusCode int) {
	if statusCode > 0 {
		l.Response.StatusCode = statusCode
	} else {
		l.Response.StatusCode = http.StatusOK
	}
}

// SetResponseBody sets the Response.Body for the LogEntry.
func (l *LogEntry) SetResponseBody(body []byte) {
	if json.Unmarshal(body, &l.Response.Body) != nil {
		l.Response.Body = string(body)
	}
}

// SetResponseError sets the Error for the LogEntry.
func (l *LogEntry) SetResponseError(err error) {
	l.Error = err
}

// SetCompletedIn sets the CompletedIn field for the LogEntry.
func (l *LogEntry) SetCompletedIn(completedIn string) {
	l.CompletedIn = completedIn
}

// SetStackTrace sets the StackTrace field for the LogEntry.
func (l *LogEntry) SetStackTrace(trace string) {
	l.StackTrace = trace
}

func (l *LogEntry) setLogger(logger *Logger) {
	l.Name = logger.name
}

func (l *LogEntry) String() string {
	result, _ := json.Marshal(l)

	return string(result)
}
