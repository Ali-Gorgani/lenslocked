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
	UserService    *models.UserService
	SessionService *models.SessionService
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
		http.Error(w, "Something went wrong in creating user. Try again?", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Something went wrong in creating user. Try again?", http.StatusInternalServerError)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
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

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Something went wrong in creating user. Try again?", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	user, err := u.SessionService.User(token)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "Current user: %v", user)
}
