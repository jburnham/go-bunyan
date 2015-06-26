package bunyan

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
	"text/template"
	"time"

	"github.com/beefsack/go-rate"
)

var (
	pendingEmails chan *LogEntry
	doneEmails    chan bool
)

// EmailStream defines the Stream and interface to output the logging data.
type EmailStream struct {
	*Stream
	recipient   string
	mailServer  string
	template    *template.Template
	rateLimiter *rate.RateLimiter
}

// NewEmailStream creates a new EmailStream with the specified logging
// level and and potential filters.
func NewEmailStream(minLogLevel LogLevel, filter StreamFilter, templateSource string, recipient, mailServer string, minimumInterval time.Duration) (result *EmailStream) {
	t, err := template.New("email").Parse(templateSource)

	if err != nil {
		panic(fmt.Sprintf("Unable to compile email template: %s", err.Error()))
	}

	pendingEmails = make(chan *LogEntry)
	doneEmails = make(chan bool)

	result = &EmailStream{
		Stream: &Stream{
			MinLogLevel: minLogLevel,
			Filter:      filter,
		},
		recipient:   recipient,
		mailServer:  mailServer,
		template:    t,
		rateLimiter: rate.New(1, minimumInterval),
	}

	// run the goroutine that processes all incoming logs to be emailed
	go dequeueEmails(result)
	return
}

func enqueueEmail(logEntry *LogEntry) {
	pendingEmails <- logEntry
}

func dequeueEmails(stream *EmailStream) {
	for {
		logEntry, more := <-pendingEmails
		if more {
			publishEmails(stream, logEntry)
		} else {
			doneEmails <- true
			return
		}
	}
}

func publishEmails(stream *EmailStream, l *LogEntry) {
	encodeRFC2047 := func(String string) string {
		addr := mail.Address{Name: String, Address: ""}
		return strings.Trim(addr.String(), " <>")
	}

	if stream.shouldPublish(l) {
		var output bytes.Buffer
		err := stream.template.ExecuteTemplate(&output, "email", l)

		if err != nil {
			println(fmt.Sprintf("Error compiling exception template: %s", err))
		}

		header := make(map[string]string)
		header["From"] = "Telemetry API <noreply@telemetryapp.com>"
		header["To"] = stream.recipient
		header["Subject"] = encodeRFC2047("Telemetry API Exception Report")
		header["MIME-Version"] = "1.0"
		header["Content-Type"] = "text/plain; charset=\"utf-8\""
		header["Content-Transfer-Encoding"] = "base64"

		message := ""
		for k, v := range header {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + base64.StdEncoding.EncodeToString(output.Bytes())

		c, err := smtp.Dial(stream.mailServer)

		println("dialed")
		if err != nil {
			println(fmt.Sprintf("Error connecting to SMTP server: %s", err))
			return
		}

		if c != nil {
			defer c.Close()
		}

		c.Mail("noreply@telemetryapp.com")
		c.Rcpt(stream.recipient)

		wc, err := c.Data()

		if err != nil {
			println(fmt.Sprintf("Error streaming to SMTP server: %s", err))
			return
		}

		if wc != nil {
			defer wc.Close()
		}

		_, err = bytes.NewBuffer([]byte(message)).WriteTo(wc)

		if err != nil {
			println(fmt.Sprintf("Error writing to SMTP server: %s", err))
			return
		}
	}
}

// Publish writes the logging data to the stream.
func (s *EmailStream) Publish(l *LogEntry) {
	if ok, _ := s.rateLimiter.Try(); !ok {
		// No more than 1 e-mail for each period!

		return
	}

	enqueueEmail(l)

}

// Close flushes and closes any pending writes to the log stream before shutdown.
func (s *EmailStream) Close() {
	close(pendingEmails)
	<-doneEmails
}
