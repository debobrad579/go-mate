package handlers

import (
	"net/http"

	"github.com/debobrad579/chessgo/internal/auth"
	"github.com/debobrad579/chessgo/internal/database"
	"github.com/lib/pq"
)

type registerData struct {
	Fields struct {
		Name            string
		Email           string
		Password        string
		ConfirmPassword string
	}
	Errors struct {
		Name            string
		Email           string
		Password        string
		ConfirmPassword string
	}
	AuthError string
}

func (cfg *Config) RegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	var data registerData
	data.Fields.Name = name
	data.Fields.Email = email
	data.Fields.Password = password
	data.Fields.ConfirmPassword = confirmPassword

	if name == "" {
		data.Errors.Name = "Required"
	}

	if email == "" {
		data.Errors.Email = "Required"
	} else if !isEmailValid(email) {
		data.Errors.Email = "Invalid email address"
	}

	if password == "" {
		data.Errors.Password = "Required"
	} else if len(password) < 8 {
		data.Errors.Password = "Must be at least 8 characters long."
	}

	if confirmPassword == "" {
		data.Errors.ConfirmPassword = "Required"
	} else if confirmPassword != password {
		data.Errors.ConfirmPassword = "Passwords do not match."
	}

	if data.Errors.Name != "" || data.Errors.Email != "" || data.Errors.Password != "" || data.Errors.ConfirmPassword != "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		RenderTemplate(w, "register.html", data)
		return
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		registerError(w, data, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{Email: email, Name: name, HashedPassword: hashedPassword})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				w.WriteHeader(http.StatusConflict)
				data.Errors.Email = "Email already in use"
				RenderTemplate(w, "register.html", data)
				return
			}
		}

		registerError(w, data, http.StatusInternalServerError, "Failed to create user")
		return
	}

	if err := cfg.login(w, r, user.ID); err != nil {
		registerError(w, data, http.StatusInternalServerError, "Failed to log in")
		return
	}

	http.Redirect(w, r, "/app", http.StatusSeeOther)
}

func registerError(w http.ResponseWriter, data registerData, code int, message string) {
	data.AuthError = message
	w.WriteHeader(code)
	RenderTemplate(w, "register.html", loginData{AuthError: message})
}
