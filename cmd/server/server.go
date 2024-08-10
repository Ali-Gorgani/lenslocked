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
	"golang.org/x/oauth2"
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
	OAuthProviders map[string]*oauth2.Config
}

func loadEnvConfig() (config, error) {
	var cfg config

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	// Postgres
	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	if cfg.PSQL.Host == "" || cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no PSQL config provided")
	}

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
	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"
	// Server
	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	// OAuth
	cfg.OAuthProviders = make(map[string]*oauth2.Config)
	dropboxConfig := &oauth2.Config{
		ClientID:     os.Getenv("DROPBOX_APP_KEY"),
		ClientSecret: os.Getenv("DROPBOX_APP_SECRET"),
		Scopes: []string{
			"files.metadata.read",
			"files.content.read",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/oauth2/authorize",
			TokenURL: "https://api.dropboxapi.com/oauth2/token",
		},
	}
	cfg.OAuthProviders["dropbox"] = dropboxConfig

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	err = run(cfg)
	if err != nil {
		panic(err)
	}
}

func run(cfg config) error {
	// Setup the database
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		return err
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

	galleryService := &models.GalleryService{
		DB: db,
	}

	// Setup the middlewares
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	CsrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
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

	galleryC := controllers.Galleries{
		GalleryService: galleryService,
	}
	galleryC.Template.New = views.Must(views.ParseFS(templates.Fs, "galleries/new.gohtml", "tailwind.gohtml"))
	galleryC.Template.Edit = views.Must(views.ParseFS(templates.Fs, "galleries/edit.gohtml", "tailwind.gohtml"))
	galleryC.Template.Show = views.Must(views.ParseFS(templates.Fs, "galleries/show.gohtml", "tailwind.gohtml"))
	galleryC.Template.Index = views.Must(views.ParseFS(templates.Fs, "galleries/index.gohtml", "tailwind.gohtml"))

	oauthC := controllers.OAuth{
		ProviderConfigs: cfg.OAuthProviders,
	}

	// Setup the chi router.
	r := chi.NewRouter()
	r.Use(CsrfMw)
	r.Use(umw.SetUser)

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.Fs, "home.gohtml", "tailwind.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.Fs, "contact.gohtml", "tailwind.gohtml"))))
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.Fs, "faq.gohtml", "tailwind.gohtml"))))
	r.Get("/signup", userC.New)
	r.Post("/signup", userC.Create)
	r.Get("/signin", userC.SignIn)
	r.Post("/signin", userC.ProcessSignIn)
	r.Post("/signout", userC.ProcessSignOut)
	r.Get("/forgot-password", userC.ForgotPassword)
	r.Post("/forgot-password", userC.ProcessForgotPassword)
	r.Get("/reset-password", userC.ResetPassword)
	r.Post("/reset-password", userC.ProcessResetPassword)
	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleryC.Show)
		r.Get("/{id}/images/{filename}", galleryC.Image)
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/", galleryC.Index)
			r.Get("/new", galleryC.New)
			r.Post("/", galleryC.Create)
			r.Get("/{id}/edit", galleryC.Edit)
			r.Post("/{id}", galleryC.Update)
			r.Post("/{id}/delete", galleryC.Delete)
			r.Post("/{id}/images/{filename}/delete", galleryC.DeleteImage)
			r.Post("/{id}/images/upload", galleryC.UploadImage)
		})
	})
	r.Route("/oauth/{provider}", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/connect", oauthC.Connect)
		r.Get("/callback", oauthC.Callback)
	})
	assetsHandler := http.FileServer(http.Dir("assets"))
	r.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Println("Server is running on port ", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, r)
}
