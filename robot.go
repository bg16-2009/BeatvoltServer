package main

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bg16_2009/beatvolt_server/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func generateCode(charset string, length int) string {
	if charset == "" {
		charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	storeMutex.Lock()
	defer storeMutex.Unlock()

	for {
		code := make([]byte, length)
		for i := range code {
			code[i] = charset[rand.Intn(len(charset))]
		}
		if _, exists := codeStore[string(code)]; !exists {
			return string(code)
		}
	}
}

func robotRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")

		robot := models.Robot{}
		if err := db.Where("username = ?", username).First(&robot).Error; err == nil && robot.CheckPassword(password) {
			_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
				"userId":   robot.ID,
				"userType": "robot",
			})

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"token":   tokenString,
			})
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid credentials.",
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			hfn := func(w http.ResponseWriter, r *http.Request) {
				_, claims, _ := jwtauth.FromContext(r.Context())
				if claims["userType"] != "robot" {
					http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusSeeOther)
					return
				}
				next.ServeHTTP(w, r)
			}
			return http.HandlerFunc(hfn)
		})

		r.Group(func(r chi.Router) {
			r.Get("/generate_code", func(w http.ResponseWriter, r *http.Request) {
				code := generateCode("", 6)

				storeMutex.Lock()
				codeStore[code] = Code{
					code:         code,
					creationTime: time.Now(),
				}
				storeMutex.Unlock()

				_, claims, _ := jwtauth.FromContext(r.Context())
				bot := models.Robot{}
				db.Where("id = ?", claims["userId"]).First(&bot)
				db.Model(&bot).Update("collected_batteries", bot.CollectedBatteries+1)

				w.Write([]byte(code))
			})
			r.Get("/random_song", func(w http.ResponseWriter, r *http.Request) {
				files, err := filepath.Glob(filepath.Join("./songs", "*.json"))
				if err != nil || len(files) == 0 {
					http.Error(w, "No JSON files found", http.StatusInternalServerError)
					return
				}

				randFile := files[rand.Intn(len(files))]
				file, err := os.Open(randFile)
				if err != nil {
					http.Error(w, "Could not open file", http.StatusInternalServerError)
					return
				}
				defer file.Close()

				w.Header().Set("Content-Type", "application/json")
				io.Copy(w, file)
			})
		})
	})

	return r
}
