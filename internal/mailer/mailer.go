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

type TemplateLoader interface {
	LoadTemplate(name string) (*template.Template, error)
}

type RealTemplateLoader struct {
	FS embed.FS
}

func NewTemplateLoader() *RealTemplateLoader {
	return &RealTemplateLoader{
		FS: templateFS,
	}
}

func (r *RealTemplateLoader) LoadTemplate(name string) (*template.Template, error) {
	tmpl, err := template.ParseFS(r.FS, "templates/"+name)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return tmpl, nil
}

type Mailer struct {
	TemplateLoader TemplateLoader
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
	return &Mailer{
		TemplateLoader: NewTemplateLoader(),
	}
}

func (m *Mailer) ListenForMail(errorChan chan error, messageChan chan *EmailMessage, mailerDoneChan chan bool) {
	for {
		select {
		case msg := <-messageChan:
			log.Printf("Sending email to %s - from %s\n", msg.To, msg.From)
			go m.SendMail(msg, mailerDoneChan, errorChan)
		case err := <-errorChan:
			log.Println(err)

		case <-mailerDoneChan:
			return

		}
	}

}

func (m *Mailer) SendMail(emailMessage *EmailMessage, doneChan chan bool, errorChan chan error) {
	defer func() {
		if r := recover(); r != nil {
			errorChan <- fmt.Errorf("panic in SendMail: %v", r)
		}
	}()

	tmpl, err := m.TemplateLoader.LoadTemplate("email-html.page.tmpl")
	if err != nil {
		errorChan <- fmt.Errorf("failed to load template: %w", err)
		return
	}

	subject, plainText, htmlBody, err := m.prepareEmailContent(tmpl, emailMessage)
	if err != nil {
		errorChan <- fmt.Errorf("failed to prepare email content: %w", err)
		return
	}

	msg := mail.NewMessage()
	msg.SetHeader("From", emailMessage.From)
	msg.SetHeader("To", emailMessage.To)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plainText)
	msg.AddAlternative("text/html", htmlBody)

	err = m.Dialer.DialAndSend(msg)
	if err != nil {
		errorChan <- fmt.Errorf("failed to send email: %w", err)
		return
	}

	doneChan <- true
}

func (m *Mailer) prepareEmailContent(tmpl *template.Template, emailMessage *EmailMessage) (string, string, string, error) {
	subject := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(subject, "subject", nil)
	if err != nil {
		return "", "", "", err
	}

	plainText := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainText, "plainText", emailMessage)
	if err != nil {
		return "", "", "", err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", emailMessage)
	if err != nil {
		return "", "", "", err
	}

	return subject.String(), plainText.String(), htmlBody.String(), nil
}

func (m *Mailer) SendPlainText(emailMessage *EmailMessage) error {
	return nil
}
