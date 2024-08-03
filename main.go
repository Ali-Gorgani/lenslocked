package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Ali-Gorgani/lenslocked/controllers"
	"github.com/Ali-Gorgani/lenslocked/migrations"
	"github.com/Ali-Gorgani/lenslocked/models"
	"github.com/Ali-Gorgani/lenslocked/templates"
	"github.com/Ali-Gorgani/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	// Postgres
	cfg.PSQL = models.DefultPostgresConfig()

	// SMTP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Port = port
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	// CSRF
	cfg.CSRF.Key = "aB1cD2eF3gH4iJ5kL6mN7oP8qR9sT0uV"
	cfg.CSRF.Secure = false

	// Server
	cfg.Server.Address = ":8080"

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// Setup the database
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup the services
	userService := &models.UserService{
		DB: db,
	}

	sessionService := &models.SessionService{
		DB: db,
	}

	passwordResetService := &models.PasswordResetService{
		DB: db,
	}

	emailService := models.NewEmailService(cfg.SMTP)

	// Setup the middlewares
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	CsrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
	)

	// Setup the controllers
	userC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: passwordResetService,
		EmailService:         emailService,
	}
	userC.Template.New = views.Must(views.ParseFS(templates.Fs, "signup.gohtml", "tailwind.gohtml"))
	userC.Template.SignIn = views.Must(views.ParseFS(templates.Fs, "signin.gohtml", "tailwind.gohtml"))
	userC.Template.ForgotPassword = views.Must(views.ParseFS(templates.Fs, "forgot-password.gohtml", "tailwind.gohtml"))
	userC.Template.CheckYourEmail = views.Must(views.ParseFS(templates.Fs, "check-your-email.gohtml", "tailwind.gohtml"))
	userC.Template.ResetPassword = views.Must(views.ParseFS(templates.Fs, "reset-password.gohtml", "tailwind.gohtml"))

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
	r.Get("/forgot-password", userC.ForgotPassword)
	r.Post("/forgot-password", userC.ProcessForgotPassword)
	r.Get("/reset-password", userC.ResetPassword)
	r.Post("/reset-password", userC.ProcessResetPassword)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", userC.CurrentUser)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Println("Server is running on port ", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
