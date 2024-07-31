package main

import (
	"fmt"
	"net/http"

	"github.com/Ali-Gorgani/lenslocked/controllers"
	"github.com/Ali-Gorgani/lenslocked/migrations"
	"github.com/Ali-Gorgani/lenslocked/models"
	"github.com/Ali-Gorgani/lenslocked/templates"
	"github.com/Ali-Gorgani/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	// Setup the database
	cfg := models.DefultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup the services
	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	// Setup the middlewares
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	CsrfAuthKey := "aB1cD2eF3gH4iJ5kL6mN7oP8qR9sT0uV"
	CsrfMw := csrf.Protect(
		[]byte(CsrfAuthKey),
		csrf.Secure(false),
	)

	// Setup the controllers
	userC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	userC.Template.New = views.Must(views.ParseFS(templates.Fs, "signup.gohtml", "tailwind.gohtml"))
	userC.Template.SignIn = views.Must(views.ParseFS(templates.Fs, "signin.gohtml", "tailwind.gohtml"))

	// Setup the chi router.
	r := chi.NewRouter()
	r.Use(CsrfMw)
	r.Use(umw.SetUser)

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.Fs, "home.gohtml", "tailwind.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.Fs, "contact.gohtml", "tailwind.gohtml"))))
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.Fs, "faq.gohtml", "tailwind.gohtml"))))
	r.Get("/signup", userC.New)
	r.Post("/users", userC.Create)
	r.Get("/signin", userC.SignIn)
	r.Post("/signin", userC.ProcessSignIn)
	r.Post("/signout", userC.ProcessSignOut)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", userC.CurrentUser)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
