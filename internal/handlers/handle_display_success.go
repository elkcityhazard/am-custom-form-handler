package handlers

import (
	"net/http"

	"github.com/elkcityhazard/am-form/internal/mailer"
	"github.com/elkcityhazard/am-form/internal/models"
	"github.com/elkcityhazard/am-form/internal/render"
)

func HandleDisplaySuccess(w http.ResponseWriter, r *http.Request) {

	var emailMessage mailer.EmailMessage

	if app.SessionManager.Exists(r.Context(), "email_message") {
		emailMessage = app.SessionManager.Pop(r.Context(), "email_message").(mailer.EmailMessage)
	}

	dataMap := make(map[string]any)
	dataMap["EmailMessage"] = emailMessage

	if err := render.RenderTemplate(w, r, "success.page.tmpl", &models.TemplateData{
		DataMap: dataMap,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
