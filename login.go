package main

import (
	"net/http"

	"github.com/bg16_2009/beatvolt_server/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func loginGroup(r chi.Router) {
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		// Check if user is already logged in
		token, _, err := jwtauth.FromContext(r.Context())
		if err == nil || token != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		renderTemplate("login", w, map[string]string{
			"Error": r.URL.Query().Get("error"),
		})
	})
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user models.User
		if err := db.Where("email = ? OR username = ?", username, username).First(&user).Error; err == nil &&
			user.CheckPassword(password) {

			_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
				"userId":   user.ID,
				"userType": "user",
			})

			http.SetCookie(w, &http.Cookie{
				Name:     "jwt",
				Value:    tokenString,
				HttpOnly: true,
				Path:     "/",
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/login?error=Invalid username or password", http.StatusSeeOther)
	})
}
