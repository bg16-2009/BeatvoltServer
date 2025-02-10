package main

import (
	"net/http"
	"regexp"

	"github.com/bg16_2009/beatvolt_server/models"
	"github.com/go-chi/chi/v5"
)

func registerGroup(r chi.Router) {
	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate("register", w, map[string]string{
			"Error": r.URL.Query().Get("error"),
		})
	})
	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		if password != confirmPassword {
			http.Redirect(w, r, "/register?error=Passwords do not match", http.StatusSeeOther)
			return
		}
		var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`) // Basic email regex
		if !emailRegex.MatchString(email) {
			http.Redirect(w, r, "/register?error=Invalid email", http.StatusSeeOther)
			return
		}
		user := models.User{}
		if err := db.Where("email = ?", email).First(&user).Error; err == nil {
			http.Redirect(w, r, "/register?error=Email is already used", http.StatusSeeOther)
			return
		}
		if err := db.Where("username = ?", username).First(&user).Error; err == nil {
			http.Redirect(w, r, "/register?error=Username is already used", http.StatusSeeOther)
			return
		}

		user = models.User{Username: username, Email: email}
		user.SetPassword(password)

		result := db.Create(&user)
		if result.Error != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/success?m=Register successful", http.StatusSeeOther)
	})
}
