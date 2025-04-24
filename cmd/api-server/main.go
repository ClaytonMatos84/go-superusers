package main

import (
	"log"
	"net/http"

	"github.com/ClaytonMatos84/go-superusers/internal"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	var users []internal.User
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		internal.UploadUsers(w, r, &users)
	}).Methods("POST")

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		internal.GetUsers(w, r, &users)
	}).Methods("GET")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(mux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
