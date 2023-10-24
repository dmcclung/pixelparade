package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dmcclung/pixelparade/controllers"
	"github.com/dmcclung/pixelparade/models"
	"github.com/dmcclung/pixelparade/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v4/stdlib"
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
	err := godotenv.Load()
	if err != nil {
		return cfg, fmt.Errorf("loading env: %w", err)
	}

	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}

	if cfg.PSQL.Host == "" || cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no postgres configuration found")
	}

	smtpConfig, err := models.GetEmailConfig()
	if err != nil {
		return cfg, err
	}
	cfg.SMTP = smtpConfig
	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return cfg, nil
}

func run(cfg config) error {
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = models.Migrate(db)
	if err != nil {
		return err
	}

	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	galleryService := models.GalleryService{
		DB: db,
	}

	emailService, err := models.GetEmailService(cfg.SMTP)
	if err != nil {
		return err
	}

	passwordResetService := models.PasswordResetService{
		DB: db,
	}

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		// TODO: Fix this before deploying
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(csrfMw)
	r.Use(umw.SetUser)

	assetsHandler := http.FileServer(http.Dir("assets"))
	r.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	r.Get("/", controllers.Static(
		views.Must(views.Parse("home.gohtml", "tailwind.gohtml")),
	))

	r.Get("/contact", controllers.Static(
		views.Must(views.Parse("contact.gohtml", "tailwind.gohtml")),
	))

	r.Get("/faq", controllers.Faq(
		views.Must(views.Parse("faq.gohtml", "tailwind.gohtml")),
	))

	userController := controllers.User{
		Templates: controllers.UserTemplates{
			SignUp:         views.Must(views.Parse("sign-up.gohtml", "tailwind.gohtml")),
			SignIn:         views.Must(views.Parse("sign-in.gohtml", "tailwind.gohtml")),
			Me:             views.Must(views.Parse("me.gohtml", "tailwind.gohtml")),
			ForgotPassword: views.Must(views.Parse("forgot-password.gohtml", "tailwind.gohtml")),
			CheckEmail:     views.Must(views.Parse("check-email.gohtml", "tailwind.gohtml")),
			ResetPassword:  views.Must(views.Parse("reset-password.gohtml", "tailwind.gohtml")),
		},
		UserService:          &userService,
		SessionService:       &sessionService,
		PasswordResetService: &passwordResetService,
		EmailService:         emailService,
	}
	r.Get("/signup", userController.SignUp)
	r.Post("/signup", userController.ProcessSignUp)
	r.Get("/signin", userController.SignIn)
	r.Post("/signin", userController.ProcessSignIn)
	r.Post("/signout", userController.ProcessSignOut)
	r.Get("/forgot-password", userController.ForgotPassword)
	r.Post("/forgot-password", userController.ProcessForgotPassword)
	r.Get("/check-email", userController.CheckEmail)
	r.Get("/reset-password", userController.ResetPassword)
	r.Post("/reset-password", userController.ProcessResetPassword)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", userController.CurrentUser)
	})

	galleryController := controllers.Gallery{
		Templates: controllers.GalleryTemplates{
			New:   views.Must(views.Parse("galleries/new.gohtml", "tailwind.gohtml")),
			Edit:  views.Must(views.Parse("galleries/edit.gohtml", "tailwind.gohtml")),
			Show:  views.Must(views.Parse("galleries/show.gohtml", "tailwind.gohtml")),
			Index: views.Must(views.Parse("galleries/index.gohtml", "tailwind.gohtml")),
		},
		GalleryService: &galleryService,
	}

	r.Route("/galleries", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/", galleryController.Index)
			r.Get("/new", galleryController.New)
			r.Post("/", galleryController.Create)
			r.Get("/{id}/edit", galleryController.Edit)
			r.Post("/{id}", galleryController.Update)
			r.Post("/{id}/delete", galleryController.Delete)
			r.Post("/{id}/{filename}/delete", galleryController.DeleteImage)
			r.Post("/{id}/images", galleryController.CreateImage)
		})
		r.Get("/{id}", galleryController.Show)
		r.Get("/{id}/{filename}", galleryController.Image)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")
	return http.ListenAndServe(cfg.Server.Address, r)
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