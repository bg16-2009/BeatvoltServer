package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bg16_2009/beatvolt_server/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
)

func adminRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			user := models.User{}
			db.Where("id = ?", claims["userId"]).First(&user)
			if !user.IsAdmin {
				http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate("admin/index", w, map[string]interface{}{
			"IsAdmin": true,
		})
	})
	r.Get("/robots", func(w http.ResponseWriter, r *http.Request) {
		var bots []models.Robot
		result := db.Find(&bots)
		if result.Error != nil {
			fmt.Println("Error fetching users:", result.Error)
		}

		renderTemplate("admin/robots", w, map[string]interface{}{
			"bots":    bots,
			"IsAdmin": true,
		})
	})
	r.Post("/remove_batteries", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		userId, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, "Bad usage", http.StatusBadRequest)
			return
		}
		nr, err := strconv.Atoi(r.FormValue("n"))
		if err != nil || nr <= 0 {
			http.Error(w, "Bad usage", http.StatusBadRequest)
			return
		}

		u := models.User{}
		result := db.Where("id = ?", userId).First(&u)
		if result.Error != nil {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		db.Model(&u).Update("recycled_batteries", u.RecycledBatteries-nr)
		http.Redirect(w, r, "/admin/user?id="+strconv.Itoa(userId), http.StatusSeeOther)
	})
	r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, "Bad usage", http.StatusBadRequest)
			return
		}

		u := models.User{}
		result := db.Where("id = ?", userId).First(&u)
		if result.Error != nil {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}

		renderTemplate("admin/user", w, map[string]interface{}{
			"user":    u,
			"IsAdmin": true,
		})
	})
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		var users []models.User
		result := db.Find(&users)
		if result.Error != nil {
			fmt.Println("Error fetching users:", result.Error)
		}

		renderTemplate("admin/users", w, map[string]interface{}{
			"users":   users,
			"IsAdmin": true,
		})
	})
	r.Get("/validate", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate("admin/validate", w, map[string]interface{}{
			"IsAdmin": true,
		})
	})
	r.Post("/validate", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		tokenbase64 := r.FormValue("token")
		token, err := base64.StdEncoding.DecodeString(tokenbase64)
		if err != nil {
			http.Error(w, "Invalid code", http.StatusBadRequest)
			return
		}
		tokenString := string(token)

		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret_key), nil
		})

		if err != nil {
			http.Error(w, "Invalid code", http.StatusBadRequest)
			return
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid && claims["id"] != nil {
			http.Redirect(w, r, "/admin/user?id="+claims["id"].(string), http.StatusSeeOther)
		} else {
			http.Error(w, "Invalid code", http.StatusBadRequest)
			return
		}
	})
	return r
}
