package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"myproject/handlers"
	"myproject/models"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func authMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err == nil {
				var user models.User
				result := db.Where("id = ?", cookie.Value).First(&user)
				if result.Error == nil {
					ctx := context.WithValue(r.Context(), "user", &user)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func main() {

	dbHost := flag.String("db-host", "localhost", "Database host")
	dbUser := flag.String("db-user", "postgres", "Database user")
	dbPassword := flag.String("db-password", "your_password", "Database password")
	dbName := flag.String("db-name", "your_dbname", "Database name")
	dbPort := flag.String("db-port", "5432", "Database port")

	// Parse the flags
	flag.Parse()

	// Construct the DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		*dbHost, *dbUser, *dbPassword, *dbName, *dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.User{}, &models.Note{})

	userHandler := &handlers.UserHandler{DB: db}
	noteHandler := &handlers.NoteHandler{DB: db}

	r := mux.NewRouter()
	r.Use(authMiddleware(db))

	r.HandleFunc("/", noteHandler.ListNotesHandler)
	r.HandleFunc("/register", userHandler.RegisterHandler)
	r.HandleFunc("/login", userHandler.LoginHandler)
	r.HandleFunc("/logout", userHandler.LogoutHandler)
	r.HandleFunc("/add-note", noteHandler.AddNoteHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
