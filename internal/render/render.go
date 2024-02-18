package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/elkcityhazard/am-form/internal/config"
	"github.com/elkcityhazard/am-form/internal/models"
	"github.com/elkcityhazard/am-form/templates"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

var pathToTemplates = "./templates/templates"

var templateFuncs = template.FuncMap{
	"sanitizeHTML": template.HTMLEscapeString,
}

func NewRenderer(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template
	if app.IsProduction {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]

	if !ok {
		return fmt.Errorf("the template %s does not exist", tmpl)
	}

	td = AddDefaultTemplateData(td, r)

	buf := new(bytes.Buffer)

	err := t.Execute(buf, td)

	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	return err
}

func CreateTemplateCache() (map[string]*template.Template, error) {

	templateSet := map[string]*template.Template{}

	if app.IsProduction {
		files := templates.GetTemplates()

		pages, err := files.ReadDir("templates/pages")

		if err != nil {
			log.Fatalln(err)
		}
		for _, page := range pages {

			tmpl, err := template.New(page.Name()).Funcs(templateFuncs).ParseFS(files, fmt.Sprintf("templates/pages/%s", page.Name()))

			if err != nil {
				return nil, err

			}

			layoutMatches, err := files.ReadDir("templates/layouts")

			if err != nil {
				return nil, err
			}

			if len(layoutMatches) > 0 {
				for _, layout := range layoutMatches {
					tmpl.ParseFS(files, fmt.Sprintf("templates/layouts/%s", layout.Name()))
				}
			}

			parialsMatches, err := files.ReadDir("templates/partials")

			if err != nil {
				return nil, err
			}

			if len(parialsMatches) > 0 {
				for _, partials := range parialsMatches {
					tmpl.ParseFS(files, fmt.Sprintf("templates/partials/%s", filepath.Base(partials.Name())))
				}
			}

			templateSet[page.Name()] = tmpl
		}

		return templateSet, nil
	} else {
		pages, _ := filepath.Glob(pathToTemplates + "/pages/*.page.tmpl")

		for _, page := range pages {

			name := filepath.Base(page)

			tmpl, err := template.New(name).Funcs(templateFuncs).ParseFiles(page)

			if err != nil {
				return nil, err
			}

			layoutMatches, err := filepath.Glob(pathToTemplates + "/layouts/*.layout.tmpl")

			if err != nil {
				return nil, err
			}

			if len(layoutMatches) > 0 {
				for _, layout := range layoutMatches {
					_, err := tmpl.ParseFiles(layout)

					if err != nil {
						return nil, err
					}
				}
			}

			partialsMatches, err := filepath.Glob(pathToTemplates + "/partials/*.partial.tmpl")

			if err != nil {
				return nil, err
			}

			if len(partialsMatches) > 0 {
				for _, partials := range partialsMatches {
					_, err := tmpl.ParseFiles(partials)

					if err != nil {
						return nil, err
					}
				}
			}

			templateSet[name] = tmpl
		}

		return templateSet, nil
	}

}

func AddDefaultTemplateData(td *models.TemplateData, r *http.Request) *models.TemplateData {

	flash := app.SessionManager.PopString(r.Context(), "flash")
	td.Flash = flash

	error := app.SessionManager.PopString(r.Context(), "error")
	td.Error = error

	td.CSRFToken = nosurf.Token(r)
	return td
}
