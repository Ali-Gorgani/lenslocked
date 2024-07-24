package main

import (
	"fmt"
	"net/http"

	"github.com/Ali-Gorgani/lenslocked/controllers"
	"github.com/Ali-Gorgani/lenslocked/templates"
	"github.com/Ali-Gorgani/lenslocked/views"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.Fs, "home.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.Fs, "contact.gohtml"))))
	r.Get("/faq", controllers.StaticHandler(views.Must(views.ParseFS(templates.Fs, "faq.gohtml"))))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
