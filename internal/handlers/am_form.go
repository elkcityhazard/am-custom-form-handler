package handlers

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"net/http"

	"github.com/elkcityhazard/am-form/internal/forms"
	"github.com/elkcityhazard/am-form/internal/mailer"
	"github.com/elkcityhazard/am-form/internal/models"
	"github.com/elkcityhazard/am-form/internal/render"
	"github.com/elkcityhazard/am-form/internal/repository/dbrepo"
)

func HandleDisplayAMForm(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:

		dataMap := make(map[string]any)

		dataMap["Form"] = forms.New(nil)

		if err := render.RenderTemplate(w, r, "am-form.page.tmpl", &models.TemplateData{
			DataMap: dataMap,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:

		err := r.ParseForm()

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		form := forms.New(r.PostForm)

		if len(form.Values.Get("password")) > 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		username := form.Values.Get("username")

		if len(username) > 0 {

			err := errors.New("bad request")

			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		form.Required("firstName", "email", "message")

		form.MinLength("firstName", 10)

		form.IsEmail("email")

		if !form.Valid() {

			app.SessionManager.Put(r.Context(), "error", "Please check your form.")

			dataMap := make(map[string]any)

			dataMap["Form"] = form

			if err := render.RenderTemplate(w, r, "am-form.page.tmpl", &models.TemplateData{
				DataMap: dataMap,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		repo := dbrepo.GetDBRepo()

		msg := &models.Message{
			Email:   form.Values.Get("email"),
			Message: form.Values.Get("message"),
		}

		if err := repo.InsertMessage(msg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		app.SessionManager.Put(r.Context(), "flash", "Your message has been sent!")

		newMsg := &mailer.EmailMessage{
			To:            app.Mailer.ToAddress,
			From:          app.Mailer.FromAddress,
			FromName:      msg.Email,
			Subject:       fmt.Sprintf("Message from %s", msg.Email),
			PlainTextbody: template.HTMLEscapeString(msg.Message),
			HTMLBody:      html.EscapeString(msg.Message),
		}

		app.Mailer.MailerChan <- newMsg

		app.SessionManager.Put(r.Context(), "email_message", newMsg)

		http.Redirect(w, r.WithContext(r.Context()), "/success", http.StatusSeeOther)

	}

}
