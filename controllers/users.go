package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ali-Gorgani/lenslocked/models"
)

type Users struct {
	Template struct {
		New    Template
		SignIn Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Template.New.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	user, err := u.UserService.Create(email, password)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user. Try again?", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "User created: %+v", user)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Template.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	if data.Email == "" || data.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		log.Printf("Error authenticating user: %v", err)
		http.Error(w, "Error authenticating user. Try again?", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "email",
		Value:    user.Email,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	fmt.Fprintf(w, "User authenticated: %+v", user)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("email")
	if err != nil {
		fmt.Fprint(w, "The email cookie could not be read.")
		return
	}
	fmt.Fprintf(w, "The current user is: %s", email.Value)
}
