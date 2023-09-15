package main

import (
	"fmt"
	"net/http"

	"github.com/dmcclung/pixelparade/controllers"
	"github.com/dmcclung/pixelparade/db"
	"github.com/dmcclung/pixelparade/models"
	"github.com/dmcclung/pixelparade/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	db, err := db.DefaultPostgresConfig.Open()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", controllers.Static(
		views.Must(views.Parse("home.gohtml", "tailwind.gohtml")),
	))

	r.Get("/contact", controllers.Static(
		views.Must(views.Parse("contact.gohtml", "tailwind.gohtml")),
	))

	r.Get("/faq", controllers.Faq(
		views.Must(views.Parse("faq.gohtml", "tailwind.gohtml")),
	))

	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	userController := controllers.User{
		Templates: controllers.UserTemplates{
			Signup: views.Must(views.Parse("signup.gohtml", "tailwind.gohtml")),
			Signin: views.Must(views.Parse("signin.gohtml", "tailwind.gohtml")),
		},
		UserService: &userService,
		SessionService: &sessionService,
	}
	r.Get("/signup", userController.GetSignup)
	r.Post("/signup", userController.PostSignup)
	r.Get("/signin", userController.GetSignin)
	r.Post("/signin", userController.PostSignin)

	r.Get("/me", userController.CurrentUser)

	galleryController := controllers.Gallery{
		Templates: struct{Get controllers.Template}{
			Get: views.Must(views.Parse("gallery.gohtml", "tailwind.gohtml")),
		},
	} 
	r.Get("/gallery/{id}", galleryController.Get)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// TODO: Fix this before deploying
		csrf.Secure(false),
	)

	fmt.Println("Starting the server on :3000...")
	err = http.ListenAndServe(":3000", csrfMw(r))
	if err != nil {
		panic(err)
	}
}
