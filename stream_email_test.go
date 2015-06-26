package bunyan

import (
	"testing"
	"time"
)

func TestEmailStream(t *testing.T) {

	s1 := NewEmailStream(Info, nil, "Test {{ .Message }}", "testing@example.com", "localhost:25", 10*time.Second)
	l1 := NewLogger("Test", []StreamInterface{s1})
	defer l1.Close()
}
