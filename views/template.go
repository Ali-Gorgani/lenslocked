package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/Ali-Gorgani/lenslocked/context"
	"github.com/Ali-Gorgani/lenslocked/models"
	"github.com/gorilla/csrf"
)

type public interface {
	Public() string
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, tplToExecute string, patterns ...string) (Template, error) {
	tpl := template.New(path.Base(tplToExecute))
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField: not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser: not implemented")
			},
			"errors": func() []string {
				return nil
			},
			"equal": func(a, b interface{}) bool {
				return false
			},
		},
	)
	patterns = append(patterns, tplToExecute)
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{htmlTpl: tpl}, nil
}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("Cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}
	errMsgs := errMessages(errs...)
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return errMsgs
			},
			"equal": func(a, b interface{}) bool {
				return a == b
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func errMessages(errs ...error) []string {
	var publicErrs []string
	for _, err := range errs {
		var publicErr public
		if errors.As(err, &publicErr) {
			publicErrs = append(publicErrs, publicErr.Public())
		} else {
			fmt.Println(err)
			publicErrs = append(publicErrs, "Something went wrong. Please try again.")
		}
	}
	return publicErrs
}
