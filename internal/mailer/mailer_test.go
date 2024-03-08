package mailer

import (
	"embed"
	"errors"
	"html/template"
	"testing"
	"time"

	"github.com/go-mail/mail/v2"
)

type MockTemplateLoader struct {
	FS embed.FS
}

func NewMockTemplateLoader() *MockTemplateLoader {
	return &MockTemplateLoader{}
}

func (m *MockTemplateLoader) LoadTemplate(name string) (*template.Template, error) {
	// Simulate an error when loading the template

	return nil, errors.New("my error")
}

func Test_New(t *testing.T) {
	var mailer = New()
	if mailer == nil {
		t.Error("mailer is nil")
	}
}

func Test_ListenForMail(t *testing.T) {

	var mailer = New()

	mailer.TemplateLoader = NewTemplateLoader()

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

	testErr := errors.New("test err")

	errorChan <- testErr

	<-doneChan

	close(doneChan)
	close(errorChan)
	close(emailMessageChan)
}

func Test_SendMail(t *testing.T) {
	// Create a new Mailer instance
	var mailer = &Mailer{}

	// Set up the mock template loader to return an error
	mailer.TemplateLoader = &MockTemplateLoader{}

	// Set up the dialer and other necessary fields for the mailer
	mailer.Dialer = mail.NewDialer("localhost", 1025, "test", "test")
	mailer.Dialer.Timeout = 5 * time.Second

	// Prepare the email message
	emailMessage := &EmailMessage{
		To:            "test@example.com",
		From:          "test@example.com",
		FromName:      "test",
		Subject:       "Test message",
		HTMLBody:      "test",
		PlainTextbody: "Test",
	}

	// Create channels for done, error, and email message
	doneChan := make(chan bool)
	errorChan := make(chan error)

	// Call SendMail in a goroutine
	go mailer.SendMail(emailMessage, doneChan, errorChan)

	// Wait for an error to be received on the errorChan
	select {
	case err := <-errorChan:
		if err == nil {
			t.Error("Expected an error but received nil")
		} else {
			t.Log("Received expected error:", err)
		}
	case <-doneChan:
		t.Error("Expected an error but received done signal")
	case <-time.After(time.Second):
		t.Error("Timeout waiting for error")
	}

	// Cleanup
	close(doneChan)
	close(errorChan)
}
