package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func executeTemplate(w http.ResponseWriter, filePath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := template.ParseFiles(filePath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, filePath)
}

func contactFunc(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, filePath)
}

func faqFunc(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join("templates", "faq.gohtml")
	executeTemplate(w, filePath)
}

func main() {
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/contact", contactFunc)
	r.Get("/faq", faqFunc)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
