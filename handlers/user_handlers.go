package handlers

import (
	"fmt"
	"html/template"
	"myproject/models"
	"net/http"

	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		login := r.FormValue("login")
		password := r.FormValue("password")

		user := models.User{Login: login, Password: password}
		result := h.DB.Create(&user)
		if result.Error != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/register.html"))
	tmpl.Execute(w, nil)
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		login := r.FormValue("login")
		password := r.FormValue("password")

		var user models.User
		result := h.DB.Where("login = ?", login).First(&user)
		if result.Error != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		if !user.ComparePassword(password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: fmt.Sprintf("%d", user.ID),
			Path:  "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/login.html"))
	tmpl.Execute(w, nil)
}

func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// In a real application, you'd want to invalidate the session here
	http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
