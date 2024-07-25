package handlers

import (
	"html/template"
	"myproject/models"
	"net/http"

	"gorm.io/gorm"
)

type NoteHandler struct {
	DB *gorm.DB
}

func (h *NoteHandler) ListNotesHandler(w http.ResponseWriter, r *http.Request) {
	var notes []models.Note
	h.DB.Find(&notes)

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/main.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Notes":    notes,
		"LoggedIn": r.Context().Value("user") != nil,
	})
}

func (h *NoteHandler) AddNoteHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		content := r.FormValue("content")

		note := models.Note{Name: name, Content: content, UserID: user.ID}
		result := h.DB.Create(&note)
		if result.Error != nil {
			http.Error(w, "Error creating note", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/note_form.html"))
	tmpl.Execute(w, nil)
}
