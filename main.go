package main

import (
	"fmt"
	"net/http"

	"github.com/dmcclung/pixelparade/controllers"
	"github.com/dmcclung/pixelparade/models"
	"github.com/dmcclung/pixelparade/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	cfg.PSQL = models.DefaultPostgresConfig
	smtpConfig, err := models.GetEmailConfig()
	if err != nil {
		return cfg, err
	}
	cfg.SMTP = smtpConfig
	cfg.CSRF.Key = "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	cfg.CSRF.Secure = false

	cfg.Server.Address = ":3000"

	return cfg, nil
}

func main() {
	config, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	db, err := models.Open(config.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.Migrate(db)
	if err != nil {
		panic(err)
	}

	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	emailService, err := models.GetEmailService(config.SMTP)
	if err != nil {
		panic(err)
	}

	passwordResetService := models.PasswordResetService{
		DB: db,
	}

	csrfMw := csrf.Protect(
		[]byte(config.CSRF.Key),
		// TODO: Fix this before deploying
		csrf.Secure(config.CSRF.Secure),
	)

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(csrfMw)
	r.Use(umw.SetUser)

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
			Signup: views.Must(views.Parse("signup.gohtml", "tailwind.gohtml")),
			Signin: views.Must(views.Parse("signin.gohtml", "tailwind.gohtml")),
			Me:     views.Must(views.Parse("me.gohtml", "tailwind.gohtml")),
			Forgot: views.Must(views.Parse("forgot.gohtml", "tailwind.gohtml")),
			CheckEmail: views.Must(views.Parse("checkemail.gohtml", "tailwind.gohtml")),
		},
		UserService:    &userService,
		SessionService: &sessionService,
		PasswordResetService: &passwordResetService,
		EmailService: emailService,
	}
	r.Get("/signup", userController.GetSignup)
	r.Post("/signup", userController.PostSignup)
	r.Get("/signin", userController.GetSignin)
	r.Post("/signin", userController.PostSignin)
	r.Post("/signout", userController.PostSignout)
	r.Get("/forgot", userController.ForgotPassword)
	r.Post("/forgot", userController.PostForgotPassword)
	r.Get("/checkemail", userController.CheckEmail)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", userController.CurrentUser)
	})

	galleryController := controllers.Gallery{
		Templates: struct{ Get controllers.Template }{
			Get: views.Must(views.Parse("gallery.gohtml", "tailwind.gohtml")),
		},
	}
	r.Get("/gallery/{id}", galleryController.Get)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")
	err = http.ListenAndServe(config.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
