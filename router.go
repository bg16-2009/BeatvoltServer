package main

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/bg16_2009/beatvolt_server/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func makeRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(jwtauth.Verifier(tokenAuth))

	r.Mount("/robot", robotRouter())

	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			hfn := func(w http.ResponseWriter, r *http.Request) {
				token, _, err := jwtauth.FromContext(r.Context())
				if err != nil || token == nil {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}
				next.ServeHTTP(w, r)
			}
			return http.HandlerFunc(hfn)
		})

		r.Mount("/admin", adminRouter())

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			u := models.User{}
			db.Where("id = ?", claims["userId"]).First(&u)

			renderTemplate("home", w, map[string]interface{}{
				"recycledBatteries": u.RecycledBatteries,
				"IsAdmin":           u.IsAdmin,
			})
		})

		r.Get("/redeem", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			u := models.User{}
			db.Where("id = ?", claims["userId"]).First(&u)
			_, redeemToken, _ := tokenAuth.Encode(map[string]interface{}{
				"id": strconv.Itoa(int(u.ID)),
			})

			renderTemplate("redeem", w, map[string]interface{}{
				"qrString": base64.StdEncoding.EncodeToString([]byte(redeemToken)),
				"IsAdmin":  u.IsAdmin,
			})
		})
		r.Get("/claim", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			u := models.User{}
			db.Where("id = ?", claims["userId"]).First(&u)

			renderTemplate("claim", w, map[string]interface{}{
				"Error":   r.URL.Query().Get("error"),
				"IsAdmin": u.IsAdmin,
			})
		})
		r.Post("/claim", func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}
			inputCode := r.FormValue("code")

			storeMutex.Lock()

			// Clean all expired codes
			expiredCodes := []string{}
			for _, code := range codeStore {
				if code.creationTime.Add(codeValidDuration).Before(time.Now()) {
					expiredCodes = append(expiredCodes, code.code)
				}
			}
			for _, code := range expiredCodes {
				delete(codeStore, code)
			}

			_, exists := codeStore[inputCode]
			if !exists {
				http.Redirect(w, r, "/claim?error=Code is invalid or claimed", http.StatusSeeOther)
				storeMutex.Unlock()
				return
			}

			_, claims, _ := jwtauth.FromContext(r.Context())
			u := models.User{}
			db.Where("id = ?", claims["userId"]).First(&u)
			db.Model(&u).Update("recycled_batteries", u.RecycledBatteries+1)
			http.Redirect(w, r, "/success?m=Code claimed succesfully", http.StatusSeeOther)
			delete(codeStore, inputCode)

			storeMutex.Unlock()
		})

		r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:     "jwt",
				Value:    "",
				Expires:  time.Unix(0, 0), // Delete cookie
				HttpOnly: true,
				Path:     "/",
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		})
	})

	r.Group(loginGroup)
	r.Group(registerGroup)
	r.Get("/success", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate("success", w, map[string]interface{}{
			"Message": r.URL.Query().Get("m"),
		})
	})

	return r
}
