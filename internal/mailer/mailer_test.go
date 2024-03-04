package mailer

import (
	"bytes"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-mail/mail/v2"
)

func Test_New(t *testing.T) {
	var mailer = New()
	if mailer == nil {
		t.Error("mailer is nil")
	}
}

func Test_ListenForMail(t *testing.T) {

	t.Log("staring listen for mail")

	tests := []struct {
		name     string
		host     int
		username string
		password string
		expected string
	}{
		{
			name:     "no host",
			host:     0,
			username: "test",
			password: "test",
			expected: "localhost",
		},
		{
			name:     "no username",
			host:     1025,
			username: "",
			password: "test",
			expected: "test",
		},
		{
			name:     "no password",
			host:     1025,
			username: "test",
			password: "",
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mailer = New()

			mailer.Dialer = mail.NewDialer("localhost", 1025, "test", "test")

			mailer.Dialer.Timeout = 5 * time.Second

			doneChan := make(chan bool)
			errorChan := make(chan error)
			emailMessageChan := make(chan *EmailMessage)

			go mailer.ListenForMail(errorChan, emailMessageChan, doneChan)

			emailMessage := &EmailMessage{
				To:            "test@example.com",
				From:          "test@example.com",
				FromName:      "test",
				Subject:       "Test message",
				HTMLBody:      "test",
				PlainTextbody: "Test",
			}

			t.Log("sending email", emailMessage)

			emailMessageChan <- emailMessage

			for {

				select {
				case err := <-errorChan:
					var buf bytes.Buffer

					log.SetOutput(&buf)

					log.Println(err)

					t.Log(buf.String())

					log.SetOutput(os.Stderr)

					if buf.String() == "" {

						t.Error("no error")
					}

				case <-doneChan:
					t.Log("mailer done")
					return
				}
			}
		})
	}

}
