package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"

	"github.com/go-mail/mail/v2"
)

//go:embed templates
var templateFS embed.FS

type Mailer struct {
	Dialer         *mail.Dialer
	Sender         string
	Domain         string
	Host           string
	Port           int
	Username       string
	Password       string
	Encryption     string
	EmailMessage   *EmailMessage
	MailerChan     chan *EmailMessage
	MailerDoneChan chan bool
	ErrorChan      chan error
	ToAddress      string
	FromAddress    string
}

type EmailMessage struct {
	From          string
	FromName      string
	To            string
	Subject       string
	HTMLBody      string
	PlainTextbody string
	Attachments   []string
}

func New() *Mailer {
	return &Mailer{}
}

func (m *Mailer) ListenForMail(errorChan chan error, messageChan chan *EmailMessage, mailerDoneChan chan bool) error {
	for {
		select {
		case msg := <-messageChan:
			log.Printf("Sending email to %s - from %s\n", msg.To, msg.From)
			go m.SendMail(msg, mailerDoneChan, errorChan)
		case err := <-errorChan:
			return err

		case <-mailerDoneChan:
			fmt.Println("mailer done")
			return nil

		}
	}
}

func (m *Mailer) SendMail(emailMessage *EmailMessage, doneChan chan bool, errorChan chan error) {

	go func() {

		tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+"email-html.page.tmpl")

		if err != nil {
			errorChan <- err
		}

		subject := new(bytes.Buffer)

		err = tmpl.ExecuteTemplate(subject, "subject", nil)

		if err != nil {
			errorChan <- err
			return
		}

		plainText := new(bytes.Buffer)

		err = tmpl.ExecuteTemplate(plainText, "plainText", emailMessage)

		if err != nil {
			errorChan <- err
			return
		}

		htmlBody := new(bytes.Buffer)

		err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", emailMessage)

		if err != nil {
			errorChan <- err
			return
		}

		msg := mail.NewMessage()

		msg.SetHeader("From", emailMessage.From)
		msg.SetHeader("To", emailMessage.To)
		msg.SetHeader("Subject", emailMessage.Subject)
		msg.SetBody("text/plain", plainText.String())
		msg.AddAlternative("text/html", htmlBody.String())

		err = m.Dialer.DialAndSend(msg)

		if err != nil {
			errorChan <- err
		}

		doneChan <- true

	}()

}

func (m *Mailer) SendPlainText(emailMessage *EmailMessage) error {
	return nil
}
